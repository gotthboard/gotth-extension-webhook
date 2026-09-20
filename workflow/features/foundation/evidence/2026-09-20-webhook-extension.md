# Webhook extension evidence — 2026-09-20

## Candidate identity

- repository: `gotthboard/gotth-extension-webhook`
- commit: `3627c43213b82afa290372a6e2d9c36978a54758`
- extension ID: `gotth.mail.notification.webhook`
- interface: `gotth.mail.interface.notification`
- candidate version: `1.0.0-alpha.1`
- canonical manifest SHA-256:
  `0f515328c232f9c0db4d1e95be4877b7819815e5ade4cd6936ea823115d5c79e`
- release archive SHA-256:
  `5cb6043ca200acfa67d4c6a85e0c1ba070c51dc550cacca7ce538021c8d9e83a`
- linked executable SHA-256:
  `616d9fd3e99543f5dbfa92f526d51551b3883815881213fd7ba44338f4c83994`
- source epoch: `1789940202`

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

The final clean source produced the release archive twice with identical
bytes. The archive has fixed `0755`/`0644` member modes, deterministic USTAR
and timestamp-free gzip, commit-epoch timestamps, a linked
`1.0.0-alpha.1` binary identity, and valid embedded checksums. Regression
tests reject an invalid version, mismatched epoch, and replacement of an
existing artifact. The final build also exercised mandatory clean-tree mode.
The final archive passed GOTTH Mail's
real start, authenticated handshake, health, route admission, bounded delivery
failure, route revocation, stop, and runtime-secret cleanup lifecycle.

The focused implementation package reports 72.4% statement coverage. The
uncovered command/process boundary is exercised through the real independent
binary lifecycle; defensive OS failures without an injectable standard-library
seam remain the explicit gap. No fake hook was added to inflate the number.

Twenty shuffled race repetitions passed. Two fresh code/artifact reviews are
clean. The admission review still denies release publication because public
distribution parity is unavailable.

## Remaining release boundary

Forgejo `main` contains the exact candidate and later evidence-only commits;
the eventual immutable release tag must resolve to the tested candidate
`3627c43213b82afa290372a6e2d9c36978a54758`. Neither the development host nor
the active agent host has GitHub CLI or SSH authentication, and unauthenticated
access confirms no usable public mirror. Public GitHub creation, one-way
mirror configuration, immutable tag/artifact publication, release parity,
deployed rollback, and network ambiguity acceptance remain release gates and
are not represented as complete. No tag or release was created.
