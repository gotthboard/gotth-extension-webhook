# Upgrade and rollback: alpha.2 to v1.0.0

GOTTH Mail owns this procedure. The extension must never update itself.

## Preconditions

1. Verify the clean-revision v1.0.0 bundle with external SHA256SUMS, SPDX SBOM, provenance subjects/materials, exact version identities, and the Mail-owned archive pin.
2. Retain the complete currently admitted artifact, pin, configuration, grant, session data, and encrypted HMAC key as rollback state.
3. Confirm the endpoint and schemas are unchanged. Recompute the canonical manifest digest because 1.0.0 replaces the alpha version identity.

## Upgrade A→B

1. Revoke routing to A before stopping it.
2. Stage B in a new digest-named, read-only artifact directory; never replace A.
3. Create a fresh owner-only runtime directory, service token, binding, and socket path. Keep secret values out of arguments and environment values.
4. Start B, then require authenticated handshake, exact manifest/grant/session binding, health, and a sandbox HTTPS/HMAC test before admitting routing.
5. On any failure, keep routing closed, stop B, preserve both artifacts and bounded audit evidence, and begin rollback.

## Rollback B→A

1. Revoke B routing, stop B, and remove B's ephemeral runtime secrets/socket.
2. Select A by its retained immutable digest; do not relabel B or mutate either artifact directory.
3. Recreate A's runtime state, handshake, and health. Restore routing only after those checks pass.
4. If A cannot be readmitted, leave routing closed and both artifacts available for diagnosis. Never report rollback success from process exit alone.

Revision-matched A→B→A proof through an unmodified Mail fixture is post-release consumer-integration acceptance after all four extension releases. The extension repository cannot satisfy it by changing or waiving Mail's immutable pin, and Forgejo publication does not wait on that later sequence.
