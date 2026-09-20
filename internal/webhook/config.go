package webhook

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	ConfigSchema = "gotth.extension.webhook.config.v1"
	ExtensionID  = "gotth.mail.notification.webhook"
	Interface    = "gotth.mail.interface.notification"
	SecretSlot   = "webhook.hmac-key"
	maxFileBytes = 64 << 10
)

var Capabilities = []string{
	"notification.alert.send",
	"notification.alert.sink",
	"notification.delivery.status",
}

var BuildVersion = "dev"

type Config struct {
	Schema         string `json:"schema"`
	Endpoint       string `json:"endpoint"`
	TimeoutSeconds int    `json:"timeout_seconds"`
}

type Binding struct {
	Schema         string `json:"schema"`
	InstanceID     string `json:"instance_id"`
	ManifestSHA256 string `json:"manifest_sha256"`
	GrantSHA256    string `json:"grant_sha256"`
	SessionSHA256  string `json:"session_sha256"`
}

type RuntimeConfig struct {
	Config       Config
	Binding      Binding
	ServiceToken string
	HMACKey      []byte
	SocketPath   string
	Version      string
}

func LoadRuntimeFromEnv() (RuntimeConfig, error) {
	configPath := strings.TrimSpace(os.Getenv("GOTTH_EXTENSION_CONFIG_FILE"))
	bindingPath := strings.TrimSpace(os.Getenv("GOTTH_EXTENSION_BINDING_FILE"))
	secretDir := strings.TrimSpace(os.Getenv("GOTTH_EXTENSION_SECRET_DIR"))
	tokenPath := strings.TrimSpace(os.Getenv("GOTTH_EXTENSION_SERVICE_TOKEN_FILE"))
	socketPath := strings.TrimSpace(os.Getenv("GOTTH_EXTENSION_SOCKET"))
	if configPath == "" || bindingPath == "" || secretDir == "" || tokenPath == "" || socketPath == "" || !validVersion(BuildVersion) {
		return RuntimeConfig{}, errors.New("complete extension runtime configuration required")
	}
	configData, err := readPrivateFile(configPath, maxFileBytes)
	if err != nil {
		return RuntimeConfig{}, fmt.Errorf("read configuration: %w", err)
	}
	bindingData, err := readPrivateFile(bindingPath, maxFileBytes)
	if err != nil {
		return RuntimeConfig{}, fmt.Errorf("read binding: %w", err)
	}
	tokenData, err := readPrivateFile(tokenPath, 1024)
	if err != nil {
		return RuntimeConfig{}, fmt.Errorf("read service token: %w", err)
	}
	secretInfo, err := os.Lstat(secretDir)
	if err != nil || secretInfo.Mode()&os.ModeSymlink != 0 || !secretInfo.IsDir() || secretInfo.Mode().Perm()&0o077 != 0 {
		return RuntimeConfig{}, errors.New("private secret directory required")
	}
	hmacKey, err := readPrivateFile(filepath.Join(secretDir, SecretSlot), 4096)
	if err != nil || len(hmacKey) < 32 {
		clear(hmacKey)
		return RuntimeConfig{}, errors.New("valid webhook HMAC key required")
	}
	config, err := decodeConfig(configData)
	if err != nil {
		clear(hmacKey)
		return RuntimeConfig{}, err
	}
	binding, err := decodeBinding(bindingData)
	if err != nil {
		clear(hmacKey)
		return RuntimeConfig{}, err
	}
	manifestDigest, err := ManifestDigest(BuildVersion)
	if err != nil || manifestDigest != binding.ManifestSHA256 {
		clear(hmacKey)
		return RuntimeConfig{}, errors.New("artifact manifest binding mismatch")
	}
	token := strings.TrimSpace(string(tokenData))
	if !validToken(token, 32, 256) {
		clear(hmacKey)
		return RuntimeConfig{}, errors.New("valid service token required")
	}
	if err := validateSocketPath(socketPath); err != nil {
		clear(hmacKey)
		return RuntimeConfig{}, err
	}
	return RuntimeConfig{Config: config, Binding: binding, ServiceToken: token, HMACKey: hmacKey, SocketPath: socketPath, Version: BuildVersion}, nil
}

func decodeConfig(data []byte) (Config, error) {
	fields, err := decodeUniqueObject(data)
	if err != nil || len(fields) != 3 {
		return Config{}, errors.New("invalid extension configuration")
	}
	var c Config
	if err := decodeField(fields, "schema", &c.Schema); err != nil || c.Schema != ConfigSchema {
		return Config{}, errors.New("invalid configuration schema")
	}
	if err := decodeField(fields, "endpoint", &c.Endpoint); err != nil {
		return Config{}, errors.New("valid webhook endpoint required")
	}
	if err := decodeField(fields, "timeout_seconds", &c.TimeoutSeconds); err != nil || c.TimeoutSeconds < 1 || c.TimeoutSeconds > 30 {
		return Config{}, errors.New("timeout must be between 1 and 30 seconds")
	}
	u, err := url.Parse(c.Endpoint)
	if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.Fragment != "" || !utf8.ValidString(c.Endpoint) || len(c.Endpoint) > 2048 {
		return Config{}, errors.New("HTTPS webhook endpoint without credentials or fragment required")
	}
	return c, nil
}

