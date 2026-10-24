# GOTTH Webhook Extension

gotth-extension-webhook is an out-of-process GOTTH Mail notification extension. It delivers bounded operational alerts to one configured HTTPS endpoint and authenticates each request with HMAC-SHA-256.

The extension has no administrator, database, mailbox, prompt, approval, or mutation authority. GOTTH Mail owns installation, configuration, encrypted secret storage, process supervision, routing, audit, update, rollback, and release admission.

The candidate version is v1.0.0. It keeps gotth-extensions v1.0.0 and gotth-sdk v1.0.0 pinned without replacement and preserves their released V1 protocol contract. See docs/prd.md, docs/architecture.md, docs/implementation-spec.md, docs/compatibility.md, and docs/v1.0.0-acceptance.md.

## Candidate status

This branch prepares an untagged, unpublished v1.0.0 candidate. Its deterministic release machinery emits a Linux/amd64 archive, SPDX SBOM, SLSA provenance statement, external checksums, source checksums, and exact linked-module license material. Tests use only isolated TLS receivers; no production webhook is contacted.

The extension-owned stable gate requires a clean committed candidate, revision-matched release evidence, two fresh clean reviews, an annotated v1.0.0 tag, and an authorized non-draft Forgejo Release with the exact four-asset bundle. GitHub publication is Bryce's later distribution handoff and does not block canonical Forgejo release. After all four extension releases, Mail owns consumer integration and the unmodified A→B→A proof; the extension release neither changes nor bypasses Mail.

No v1.0.0 tag, hosted Release, upload, deployment, production delivery, or owner acceptance exists yet. Historical alpha.1 and alpha.2 records remain in their dated release notes and workflow evidence; they are not relabelled as stable proof.

## Build and verify

Run make verify for focused repository gates. Final release evidence additionally requires the development-host repeat/shuffle/fuzz/security matrix and the external lifecycle/review/publication gates listed in docs/v1.0.0-acceptance.md.
