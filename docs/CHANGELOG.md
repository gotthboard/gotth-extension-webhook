# Changelog

## Unreleased

- Hardened release artifact construction: exact version grammar, commit-time
  epoch binding, optional clean-tree enforcement, linked binary-version
  verification, fixed archive modes, deterministic USTAR/gzip output,
  no-replace publication, and file/directory synchronization.
- Added regression proof for archive reproducibility, embedded checksums,
  linked version identity, epoch mismatch rejection, malformed-version
  rejection, and preservation of an existing artifact.
- Defined the bounded non-Telegram HTTPS webhook alert extension.
- Added protected file-based configuration, binding, service-token, and HMAC
  secret transport.
- Added authenticated Unix gRPC foundation control and alert-only notification
  RPC services.
- Added deterministic JSON/HMAC delivery, redirect rejection, bounded failure
  reporting, immutable version identity, and reproducible Linux/amd64 archives.
- Added cross-process GOTTH Mail supervisor lifecycle proof.
