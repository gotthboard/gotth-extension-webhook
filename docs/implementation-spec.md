# GOTTH Webhook Extension Implementation Specification

## Runtime files

The process requires:

- `GOTTH_EXTENSION_CONFIG_FILE`: owner-only regular JSON file;
- `GOTTH_EXTENSION_BINDING_FILE`: owner-only regular JSON file containing the
  instance ID, canonical grant, authoritative host protocol/interface ranges,
  and exact manifest, grant, and negotiated-session digests;
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

Control request/response and gRPC service types come only from
`github.com/gotthboard/gotth-sdk/gen/go/gotth/extensions/v1` at module version
`v1.0.0`. The module's copied source is bound byte-for-byte to the
`gotth-extensions v1.0.0` control source at SHA-256
`846dbdd487e9688a338c7885478e1296456f994813330ecccdc3f4b225f10cfe`;
the deterministic descriptor-set SHA-256 is
`62a549b29a030f0decb46b2b1c2f382bbbc75430ba0b3d7c729bd21528b505cf`.


A unary interceptor requires a deadline and constant-time comparison of
`x-gotth-mail-service-token` metadata before dispatching every RPC handler. Each
handler then requires its bounded correlation ID and method-specific request
admission. Startup parses the canonical grant and authoritative host profile,
reruns gotth-extensions v1 negotiation, compares instance/manifest/grant/session
identity, and retains the complete session. Handshake and legacy capability
discovery return only that session; alerts require its negotiated
`notification.alert.send` capability.

`SendPrompt` performs normal interceptor and correlation-ID admission before it
returns `Unimplemented`. No RPC accepts arbitrary bytes,
URLs, headers, executable expressions, or secret values.

## Artifact

The stable archive contains the static executable, canonical manifest,
configuration metadata, project license, security policy, tracked-source
checksums, SPDX 2.3 SBOM, third-party notices, verbatim linked-module license
files, and internal `SHA256SUMS`. The release bundle adds the archive-named SPDX
SBOM, an in-toto Statement v1 / SLSA provenance v1 JSONL file, and external
`SHA256SUMS`. Artifact identity is the SHA-256 digest of the final archive.

The v1.0.0 gate builds the complete bundle twice and requires identical bytes,
exact linked-module/license allowlist, embedded and external checksums, linked
version identity, fixed modes, source-epoch binding, source commit/tree and
manifest provenance, malformed-version rejection, clean-tree final mode, and
same-filesystem atomic directory publication to a previously absent output path.
The build runs only from an immutable exported HEAD, rejects ambient Go
workspace/config/replace influence, and records module ZIP SHA-256 separately
from Go h1 directory hashes. Optional invocation-specific SLSA metadata is
omitted rather than fabricated.
