package webhook

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDecodeConfigAndBindingRejectHostileInput(t *testing.T) {
	valid := `{"schema":"gotth.extension.webhook.config.v1","endpoint":"https://receiver.example.test/hook","timeout_seconds":10}`
	if config, err := decodeConfig([]byte(valid)); err != nil || config.TimeoutSeconds != 10 {
		t.Fatalf("valid config rejected: %#v %v", config, err)
	}
	for name, input := range map[string]string{
		"duplicate":   `{"schema":"gotth.extension.webhook.config.v1","endpoint":"https://receiver.example.test","endpoint":"https://other.example.test","timeout_seconds":10}`,
		"unknown":     `{"schema":"gotth.extension.webhook.config.v1","endpoint":"https://receiver.example.test","timeout_seconds":10,"extra":true}`,
		"http":        `{"schema":"gotth.extension.webhook.config.v1","endpoint":"http://receiver.example.test","timeout_seconds":10}`,
		"credentials": `{"schema":"gotth.extension.webhook.config.v1","endpoint":"https://user:pass@receiver.example.test","timeout_seconds":10}`,
		"fragment":    `{"schema":"gotth.extension.webhook.config.v1","endpoint":"https://receiver.example.test/#secret","timeout_seconds":10}`,
		"timeout":     `{"schema":"gotth.extension.webhook.config.v1","endpoint":"https://receiver.example.test","timeout_seconds":31}`,
		"trailing":    valid + `{}`,
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := decodeConfig([]byte(input)); err == nil {
				t.Fatal("hostile configuration accepted")
			}
		})
	}
	binding := `{"schema":"gotth.extension.runtime-binding.v1","instance_id":"00000000-0000-4000-8000-000000000019","manifest_sha256":"` + strings.Repeat("1", 64) + `","grant_sha256":"` + strings.Repeat("2", 64) + `","session_sha256":"` + strings.Repeat("3", 64) + `"}`
	if _, err := decodeBinding([]byte(binding)); err != nil {
		t.Fatalf("valid binding rejected: %v", err)
	}
	if _, err := decodeBinding([]byte(strings.Replace(binding, strings.Repeat("3", 64), strings.Repeat("A", 64), 1))); err == nil {
		t.Fatal("uppercase digest accepted")
	}
}

func TestLoadRuntimeFromEnvPrivateFileBoundary(t *testing.T) {
	dir := t.TempDir()
	if err := os.Chmod(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	secretDir := filepath.Join(dir, "secrets")
	if err := os.Mkdir(secretDir, 0o700); err != nil {
		t.Fatal(err)
	}
	writePrivate := func(name, value string) string {
		t.Helper()
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(value), 0o600); err != nil {
			t.Fatal(err)
		}
		return path
	}
	config := writePrivate("config.json", `{"schema":"gotth.extension.webhook.config.v1","endpoint":"https://receiver.example.test/hook","timeout_seconds":10}`)
	previousVersion := BuildVersion
	BuildVersion = "1.0.0-alpha.2"
	t.Cleanup(func() { BuildVersion = previousVersion })
	manifestDigest, err := ManifestDigest(BuildVersion)
	if err != nil {
		t.Fatal(err)
	}
	binding := writePrivate("binding.json", `{"schema":"gotth.extension.runtime-binding.v1","instance_id":"00000000-0000-4000-8000-000000000019","manifest_sha256":"`+manifestDigest+`","grant_sha256":"`+strings.Repeat("2", 64)+`","session_sha256":"`+strings.Repeat("3", 64)+`"}`)
	token := writePrivate("token", strings.Repeat("t", 32))
	if err := os.WriteFile(filepath.Join(secretDir, SecretSlot), []byte(strings.Repeat("k", 32)), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GOTTH_EXTENSION_CONFIG_FILE", config)
	t.Setenv("GOTTH_EXTENSION_BINDING_FILE", binding)
	t.Setenv("GOTTH_EXTENSION_SECRET_DIR", secretDir)
	t.Setenv("GOTTH_EXTENSION_SERVICE_TOKEN_FILE", token)
	t.Setenv("GOTTH_EXTENSION_SOCKET", filepath.Join(dir, "extension.sock"))
	runtime, err := LoadRuntimeFromEnv()
	if err != nil {
		t.Fatalf("valid runtime rejected: %v", err)
	}
	defer clear(runtime.HMACKey)
	if runtime.Binding.InstanceID == "" || len(runtime.HMACKey) != 32 || runtime.ServiceToken == "" {
		t.Fatalf("incomplete runtime: %#v", runtime)
	}

	if err := os.Chmod(config, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadRuntimeFromEnv(); err == nil {
		t.Fatal("world-readable configuration accepted")
	}
	if err := os.Chmod(config, 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "config-link")
	if err := os.Symlink(config, link); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GOTTH_EXTENSION_CONFIG_FILE", link)
	if _, err := LoadRuntimeFromEnv(); err == nil {
		t.Fatal("symlinked configuration accepted")
	}
}

func TestVersionGrammar(t *testing.T) {
	for _, value := range []string{"dev", "1.0.0-alpha.1", "1.0.0-beta.12", "1.0.0"} {
		if !validVersion(value) {
			t.Fatalf("valid version rejected: %s", value)
		}
	}
	for _, value := range []string{"", "1.0.0-alpha.0", "1.0.0-alpha.01", "v1.0.0", "1.0.1", "1.0.0-rc.1"} {
		if validVersion(value) {
			t.Fatalf("invalid version accepted: %s", value)
		}
	}
}
