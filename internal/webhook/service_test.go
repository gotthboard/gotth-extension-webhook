package webhook

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	extensionsv1 "github.com/gotthboard/gotth-extension-webhook/proto/gotth/extensions/v1"
	pluginv1 "github.com/gotthboard/gotth-extension-webhook/proto/gotth/mail/plugin/v1"
	extensioncore "github.com/gotthboard/gotth-extensions/pkg/extensions"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func testRuntime(endpoint string) RuntimeConfig {
	return RuntimeConfig{
		Config:       Config{Schema: ConfigSchema, Endpoint: endpoint, TimeoutSeconds: 5},
		Binding:      Binding{Schema: "gotth.extension.runtime-binding.v1", InstanceID: "00000000-0000-4000-8000-000000000019", ManifestSHA256: strings.Repeat("1", 64), GrantSHA256: strings.Repeat("2", 64), SessionSHA256: strings.Repeat("3", 64)},
		ServiceToken: strings.Repeat("t", 32), HMACKey: []byte(strings.Repeat("k", 32)), Version: "1.0.0-alpha.1",
	}
}

func authorizedContext(t *testing.T, token string) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	return metadata.NewIncomingContext(ctx, metadata.Pairs("x-gotth-mail-service-token", token))
}

func testAlert() *pluginv1.AlertMessage {
	return &pluginv1.AlertMessage{Id: "alert-1", Class: "backup", Severity: "warning", Title: "Backup failed", Summary: "The bounded backup check failed.", CorrelationId: "corr-1", Resource: &pluginv1.AlertResource{Type: "backup", Id: "backup-1"}, Details: map[string]string{"z-last": "two", "a-first": "one"}}
}

func TestSendAlertHTTPSAndHMAC(t *testing.T) {
	key := []byte(strings.Repeat("k", 32))
	var gotBody []byte
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		gotBody, err = io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		mac := hmac.New(sha256.New, key)
		_, _ = mac.Write(gotBody)
		want := "sha256=" + hex.EncodeToString(mac.Sum(nil))
		if r.Method != http.MethodPost || r.Header.Get("Content-Type") != "application/json" || r.Header.Get("X-GOTTH-Signature") != want || r.Header.Get("X-GOTTH-Alert-ID") != "alert-1" || r.Header.Get("X-GOTTH-Correlation-ID") != "corr-1" {
			t.Errorf("invalid webhook request: method=%s headers=%v", r.Method, r.Header)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	runtime := testRuntime(server.URL)
	service, err := NewService(runtime)
	if err != nil {
		t.Fatal(err)
	}
	defer service.Close()
	service.client = server.Client()
	service.client.Timeout = runtime.Config.Timeout()
	service.now = func() time.Time { return time.Date(2026, 9, 20, 16, 0, 0, 0, time.UTC) }
	response, err := service.SendAlert(authorizedContext(t, runtime.ServiceToken), &pluginv1.SendAlertRequest{Alert: testAlert()})
	if err != nil {
		t.Fatal(err)
	}
	if response.GetStatus() != "delivered" || response.GetReason() != "receiver_accepted" || response.GetEvidence().GetTransport() != "https-webhook" || response.GetEvidence().GetVerificationResult() != "hmac_sha256" {
		t.Fatalf("unexpected delivery response: %#v", response)
	}
	body := string(gotBody)
	if strings.Index(body, "a-first") > strings.Index(body, "z-last") || strings.Contains(body, strings.Repeat("k", 32)) || !strings.Contains(body, `"schema":"gotth.extension.webhook.alert.v1"`) {
		t.Fatalf("body is not deterministic and redacted: %s", body)
	}
}

func TestSendAlertRedirectIsUnavailableAndTargetIsNotReached(t *testing.T) {
	targetHit := make(chan struct{}, 1)
	target := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		select {
		case targetHit <- struct{}{}:
		default:
		}
	}))
	defer target.Close()

	source := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Location", target.URL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}))
	defer source.Close()

	runtime := testRuntime(source.URL)
	service, err := NewService(runtime)
	if err != nil {
		t.Fatal(err)
	}
	redirectPolicy := service.client.CheckRedirect
	client := source.Client()
	client.Timeout = runtime.Config.Timeout()
	client.CheckRedirect = redirectPolicy
	service.client = client
	defer service.Close()

	response, err := service.SendAlert(authorizedContext(t, runtime.ServiceToken), &pluginv1.SendAlertRequest{Alert: testAlert()})
	if response != nil || status.Code(err) != codes.Unavailable {
		t.Fatalf("redirect did not fail closed: response=%#v err=%v", response, err)
	}
	select {
	case <-targetHit:
		t.Fatal("redirect target was reached")
	default:
	}
}

