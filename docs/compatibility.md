# Compatibility policy

v1.0.0 freezes the extension's supported V1 surface:

- GOTTH Extensions foundation v1.0.0 and Go SDK v1.0.0;
- control protocol gotth.extensions.control major 1, minor 0;
- notification interface gotth.mail.interface.notification major 1, minor 0;
- configuration schema gotth.extension.webhook.config.v1;
- runtime binding schema gotth.extension.runtime-binding.v1;
- alert body schema gotth.extension.webhook.alert.v1;
- secret slot webhook.hmac-key and the documented fixed HTTP headers.

Alpha.2 and v1.0.0 use the same protocol, interface, configuration, alert, and secret contracts. Stable v1 tightens the protected runtime-binding v1 document: the host must include the canonical grant and authoritative host profile in addition to their digests so the extension can derive and retain the complete negotiated session. Digest-only alpha.2 binding files are not accepted and must be regenerated during staging. The versions also have different manifest and artifact digests because the version is part of both identities. Upgrade must therefore stage and admit v1.0.0 as a new immutable artifact; relabelling an alpha archive is forbidden.

Patch releases on the v1.0.x line may fix defects without breaking these contracts. A change to RPC identity, schema meaning, required authority, or accepted configuration requires a separately reviewed compatibility decision; a breaking change requires a new major line. Compatibility does not promise receiver provisioning, retries, prompts, mutations, or arbitrary payloads.
