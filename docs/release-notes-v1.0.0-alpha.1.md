# gotth-extension-webhook v1.0.0-alpha.1

First immutable Alpha prerelease of the out-of-process, HMAC-authenticated HTTPS
alert extension for GOTTH Mail.

## Exact identity

- Annotated tag object: `a126253f637ee9cf8b60482640f5f99c128a2d24`
- Peeled source commit: `3627c43213b82afa290372a6e2d9c36978a54758`
- Linux/amd64 archive:
  `gotth-extension-webhook-1.0.0-alpha.1-linux-amd64.tar.gz`
- Archive SHA-256:
  `5cb6043ca200acfa67d4c6a85e0c1ba070c51dc550cacca7ce538021c8d9e83a`
- Embedded executable SHA-256:
  `616d9fd3e99543f5dbfa92f526d51551b3883815881213fd7ba44338f4c83994`
- Canonical manifest SHA-256:
  `0f515328c232f9c0db4d1e95be4877b7819815e5ade4cd6936ea823115d5c79e`

Fresh Go 1.26.6 verification on 2026-10-21 passed normal, race, vet,
build, deterministic-artifact, archive-mode, checksum, and no-replace gates.
The original archive was reproduced byte-for-byte from the public tag. Later
GOTTH Mail evidence binds that exact archive to real packaged HTTPS/HMAC
delivery and A→B→A administrator rollback on a development test equipment.
The closeout successor adds a post-body transport-close regression proving the
extension returns an unavailable error and never claims delivery after the
receiver consumes the body but closes before a response.

## Boundary

The extension accepts protected-file configuration, one service token, and one
HMAC key; it exposes an owner-only Unix gRPC socket and sends bounded HTTPS-only
JSON without redirects or response-body disclosure. GOTTH Mail retains
installation, secrets, lifecycle, routing, retries, audit, and release authority.
This prerelease does not grant mailbox, database, administrator, prompt,
mutation, or Telegram authority. It is not a stable compatibility promise.
