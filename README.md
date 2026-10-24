# GOTTH Webhook Extension

`gotth-extension-webhook` is an out-of-process GOTTH Mail notification
extension. It delivers bounded operational alerts to one configured HTTPS
endpoint and authenticates each request with an HMAC-SHA-256 signature.

The extension has no administrator, database, mailbox, prompt, approval, or
mutation authority. GOTTH Mail owns installation, configuration, encrypted
secret storage, process supervision, routing, audit, update, and rollback.

The first release is being developed on the `1.0.0-alpha.N` line. See
[`docs/prd.md`](docs/prd.md), [`docs/architecture.md`](docs/architecture.md),
and [`docs/implementation-spec.md`](docs/implementation-spec.md).

The next extension version is `v1.0.0-alpha.2`. Its candidate replaces the
repository-local extension-control protobuf with the published
`github.com/gotthboard/gotth-sdk/gen/go/gotth/extensions/v1` package. The
foundation manifest, grant, and session contracts remain imported from
`github.com/gotthboard/gotth-extensions/pkg/extensions`. Both modules are
selected at `v1.0.0` without a `replace` directive.

## Release status

The alpha.2 adoption implementation is committed and pushed on
`workflow/feature/extensions-v1-adoption` at
`7ac45df3d9712978868f15d14bc8c402e0d1b0ed` (tree
`c8d62f9ad598409071dce76bdfc89ae2d9a110b3`). It has not been tagged,
published, deployed, or accepted by GOTTH Mail. Bounded local verification and
the designated development-host full/race/repeat/shuffle/vet/build/process/
hostile/reproducibility/govulncheck/gitleaks matrix pass for that exact
checkpoint. Its reproducible local alpha.2 archive has SHA-256
`0b09bd192e517a79d0088e9c639dd774a81d7089c26c655765f6095103e6a8a2`.

Release admission remains blocked. The unmodified Mail harness at
`4a83dee139f3735a4b2f925c54c339cc1afd19df` accepts an alpha.2 path as its
fourth runner argument but admits only the source-pinned archive SHA-256
`d3b79e577bc68aedcb83b8c3cc378341cf8610db514c751be3ff331df9a8937a`.
It has no documented digest override, so it cannot exercise the exact adoption
archive without a Mail-owned pin change. The lifecycle was not run against a
bypassed harness. Alpha.1 evidence below remains historical and is not
relabelled as alpha.2 proof.

### Historical alpha.1 status

The public GitHub mirror now exists, and its annotated
`v1.0.0-alpha.1` tag exactly matches Forgejo: tag object
`a126253f637ee9cf8b60482640f5f99c128a2d24` peels to the tested candidate
`3627c43213b82afa290372a6e2d9c36978a54758` on both remotes. Earlier records
that report no public mirror or no tag describe the historical state when they
were written.

Forgejo now publishes the non-draft prerelease with the exact archive and
`SHA256SUMS`. No GitHub Release is publicly visible; draft state is unknown
without GitHub authentication. The tag
fallback URL and automatic source archives are not publication. Fresh exact-tag
verification on 2026-10-21 reproduced the
recorded archive byte-for-byte at SHA-256
`5cb6043ca200acfa67d4c6a85e0c1ba070c51dc550cacca7ce538021c8d9e83a`.
Later GOTTH Mail evidence also proves exact packaged delivery and real A→B→A
administrator rollback on the development test equipment; it does not claim
a live production change.

The closeout candidate adds the missing post-body transport-ambiguity
regression; the complete fresh normal/race/vet/build/reproducibility gate
passed on development, including redirect and post-body-close regressions, and
two independent reviews are clean. Forgejo publication is complete; matching
GitHub prerelease publication and cross-host byte parity remain pending. Run
`scripts/check-release-closeout.sh` before and after publication; it verifies
remote tag identity and the expected Forgejo/GitHub hosted-release state without moving refs.
