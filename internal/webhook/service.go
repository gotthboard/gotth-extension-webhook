package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	pluginv1 "github.com/gotthboard/gotth-extension-webhook/proto/gotth/mail/plugin/v1"
	extensionsv1 "github.com/gotthboard/gotth-sdk/gen/go/gotth/extensions/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const maxBodyBytes = 64 << 10

type Service struct {
	extensionsv1.UnimplementedExtensionControlServer
	pluginv1.UnimplementedPluginControlServer
	pluginv1.UnimplementedNotificationBackendServer

	config  Config
	binding Binding
	token   string
	key     []byte
	version string
	client  *http.Client
	now     func() time.Time
}

func NewService(runtime RuntimeConfig) (*Service, error) {
	if runtime.ServiceToken == "" || len(runtime.HMACKey) < 32 || !validVersion(runtime.Version) {
		return nil, errors.New("complete webhook runtime required")
	}
	transport := &http.Transport{TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12}}
	client := &http.Client{
		Transport: transport,
		Timeout:   runtime.Config.Timeout(),
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return errors.New("webhook redirects are disabled")
		},
	}
	return &Service{config: runtime.Config, binding: runtime.Binding, token: runtime.ServiceToken, key: append([]byte(nil), runtime.HMACKey...), version: runtime.Version, client: client, now: time.Now}, nil
}

func (s *Service) Close() {
	clear(s.key)
	if transport, ok := s.client.Transport.(*http.Transport); ok {
		transport.CloseIdleConnections()
	}
}

func (s *Service) authorize(ctx context.Context, correlationID string) error {
	if _, ok := ctx.Deadline(); !ok {
		return status.Error(codes.InvalidArgument, "deadline required")
	}
	if !validScalar(correlationID, 1, 256) {
		return status.Error(codes.InvalidArgument, "valid correlation ID required")
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "service identity required")
	}
	values := md.Get("x-gotth-mail-service-token")
	if len(values) != 1 || subtle.ConstantTimeCompare([]byte(values[0]), []byte(s.token)) != 1 {
		return status.Error(codes.Unauthenticated, "invalid service identity")
	}
	return nil
}

func (s *Service) Handshake(ctx context.Context, in *extensionsv1.HandshakeRequest) (*extensionsv1.HandshakeResponse, error) {
	if in == nil || len(in.GetChallenge()) != 32 {
		return nil, status.Error(codes.InvalidArgument, "valid extension handshake required")
	}
	if err := s.authorize(ctx, in.GetCorrelationId()); err != nil {
		return nil, err
	}
	if in.GetExpectedExtensionId() != ExtensionID || in.GetExpectedManifestSha256() != s.binding.ManifestSHA256 || in.GetExpectedGrantSha256() != s.binding.GrantSHA256 || in.GetExpectedSessionSha256() != s.binding.SessionSHA256 {
		return nil, status.Error(codes.FailedPrecondition, "extension binding mismatch")
	}
	return &extensionsv1.HandshakeResponse{
		Challenge:        append([]byte(nil), in.GetChallenge()...),
		ExtensionId:      ExtensionID,
		ExtensionVersion: s.version,
		ManifestSha256:   s.binding.ManifestSHA256,
		GrantSha256:      s.binding.GrantSHA256,
		SessionSha256:    s.binding.SessionSHA256,
		Control:          &extensionsv1.ProtocolVersion{Name: "gotth.extensions.control", Major: 1, Minor: 0},
		Interfaces:       []*extensionsv1.InterfaceVersion{{Name: Interface, Major: 1, Minor: 0}},
		Capabilities:     append([]string(nil), Capabilities...),
	}, nil
}

func (s *Service) Health(ctx context.Context, in *extensionsv1.HealthRequest) (*extensionsv1.HealthResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "valid health request required")
	}
	if err := s.authorize(ctx, in.GetCorrelationId()); err != nil {
		return nil, err
	}
	if in.GetSessionSha256() != s.binding.SessionSHA256 {
		return nil, status.Error(codes.FailedPrecondition, "extension session mismatch")
	}
	return &extensionsv1.HealthResponse{State: extensionsv1.HealthState_HEALTH_STATE_READY, Code: "extension.ready"}, nil
}

func (s *Service) HealthPlugin(ctx context.Context, in *pluginv1.HealthRequest) (*pluginv1.HealthResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "valid health request required")
	}
	if err := s.authorize(ctx, in.GetCorrelationId()); err != nil {
		return nil, err
	}
	return &pluginv1.HealthResponse{Healthy: true, Message: "extension.ready"}, nil
}

func (s *Service) Version(ctx context.Context, in *pluginv1.VersionRequest) (*pluginv1.VersionResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "valid version request required")
	}
	if err := s.authorize(ctx, in.GetCorrelationId()); err != nil {
		return nil, err
	}
	return &pluginv1.VersionResponse{Name: ExtensionID, Version: s.version}, nil
}