func TestSendAlertFailuresAreBounded(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Secret", "must-not-return")
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte("private receiver failure"))
	}))
	defer server.Close()
	runtime := testRuntime(server.URL)
	service, _ := NewService(runtime)
	defer service.Close()
	service.client = server.Client()
	response, err := service.SendAlert(authorizedContext(t, runtime.ServiceToken), &pluginv1.SendAlertRequest{Alert: testAlert()})
	if err != nil || response.GetStatus() != "failed_permanent" || response.GetReason() != "receiver_rejected" || response.GetEvidence() != nil {
		t.Fatalf("receiver rejection was not bounded: %#v %v", response, err)
	}
	if _, err := service.SendAlert(authorizedContext(t, "wrong-token-wrong-token-wrong-token"), &pluginv1.SendAlertRequest{Alert: testAlert()}); status.Code(err) != codes.Unauthenticated {
		t.Fatalf("wrong token accepted: %v", err)
	}
	if _, err := service.SendAlert(context.Background(), &pluginv1.SendAlertRequest{Alert: testAlert()}); status.Code(err) != codes.InvalidArgument {
		t.Fatalf("missing deadline accepted: %v", err)
	}
	hostile := testAlert()
	hostile.Summary = "line one\nline two"
	if _, err := service.SendAlert(authorizedContext(t, runtime.ServiceToken), &pluginv1.SendAlertRequest{Alert: hostile}); status.Code(err) != codes.InvalidArgument {
		t.Fatalf("hostile alert accepted: %v", err)
	}
	if _, err := service.SendPrompt(authorizedContext(t, runtime.ServiceToken), &pluginv1.SendPromptRequest{}); status.Code(err) != codes.Unimplemented {
		t.Fatalf("prompt authority unexpectedly available: %v", err)
	}
}

func TestSendAlertPostBodyCloseIsNotDelivered(t *testing.T) {
	received := make(chan []byte, 1)
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		received <- body
		hijacker, ok := w.(http.Hijacker)
		if !ok {
			t.Error("response writer cannot inject a post-body close")
			return
		}
		conn, _, err := hijacker.Hijack()
		if err != nil {
			t.Error(err)
			return
		}
		_ = conn.Close()
	}))
	defer server.Close()

	runtime := testRuntime(server.URL)
	service, err := NewService(runtime)
	if err != nil {
		t.Fatal(err)
	}
	defer service.Close()
	service.client = server.Client()
	service.client.Timeout = runtime.Config.Timeout()

	response, err := service.SendAlert(authorizedContext(t, runtime.ServiceToken), &pluginv1.SendAlertRequest{Alert: testAlert()})
	if response != nil || status.Code(err) != codes.Unavailable {
		t.Fatalf("ambiguous post-body close claimed delivery: response=%#v err=%v", response, err)
	}
	select {
	case body := <-received:
		if len(body) == 0 || !strings.Contains(string(body), `"id":"alert-1"`) {
			t.Fatalf("receiver did not consume the complete alert before closing: %q", body)
		}
	case <-time.After(time.Second):
		t.Fatal("receiver did not consume the alert before the transport failure")
	}
}

func TestControlBindingAndManifest(t *testing.T) {
	runtime := testRuntime("https://receiver.example.test/hook")
	service, err := NewService(runtime)
	if err != nil {
		t.Fatal(err)
	}
	defer service.Close()
	ctx := authorizedContext(t, runtime.ServiceToken)
	request := &extensionsv1.HandshakeRequest{CorrelationId: "corr-1", Challenge: []byte(strings.Repeat("c", 32)), ExpectedExtensionId: ExtensionID, ExpectedManifestSha256: runtime.Binding.ManifestSHA256, ExpectedGrantSha256: runtime.Binding.GrantSHA256, ExpectedSessionSha256: runtime.Binding.SessionSHA256}
	response, err := service.Handshake(ctx, request)
	if err != nil || response.GetExtensionId() != ExtensionID || response.GetExtensionVersion() != runtime.Version || len(response.GetInterfaces()) != 1 || len(response.GetCapabilities()) != len(Capabilities) {
		t.Fatalf("valid handshake rejected: %#v %v", response, err)
	}
	request.ExpectedSessionSha256 = strings.Repeat("0", 64)
	if _, err := service.Handshake(ctx, request); status.Code(err) != codes.FailedPrecondition {
		t.Fatalf("binding mismatch accepted: %v", err)
	}
	health, err := service.Health(ctx, &extensionsv1.HealthRequest{CorrelationId: "corr-1", SessionSha256: runtime.Binding.SessionSHA256})
	if err != nil || health.GetState() != extensionsv1.HealthState_HEALTH_STATE_READY || health.GetCode() != "extension.ready" {
		t.Fatalf("health unavailable: %#v %v", health, err)
	}
	manifest, err := Manifest(runtime.Version)
	if err != nil || extensioncore.ValidateManifest(manifest) != nil || manifest.ID != ExtensionID || manifest.Secrets[0].ID != SecretSlot {
		t.Fatalf("invalid foundation manifest: %#v %v", manifest, err)
	}
}