func decodeBinding(data []byte) (Binding, error) {
	fields, err := decodeUniqueObject(data)
	if err != nil || len(fields) != 5 {
		return Binding{}, errors.New("invalid extension binding")
	}
	var b Binding
	if err := decodeField(fields, "schema", &b.Schema); err != nil || b.Schema != "gotth.extension.runtime-binding.v1" {
		return Binding{}, errors.New("invalid extension binding schema")
	}
	if err := decodeField(fields, "instance_id", &b.InstanceID); err != nil || !validUUID(b.InstanceID) {
		return Binding{}, errors.New("valid instance binding required")
	}
	for name, target := range map[string]*string{"manifest_sha256": &b.ManifestSHA256, "grant_sha256": &b.GrantSHA256, "session_sha256": &b.SessionSHA256} {
		if err := decodeField(fields, name, target); err != nil || !validHex(*target, 64) {
			return Binding{}, errors.New("valid digest binding required")
		}
	}
	return b, nil
}

func decodeUniqueObject(data []byte) (map[string]json.RawMessage, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	token, err := decoder.Token()
	if err != nil || token != json.Delim('{') {
		return nil, errors.New("JSON object required")
	}
	fields := map[string]json.RawMessage{}
	for decoder.More() {
		nameToken, err := decoder.Token()
		name, ok := nameToken.(string)
		if err != nil || !ok {
			return nil, errors.New("invalid JSON object")
		}
		if _, exists := fields[name]; exists {
			return nil, errors.New("duplicate JSON field")
		}
		var value json.RawMessage
		if err := decoder.Decode(&value); err != nil {
			return nil, err
		}
		fields[name] = value
	}
	if token, err = decoder.Token(); err != nil || token != json.Delim('}') {
		return nil, errors.New("invalid JSON object")
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return nil, errors.New("trailing JSON input")
	}
	return fields, nil
}

func decodeField(fields map[string]json.RawMessage, name string, target any) error {
	data, ok := fields[name]
	if !ok {
		return errors.New("missing field")
	}
	return json.Unmarshal(data, target)
}

func readPrivateFile(path string, limit int64) ([]byte, error) {
	info, err := os.Lstat(path)
	if err != nil || info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() || info.Mode().Perm()&0o077 != 0 {
		return nil, errors.New("private regular file required")
	}
	handle, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer handle.Close()
	opened, err := handle.Stat()
	if err != nil || !os.SameFile(info, opened) || !opened.Mode().IsRegular() || opened.Mode().Perm()&0o077 != 0 {
		return nil, errors.New("private regular file changed during open")
	}
	data, err := io.ReadAll(io.LimitReader(handle, limit+1))
	if err != nil || int64(len(data)) > limit {
		return nil, errors.New("bounded private file required")
	}
	return data, nil
}

func validateSocketPath(path string) error {
	if !filepath.IsAbs(path) || len(path) > 100 {
		return errors.New("bounded absolute Unix socket path required")
	}
	parent, err := os.Lstat(filepath.Dir(path))
	if err != nil || parent.Mode()&os.ModeSymlink != 0 || !parent.IsDir() || parent.Mode().Perm()&0o077 != 0 {
		return errors.New("private Unix socket directory required")
	}
	if _, err := os.Lstat(path); err == nil || !errors.Is(err, os.ErrNotExist) {
		return errors.New("Unix socket path must not exist")
	}
	return nil
}

func validVersion(value string) bool {
	if value == "dev" || value == "1.0.0" {
		return true
	}
	for _, prefix := range []string{"1.0.0-alpha.", "1.0.0-beta."} {
		if strings.HasPrefix(value, prefix) {
			n := value[len(prefix):]
			return n != "" && n[0] != '0' && strings.Trim(n, "0123456789") == ""
		}
	}
	return false
}

func validToken(value string, min, max int) bool {
	if len(value) < min || len(value) > max || !utf8.ValidString(value) {
		return false
	}
	for _, r := range value {
		if r <= 0x20 || r == 0x7f {
			return false
		}
	}
	return true
}

func validHex(value string, length int) bool {
	if len(value) != length {
		return false
	}
	for _, r := range value {
		if !(r >= '0' && r <= '9' || r >= 'a' && r <= 'f') {
			return false
		}
	}
	return true
}

func validUUID(value string) bool {
	return len(value) == 36 && value[8] == '-' && value[13] == '-' && value[18] == '-' && value[23] == '-' && validHex(strings.ReplaceAll(value, "-", ""), 32)
}

func (c Config) Timeout() time.Duration { return time.Duration(c.TimeoutSeconds) * time.Second }
