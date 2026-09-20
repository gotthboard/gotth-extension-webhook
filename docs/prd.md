# GOTTH Webhook Extension PRD

## Problem

GOTTH Mail has an admitted notification RPC seam and a host-owned extension
administrator, but it lacks one independently packaged, non-Telegram extension
whose artifact, protected configuration, secret injection, supervision,
health, routing, update, and rollback can be proven end to end.

## Product boundary

The extension delivers one class of data: bounded operational alert messages
already admitted by GOTTH Mail. It sends those messages to one configured
HTTPS webhook endpoint and authenticates the exact transmitted body with an
HMAC-SHA-256 signature.

GOTTH Mail remains authoritative for installation, configuration, secret
storage, process lifecycle, routing, alert construction, redaction, retry
policy, audit, and release admission. The receiver remains authoritative for
accepting or rejecting the HTTP request.

## Requirements

- `WEBHOOK-001`: Run as a separate unprivileged process and communicate with
  GOTTH Mail over an owner-only Unix-domain socket.
- `WEBHOOK-002`: Implement GOTTH extension control v1 and the GOTTH Mail
  notification v1 alert method. Prompt and mutation methods fail closed.
- `WEBHOOK-003`: Accept exactly one HTTPS endpoint, a bounded timeout, one
  service-token file, and one HMAC-key file. Secret values never appear in
  manifests, arguments, environment values, logs, errors, health, or RPC
  responses.
- `WEBHOOK-004`: Serialize one deterministic bounded JSON request and sign the
  exact body with HMAC-SHA-256. Send the signature, alert ID, and correlation
  ID in fixed headers.
- `WEBHOOK-005`: Reject redirects, non-HTTPS endpoints, embedded credentials,
  invalid or oversized alerts, unsupported capabilities, missing deadlines,
  unauthenticated calls, and non-2xx responses.
- `WEBHOOK-006`: Return bounded delivery evidence without response bodies,
  headers, URLs, credentials, or receiver error details.
- `WEBHOOK-007`: Publish immutable Linux/amd64 artifacts with checksums,
  manifest, configuration metadata, and exact source/version identity.
- `WEBHOOK-008`: Prove test, enable, disable, update, rollback, redaction,
  crash recovery, and failure isolation through the GOTTH Mail administrator.

## Non-goals

- Telegram, chat commands, interactive prompts, or approval workflows.
- A generic event bus or arbitrary payload transport.
- Receiver provisioning or remote configuration.
- Mail database, mailbox, administrator, or Authentik access.
- Hidden retries or delivery-success claims after ambiguous transport failure.

## Acceptance

All changed behavior, edge cases, failure paths, and secret boundaries have
tests. Full, race, vet, build, artifact reproducibility, process integration,
and two fresh admission reviews pass before the first immutable release.
