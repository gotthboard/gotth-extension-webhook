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
