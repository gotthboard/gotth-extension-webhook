# Security

Report vulnerabilities privately to the maintainers of the `gotthboard`
organization. Do not include production webhook URLs, HMAC keys, service
tokens, alert bodies, or deployment credentials in public issues.

The extension accepts only protected local runtime files and authenticated
Unix-socket RPC. It has no administrator, database, mailbox, OIDC, SCIM,
Docker-socket, prompt, or mutation authority. Suspected secret disclosure,
signature bypass, redirect following, arbitrary request construction,
artifact substitution, dependency/license drift, SBOM/provenance mismatch,
Unicode header-boundary bypass, or lifecycle escape is release-blocking.

The v1.0.x line becomes the supported stable line only after v1.0.0 is
published. Alpha artifacts remain historical prereleases and receive no stable
compatibility promise. Security reports must identify the affected exact
version or archive digest without including secrets or production payloads.
