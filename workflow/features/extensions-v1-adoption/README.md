# Published Extensions V1 adoption

State: `complete`

This alpha.2 candidate selects the published `gotth-extensions v1.0.0` and
`gotth-sdk v1.0.0` modules without replacement. Control messages and gRPC
services come from the SDK generated package; manifest, grant, and session
contracts remain foundation imports. The local control schema and generated
copy are removed.

Source, descriptor, full-method-path, module-selection, and no-local-copy gates
must pass with the functional, process, hostile-input, race/repeat/shuffle,
build, security-scan, and reproducible alpha.2 artifact gates.
All of those extension-repository gates passed locally or on the designated
development host on 2026-10-23 and are sealed in the requested evidence
bundle.

The implementation checkpoint is committed and pushed at `7ac45df3` (tree
`c8d62f9a`), and its reproducible local alpha.2 archive is
`sha256:0b09bd192e517a79d0088e9c639dd774a81d7089c26c655765f6095103e6a8a2`.
Published-V1 adoption is complete at the committed implementation checkpoint.
The unmodified Mail `4a83dee1` harness still pins a different historical alpha.2
archive, but that is post-release consumer-integration history rather than an
adoption or producer-release blocker. Historical alpha.1 evidence remains
unchanged. No alpha.2 tag, hosted Release, deployment, or Mail source change
occurred; the branch sync used only existing configured Git authentication.
