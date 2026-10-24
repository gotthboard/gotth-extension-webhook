package webhook

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testBindingJSON(t *testing.T) string {
	t.Helper()
	binding := testRuntime("https://receiver.example.test/hook").Binding
	document := map[string]any{
		"schema": binding.Schema, "instance_id": binding.InstanceID,
		"manifest_sha256": binding.ManifestSHA256, "grant_sha256": binding.GrantSHA256, "session_sha256": binding.SessionSHA256,
		"grant": binding.Grant,
		"host":  map[string]any{"protocols": binding.Host.Protocols, "interfaces": binding.Host.Interfaces},
	}
	encoded, err := json.Marshal(document)
	if err != nil {
		t.Fatal(err)
	}
	return string(encoded)
}

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
	invalidUTF8Config := append([]byte(`{"schema":"gotth.extension.webhook.config.v1","endpoint":"https://receiver.example.test/`), 0xff)
	invalidUTF8Config = append(invalidUTF8Config, []byte(`","timeout_seconds":10}`)...)
	if _, err := decodeConfig(invalidUTF8Config); err == nil {
		t.Fatal("configuration containing raw invalid UTF-8 accepted")
	}
	binding := testBindingJSON(t)
	decoded, err := decodeBinding([]byte(binding))
	if err != nil {
		t.Fatalf("valid binding rejected: %v", err)
	}
	if _, err := decodeBinding([]byte(strings.Replace(binding, decoded.SessionSHA256, strings.ToUpper(decoded.SessionSHA256), 1))); err == nil {
		t.Fatal("uppercase digest accepted")
	}
	duplicateHostRange := strings.Replace(binding, `"max_minor":0`, `"max_minor":0,"max_minor":0`, 1)
	if _, err := decodeBinding([]byte(duplicateHostRange)); err == nil {
		t.Fatal("duplicate host-profile field accepted")
	}
	bindingBytes := []byte(binding)
	instanceMarker := []byte(`"instance_id":"`)
	instanceOffset := bytes.Index(bindingBytes, instanceMarker)
	if instanceOffset < 0 {
		t.Fatal("test binding lacks instance_id")
	}
	instanceOffset += len(instanceMarker)
	invalidUTF8Binding := append([]byte(nil), bindingBytes[:instanceOffset]...)
	invalidUTF8Binding = append(invalidUTF8Binding, 0xff)
	invalidUTF8Binding = append(invalidUTF8Binding, bindingBytes[instanceOffset:]...)
	if _, err := decodeBinding(invalidUTF8Binding); err == nil {
		t.Fatal("binding containing raw invalid UTF-8 accepted")
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
	BuildVersion = "1.0.0"
	t.Cleanup(func() { BuildVersion = previousVersion })
	binding := writePrivate("binding.json", testBindingJSON(t))
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
	if runtime.Binding.InstanceID == "" || runtime.Session.Fingerprint != runtime.Binding.SessionSHA256 || len(runtime.HMACKey) != 32 || runtime.ServiceToken == "" {
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
