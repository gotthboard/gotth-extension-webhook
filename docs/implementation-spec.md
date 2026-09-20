# GOTTH Webhook Extension Implementation Specification

## Runtime files

The process requires:

- `GOTTH_EXTENSION_CONFIG_FILE`: owner-only regular JSON file;
- `GOTTH_EXTENSION_BINDING_FILE`: owner-only regular JSON file containing the
  instance ID and exact manifest, grant, and negotiated-session digests;
- `GOTTH_EXTENSION_SECRET_DIR`: owner-only directory containing
  `webhook.hmac-key`;
- `GOTTH_EXTENSION_SERVICE_TOKEN_FILE`: owner-only regular token file;
- `GOTTH_EXTENSION_SOCKET`: absent Unix socket path below an owner-only runtime
  directory.

The version is linked into the executable at build time and exposed by
`gotth-extension-webhook --version`; runtime environment cannot rewrite it.

Symlinks, group/world-readable files, oversized files, unknown JSON fields,
duplicate JSON names, trailing JSON, embedded endpoint credentials,
non-HTTPS endpoints, fragments, and invalid timeouts are rejected.

Configuration schema:

```json
{
  "schema": "gotth.extension.webhook.config.v1",
  "endpoint": "https://receiver.example.test/gotth-alerts",
  "timeout_seconds": 10
}
```

## HTTP delivery

The body is deterministic JSON with schema
`gotth.extension.webhook.alert.v1`, the admitted alert fields, and sorted detail
keys. The maximum encoded body is 64 KiB. The process sends:

- `Content-Type: application/json`
- `X-GOTTH-Alert-ID: <bounded ID>`
- `X-GOTTH-Correlation-ID: <bounded correlation ID>`
- `X-GOTTH-Signature: sha256=<lowercase hex HMAC>`

The client permits TLS 1.2 or newer, disables redirects, bounds total request
time by the earlier RPC deadline or configured timeout, and does not retry.

## RPC admission

Every RPC requires a deadline, a bounded correlation ID, and constant-time
comparison of `x-gotth-mail-service-token` metadata against the protected token
file. Control handshake additionally binds the expected extension ID,
instance ID, manifest/grant/session digests, and challenge. Alerts require the
negotiated `notification.alert.send` capability.

`SendPrompt` always returns `Unimplemented`. No RPC accepts arbitrary bytes,
URLs, headers, executable expressions, or secret values.

## Artifact

The release archive contains only the static executable, `manifest.json`,
`configuration-metadata.json`, `LICENSE`, and `SHA256SUMS`. The archive name is
`gotth-extension-webhook-<version>-linux-amd64.tar.gz`. Artifact identity is
the SHA-256 digest of the final archive and is recorded as `sha256:<hex>`.