func (s *Service) Capabilities(ctx context.Context, in *pluginv1.CapabilitiesRequest) (*pluginv1.CapabilitiesResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "valid capabilities request required")
	}
	if err := s.authorize(ctx, in.GetCorrelationId()); err != nil {
		return nil, err
	}
	return &pluginv1.CapabilitiesResponse{Capabilities: append([]string(nil), Capabilities...)}, nil
}

type detail struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type alertBody struct {
	Schema        string   `json:"schema"`
	ID            string   `json:"id"`
	Class         string   `json:"class"`
	Severity      string   `json:"severity"`
	Title         string   `json:"title"`
	Summary       string   `json:"summary"`
	CorrelationID string   `json:"correlation_id"`
	ResourceType  string   `json:"resource_type,omitempty"`
	ResourceID    string   `json:"resource_id,omitempty"`
	Details       []detail `json:"details,omitempty"`
}

func (s *Service) SendAlert(ctx context.Context, in *pluginv1.SendAlertRequest) (*pluginv1.DeliveryResponse, error) {
	if in == nil || in.GetAlert() == nil {
		return nil, status.Error(codes.InvalidArgument, "valid alert required")
	}
	a := in.GetAlert()
	if err := s.authorize(ctx, a.GetCorrelationId()); err != nil {
		return nil, err
	}
	body, err := canonicalAlert(a)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "bounded alert required")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.config.Endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, status.Error(codes.Internal, "webhook request unavailable")
	}
	mac := hmac.New(sha256.New, s.key)
	_, _ = mac.Write(body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GOTTH-Alert-ID", a.GetId())
	req.Header.Set("X-GOTTH-Correlation-ID", a.GetCorrelationId())
	req.Header.Set("X-GOTTH-Signature", "sha256="+hex.EncodeToString(mac.Sum(nil)))
	response, err := s.client.Do(req)
	if err != nil {
		return nil, status.Error(codes.Unavailable, "webhook delivery failed")
	}
	defer response.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, 4096))
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return &pluginv1.DeliveryResponse{Status: "failed_permanent", Reason: "receiver_rejected"}, nil
	}
	generated := s.now().UTC().Format(time.RFC3339Nano)
	return &pluginv1.DeliveryResponse{Status: "delivered", Reason: "receiver_accepted", Evidence: &pluginv1.DeliveryEvidence{Transport: "https-webhook", MessageId: a.GetId(), GeneratedAt: generated, VerificationResult: "hmac_sha256", Workflow: "notification.alert"}}, nil
}

func (s *Service) SendPrompt(context.Context, *pluginv1.SendPromptRequest) (*pluginv1.PromptResponse, error) {
	return nil, status.Error(codes.Unimplemented, "webhook extension does not support prompts")
}

func canonicalAlert(a *pluginv1.AlertMessage) ([]byte, error) {
	if !validScalar(a.GetId(), 1, 256) || !validScalar(a.GetClass(), 1, 128) || !validScalar(a.GetSeverity(), 1, 32) || !validScalar(a.GetTitle(), 1, 256) || !validScalar(a.GetSummary(), 1, 4096) || !validScalar(a.GetCorrelationId(), 1, 256) || len(a.GetDetails()) > 32 {
		return nil, errors.New("invalid alert")
	}
	body := alertBody{Schema: "gotth.extension.webhook.alert.v1", ID: a.GetId(), Class: a.GetClass(), Severity: a.GetSeverity(), Title: a.GetTitle(), Summary: a.GetSummary(), CorrelationID: a.GetCorrelationId()}
	if resource := a.GetResource(); resource != nil {
		if !validScalar(resource.GetType(), 1, 128) || !validScalar(resource.GetId(), 1, 256) {
			return nil, errors.New("invalid resource")
		}
		body.ResourceType, body.ResourceID = resource.GetType(), resource.GetId()
	}
	keys := make([]string, 0, len(a.GetDetails()))
	for key, value := range a.GetDetails() {
		if !validScalar(key, 1, 128) || !validScalar(value, 0, 2048) {
			return nil, errors.New("invalid details")
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		body.Details = append(body.Details, detail{Key: key, Value: a.GetDetails()[key]})
	}
	encoded, err := json.Marshal(body)
	if err != nil || len(encoded) > maxBodyBytes {
		return nil, fmt.Errorf("bounded alert body required")
	}
	return encoded, nil
}

func validScalar(value string, min, max int) bool {
	if len(value) < min || len(value) > max || !utf8.ValidString(value) {
		return false
	}
	for _, r := range value {
		if r == 0 || r == '\r' || r == '\n' || r < 0x20 && r != '\t' {
			return false
		}
	}
	return !strings.Contains(value, "\u2028") && !strings.Contains(value, "\u2029")
}
