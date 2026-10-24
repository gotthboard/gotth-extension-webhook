# GOTTH Webhook Extension Architecture

## Boundary

```text
GOTTH Mail
  host-owned administrator
    encrypted configuration and HMAC key
    artifact pin and lifecycle state
    owner-only runtime directory
           |
           | Unix gRPC + per-instance service token
           v
gotth-extension-webhook
  ExtensionControl v1
  NotificationBackend v1 (alerts only)
           |
           | HTTPS POST + HMAC-SHA-256
           v
configured webhook receiver
```

The extension-control service implementation imports the generated V1 package
from `gotth-sdk v1.0.0`. Foundation manifest construction and validation stay
in `gotth-extensions v1.0.0`; the extension repository does not own or vendor a
second control schema. This changes source ownership, not protobuf package
names, field numbers, service names, or full RPC paths.

The process receives paths, not secret values, in its environment. The Mail
supervisor creates an owner-only instance directory and writes the admitted JSON
configuration, canonical grant, authoritative host profile, binding digests, and
secret files with mode `0600`. The extension reruns foundation negotiation and
retains the derived session before it opens the socket. The supervisor then
starts the digest-pinned artifact with a bounded environment, waits for the Unix
socket, performs the foundation handshake and health probe, and only then admits
routing.

Disable revokes routing before process termination. Test uses a separate
ephemeral runtime directory and always stops the process. Update and rollback
select complete immutable artifact directories; they never replace a running
binary in place.

## Trust model

The extension is untrusted relative to Mail. It receives one scoped service
token, one endpoint, one HMAC key, and already-redacted bounded alerts. It does
not receive the Mail database, administrator session, OIDC/SCIM credentials,
mailbox contents, Docker socket, host network namespace control, or arbitrary
filesystem access from the supervisor.

The webhook receiver is remote and untrusted. Redirects are disabled. Response
bodies and headers are discarded. A 2xx status proves only receiver admission;
connection loss after body transmission remains an ambiguous failure and is
reported as such.

## Failure model

- Invalid local files or permissions prevent startup.
- Failed handshake or health prevents routing.
- Non-2xx HTTP status is a permanent rejected delivery.
- Timeout, TLS, DNS, connection, or ambiguous transport errors are reported as
  failed without fabricated delivery evidence.
- Process exit removes health and prevents new routing; Mail records the
  failure and retains the prior immutable artifact for rollback.
