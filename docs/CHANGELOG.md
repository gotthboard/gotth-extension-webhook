# Changelog

## Unreleased

- Prepared the independent `v1.0.0` stable candidate identity and exact release
  acceptance contract. The deterministic bundle now includes the archive, SPDX
  2.3 SBOM, in-toto/SLSA provenance, external checksums, tracked-source
  checksums, exact linked-module notices, and verbatim dependency licenses.
- Added pre-handler deadline/token admission for every RPC, normal SendPrompt
  admission, authoritative grant/host session derivation, negotiated capability
  responses, and SendAlert capability enforcement with real Unix-gRPC regression
  coverage.
- Reworked release output around an immutable HEAD export, sanitized offline Go
  environment, actual module-ZIP SHA-256 plus explicitly labelled Go h1 hashes,
  truthful deterministic SLSA fields, and one atomic bundle-directory rename.
  Added stable Forgejo-only closeout verification while retaining alpha history.
- Added stable compatibility, dependency/license, security, and post-release
  Mail-owned A→B→A upgrade/rollback requirements; bounded parser fuzzing; and
  rejection of actual Unicode line/paragraph separators at JSON/HTTP boundaries.
  The clean committed candidate and regenerated revision-matched evidence now
  pass, including ordinary relative-output regression coverage. Producer
  admission remains blocked on two fresh clean reviews and authorized Forgejo
  publication. GitHub publication is Bryce's handoff; Mail integration follows
  all four extension releases.
- Reject malformed UTF-8 bytes before any shared JSON object decoding, with
  raw-byte configuration and runtime-binding regressions. Release workflow
  completion now explicitly requires a revision-matched clean candidate and
  complete sealed matrix, two clean reviews of that identity and bundle, an
  annotated `v1.0.0` tag, and successful published-mode canonical Forgejo
  closeout.

- Sealed the published-V1 adoption implementation at commit `7ac45df3` and
  reproduced its alpha.2 archive at SHA-256 `0b09bd192e517a79d0088e9c639dd774a81d7089c26c655765f6095103e6a8a2`.
  Release admission remains blocked because the unmodified current Mail
  lifecycle harness pins a different historical alpha.2 candidate digest and
  exposes no digest override. No tag, hosted Release, or deployment was made.

- Prepared `v1.0.0-alpha.2`: select published `gotth-extensions v1.0.0` and
  `gotth-sdk v1.0.0` without replacements, import extension control from the
  SDK, and remove the repository-local control schema and generated Go copy.
  Exact module, source, descriptor, RPC-path, and no-local-copy gates preserve
  the V1 wire contract. The alpha.2 artifact and process gates replace only
  current candidate machinery; alpha.1 release records remain historical.

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
