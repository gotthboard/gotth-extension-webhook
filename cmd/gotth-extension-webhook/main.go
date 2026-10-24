package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gotthboard/gotth-extension-webhook/internal/webhook"
	pluginv1 "github.com/gotthboard/gotth-extension-webhook/proto/gotth/mail/plugin/v1"
	extensionsv1 "github.com/gotthboard/gotth-sdk/gen/go/gotth/extensions/v1"
	"google.golang.org/grpc"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Println(webhook.BuildVersion)
		return
	}
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "gotth-extension-webhook: startup failed")
		os.Exit(1)
	}
}

func run() error {
	runtime, err := webhook.LoadRuntimeFromEnv()
	if err != nil {
		return err
	}
	defer clear(runtime.HMACKey)
	service, err := webhook.NewService(runtime)
	if err != nil {
		return err
	}
	defer service.Close()
	listener, err := net.Listen("unix", runtime.SocketPath)
	if err != nil {
		return err
	}
	defer listener.Close()
	if err := os.Chmod(runtime.SocketPath, 0o600); err != nil {
		return err
	}
	server := grpc.NewServer()
	extensionsv1.RegisterExtensionControlServer(server, service)
	pluginv1.RegisterPluginControlServer(server, pluginControlAdapter{service: service})
	pluginv1.RegisterNotificationBackendServer(server, service)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signals)
	go func() {
		<-signals
		server.GracefulStop()
	}()
	if err := server.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return err
	}
	return nil
}

type pluginControlAdapter struct {
	service *webhook.Service
	pluginv1.UnimplementedPluginControlServer
}

func (a pluginControlAdapter) Health(ctx context.Context, in *pluginv1.HealthRequest) (*pluginv1.HealthResponse, error) {
	return a.service.HealthPlugin(ctx, in)
}

func (a pluginControlAdapter) Version(ctx context.Context, in *pluginv1.VersionRequest) (*pluginv1.VersionResponse, error) {
	return a.service.Version(ctx, in)
}

func (a pluginControlAdapter) Capabilities(ctx context.Context, in *pluginv1.CapabilitiesRequest) (*pluginv1.CapabilitiesResponse, error) {
	return a.service.Capabilities(ctx, in)
}
