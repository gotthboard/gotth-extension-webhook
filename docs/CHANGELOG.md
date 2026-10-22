# Changelog

## Unreleased

- Recorded the superseding release-closeout state: the public GitHub mirror
  and exact `v1.0.0-alpha.1` tag parity now exist. Fresh exact-tag verification
  reproduced the original archive byte-for-byte, and later Mail evidence binds
  the same archive to real packaged delivery and A→B→A administrator rollback
  on the development test equipment. Added the missing post-body-close
  ambiguity regression and a closeout check for exact tag and hosted-release
  state. The complete fresh gate and two independent final reviews passed.
  Forgejo now publishes the exact prerelease assets; matching GitHub publication
  and cross-host parity remain open. The webhook transport has explicit
  redirect fail-closed regression coverage.
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
