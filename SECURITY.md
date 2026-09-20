# Security

Report vulnerabilities privately to the maintainers of the `gotthboard`
organization. Do not include production webhook URLs, HMAC keys, service
tokens, alert bodies, or deployment credentials in public issues.

The extension accepts only protected local runtime files and authenticated
Unix-socket RPC. It has no administrator, database, mailbox, OIDC, SCIM,
Docker-socket, prompt, or mutation authority. Suspected secret disclosure,
signature bypass, redirect following, arbitrary request construction,
artifact substitution, or lifecycle escape is release-blocking.
