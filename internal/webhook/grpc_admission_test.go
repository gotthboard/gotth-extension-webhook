package webhook

import (
	"context"
	"net"
	"path/filepath"
	"testing"
	"time"

	pluginv1 "github.com/gotthboard/gotth-extension-webhook/proto/gotth/mail/plugin/v1"
	extensionsv1 "github.com/gotthboard/gotth-sdk/gen/go/gotth/extensions/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type testPluginControlAdapter struct {
	pluginv1.UnimplementedPluginControlServer
	service *Service
}

func (a testPluginControlAdapter) Health(ctx context.Context, in *pluginv1.HealthRequest) (*pluginv1.HealthResponse, error) {
	return a.service.HealthPlugin(ctx, in)
}

func (a testPluginControlAdapter) Version(ctx context.Context, in *pluginv1.VersionRequest) (*pluginv1.VersionResponse, error) {
	return a.service.Version(ctx, in)
}

func (a testPluginControlAdapter) Capabilities(ctx context.Context, in *pluginv1.CapabilitiesRequest) (*pluginv1.CapabilitiesResponse, error) {
	return a.service.Capabilities(ctx, in)
}

func TestUnaryAdmissionPrecedesEveryRPCHandler(t *testing.T) {
	runtime := testRuntime("https://receiver.example.test/hook")
	service, err := NewService(runtime)
	if err != nil {
		t.Fatal(err)
	}
	defer service.Close()

	socket := filepath.Join(t.TempDir(), "webhook.sock")
	listener, err := net.Listen("unix", socket)
	if err != nil {
		t.Fatal(err)
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(service.UnaryServerInterceptor()))
	extensionsv1.RegisterExtensionControlServer(server, service)
	pluginv1.RegisterPluginControlServer(server, testPluginControlAdapter{service: service})
	pluginv1.RegisterNotificationBackendServer(server, service)
	serveDone := make(chan error, 1)
	go func() { serveDone <- server.Serve(listener) }()
	defer func() {
		server.Stop()
		if err := <-serveDone; err != nil && err != grpc.ErrServerStopped {
			t.Errorf("serve: %v", err)
		}
	}()

	connection, err := grpc.NewClient(
		"passthrough:///webhook",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return net.Dial("unix", socket) }),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Close()
	control := extensionsv1.NewExtensionControlClient(connection)
	plugin := pluginv1.NewPluginControlClient(connection)
	notification := pluginv1.NewNotificationBackendClient(connection)

	calls := []struct {
		name string
		call func(context.Context) error
	}{
		{"extension handshake", func(ctx context.Context) error {
			_, err := control.Handshake(ctx, &extensionsv1.HandshakeRequest{})
			return err
		}},
		{"extension health", func(ctx context.Context) error {
			_, err := control.Health(ctx, &extensionsv1.HealthRequest{})
			return err
		}},
		{"plugin health", func(ctx context.Context) error { _, err := plugin.Health(ctx, &pluginv1.HealthRequest{}); return err }},
		{"plugin version", func(ctx context.Context) error { _, err := plugin.Version(ctx, &pluginv1.VersionRequest{}); return err }},
		{"plugin capabilities", func(ctx context.Context) error {
			_, err := plugin.Capabilities(ctx, &pluginv1.CapabilitiesRequest{})
			return err
		}},
		{"send alert", func(ctx context.Context) error {
			_, err := notification.SendAlert(ctx, &pluginv1.SendAlertRequest{})
			return err
		}},
		{"send prompt", func(ctx context.Context) error {
			_, err := notification.SendPrompt(ctx, &pluginv1.SendPromptRequest{})
			return err
		}},
	}
	for _, test := range calls {
		t.Run(test.name+" rejects missing identity before handler", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if code := status.Code(test.call(ctx)); code != codes.Unauthenticated {
				t.Fatalf("got %s, want Unauthenticated", code)
			}
		})
		t.Run(test.name+" rejects missing deadline before handler", func(t *testing.T) {
			ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("x-gotth-mail-service-token", runtime.ServiceToken))
			if code := status.Code(test.call(ctx)); code != codes.InvalidArgument {
				t.Fatalf("got %s, want InvalidArgument", code)
			}
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("x-gotth-mail-service-token", runtime.ServiceToken))
	_, err = notification.SendPrompt(ctx, &pluginv1.SendPromptRequest{CorrelationId: "corr-1"})
	if status.Code(err) != codes.Unimplemented {
		t.Fatalf("admitted prompt got %v, want Unimplemented", err)
	}
}
