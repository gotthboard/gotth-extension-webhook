# Webhook extension evidence — 2026-09-20

## Candidate identity

- repository: `gotthboard/gotth-extension-webhook`
- commit: `95af997`
- extension ID: `gotth.mail.notification.webhook`
- interface: `gotth.mail.interface.notification`
- candidate version: `1.0.0-alpha.1`
- canonical manifest SHA-256:
  `0f515328c232f9c0db4d1e95be4877b7819815e5ade4cd6936ea823115d5c79e`

## Implemented boundary

The extension accepts only protected-file configuration, binding, service
token, and HMAC key inputs. It exposes an owner-only Unix gRPC socket, enforces
authenticated deadlines, provides foundation control and health services, and
implements the notification alert method only. Prompt methods remain
unimplemented.

Outbound delivery is HTTPS-only, refuses redirects and URL credentials,
serializes bounded deterministic JSON, authenticates the request with
HMAC-SHA-256, and discards response headers and bodies. It has no Mail
administrator, database, mailbox, prompt, approval, or mutation authority.

## Verification

`make verify` passed locally and on development host `10.0.0.97`, covering:

```text
go test ./...
go test -race ./...
go vet ./...
go build ./...
format checks
reproducible Linux/amd64 archive comparison
```

Hostile JSON, unknown/duplicate input, file ownership/mode/symlink rejection,
version grammar, TLS/HMAC delivery, redirects, non-2xx mapping, bounded network
failures, service authentication/deadlines, handshake, health, and manifest
generation are covered. GOTTH Mail's production-supervisor integration also
passes against the exact independently built process and manifest digest.

## Remaining release boundary

The canonical private Forgejo repository exists. Public GitHub creation,
one-way mirror configuration, immutable tag/artifact publication, release
parity, deployed rollback, and network ambiguity acceptance remain release
gates and are not represented as complete.
