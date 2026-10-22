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

## Release status

The public GitHub mirror now exists, and its annotated
`v1.0.0-alpha.1` tag exactly matches Forgejo: tag object
`a126253f637ee9cf8b60482640f5f99c128a2d24` peels to the tested candidate
`3627c43213b82afa290372a6e2d9c36978a54758` on both remotes. Earlier records
that report no public mirror or no tag describe the historical state when they
were written.

Authenticated Forgejo has no Release object for the tag. No GitHub Release is
publicly visible; draft state is unknown without GitHub authentication. The tag
fallback URL and automatic source archives are not publication. Fresh exact-tag
verification on 2026-10-21 reproduced the
recorded archive byte-for-byte at SHA-256
`5cb6043ca200acfa67d4c6a85e0c1ba070c51dc550cacca7ce538021c8d9e83a`.
Later GOTTH Mail evidence also proves exact packaged delivery and real A→B→A
administrator rollback on the development test equipment; it does not claim
a live production change.

The closeout candidate adds the missing post-body transport-ambiguity
regression; the complete fresh normal/race/vet/build/reproducibility gate
passed on development. Two clean final reviews, Forgejo/GitHub prerelease publication,
and exact asset upload remain pending. Run
`scripts/check-release-closeout.sh` before and after publication; it verifies
remote tag identity and the expected Forgejo/GitHub hosted-release state without moving refs.
