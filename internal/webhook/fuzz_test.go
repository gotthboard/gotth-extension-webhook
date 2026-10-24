package webhook

import (
	"bytes"
	"encoding/json"
	"testing"

	pluginv1 "github.com/gotthboard/gotth-extension-webhook/proto/gotth/mail/plugin/v1"
)

func FuzzDecodeConfig(f *testing.F) {
	for _, seed := range [][]byte{
		[]byte("{\"schema\":\"gotth.extension.webhook.config.v1\",\"endpoint\":\"https://receiver.example.test/hook\",\"timeout_seconds\":5}"),
		[]byte("{\"schema\":\"gotth.extension.webhook.config.v1\",\"endpoint\":\"http://receiver.example.test\",\"timeout_seconds\":5}"),
		[]byte("{\"schema\":\"gotth.extension.webhook.config.v1\",\"schema\":\"duplicate\",\"endpoint\":\"https://receiver.example.test\",\"timeout_seconds\":5}"),
		[]byte("not-json"),
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		config, err := decodeConfig(data)
		if err != nil {
			return
		}
		if config.Schema != ConfigSchema || config.TimeoutSeconds < 1 || config.TimeoutSeconds > 30 {
			t.Fatalf("decoder accepted invalid config: %#v", config)
		}
		encoded, err := json.Marshal(config)
		if err != nil {
			t.Fatal(err)
		}
		roundTrip, err := decodeConfig(encoded)
		if err != nil || roundTrip != config {
			t.Fatalf("accepted config did not round trip: %#v %v", roundTrip, err)
		}
	})
}

func FuzzCanonicalAlert(f *testing.F) {
	for _, seed := range [][]byte{
		[]byte("alert-1"),
		[]byte("line one\nline two"),
		[]byte{0xff, 0xfe},
		[]byte("separator\u2028value"),
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		value := string(data)
		alert := &pluginv1.AlertMessage{
			Id: value, Class: "backup", Severity: "warning", Title: "Backup failed",
			Summary: "bounded summary", CorrelationId: "corr-1",
			Details: map[string]string{"input": value},
		}
		first, err := canonicalAlert(alert)
		if err != nil {
			return
		}
		second, err := canonicalAlert(alert)
		if err != nil || !bytes.Equal(first, second) {
			t.Fatalf("accepted alert serialization is not deterministic: %v", err)
		}
		if len(first) > maxBodyBytes || !json.Valid(first) {
			t.Fatalf("accepted alert produced invalid bounded JSON: %d bytes", len(first))
		}
	})
}
