# Verification

make verify runs normal and race tests, bounded fuzz targets, vet, build, format, module verification, published-control binding, release-script syntax, a local stable-closeout mock, real-process lifecycle, and the reproducible v1.0.0 bundle gate. The bundle gate checks immutable-HEAD export, sanitized offline Go configuration, replacement rejection, actual module-ZIP SHA-256, separately labelled Go h1 directory hashes, truthful deterministic SLSA fields, internal/external checksums, fixed identities/modes, atomic directory publication, invalid inputs, and non-replacement.

The full release matrix additionally runs 20 normal repeats, 20 shuffled race repeats, 20 hostile-input repeats, 20-second fuzz campaigns for both fuzz targets, govulncheck, and gitleaks on development. HTTPS behavior uses only isolated TLS receivers.

## Superseding v1.0.0 remediation — 2026-10-23

The earlier .artifacts/gotth-extension-webhook-v1-release-20261023 candidate is historical and not admission evidence: two cold reviews found pre-handler authentication, negotiated-authority, checksum-semantics, provenance, hermeticity, immutable-source, publication, and closeout defects.

The repaired worktree passed focused local Go tests, process lifecycle, stable closeout regression, script syntax, formatting, and diff checks. A sealed snapshot was committed inside an isolated development scratch clone so immutable-HEAD release machinery could be exercised before the release commit. The complete development matrix passed: source/diff/format/script checks, module and published-control binding, normal/race/repeat/shuffled-race/hostile suites, both fuzz campaigns, vet, build, process lifecycle, reproducible bundle, clean-mode bundle, govulncheck, and gitleaks. Exact historical logs, candidate identity, bundle, report, and seals are under .artifacts/gotth-extension-webhook-v1-remediation-20261023.

That isolated validation commit and its bundle remain provisional historical evidence. Later final review proved that the complete matrix was not revision-matched to the relative-output repair and that malformed UTF-8 could be rewritten before string validation. Those prior final-candidate directories are historical, not admission evidence. The shared object decoder now rejects malformed UTF-8 before any JSON token or value decoding, with raw `0xff` configuration and binding regressions. The complete matrix and exact clean bundle must be generated and sealed against the corrected commit; no older tree result may be carried forward.

## Remaining producer gates

- The corrected commit requires a revision-matched clean candidate, complete sealed matrix, exact clean bundle/provenance, and two fresh independent clean reviews of that identity and bundle.
- Release completion requires an annotated `v1.0.0` tag peeling to the reviewed commit and a successful published-mode canonical Forgejo closeout. Tag/Release publication remains separately authorized; the closeout implementation is regression-tested against a local Forgejo-compatible fixture.

Mail was not modified. Its A→B→A proof follows all four extension releases as consumer integration, using a normal Mail-owned exact pin with no patch, overlay, or bypass. GitHub publication is Bryce's handoff and does not block canonical Forgejo release.
