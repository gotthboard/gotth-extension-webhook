# Published Extensions V1 adoption

State: `blocked`

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
Admission remains blocked: the unmodified Mail `4a83dee1` harness accepts a
candidate path but pins the different historical alpha.2 archive
`sha256:d3b79e577bc68aedcb83b8c3cc378341cf8610db514c751be3ff331df9a8937a`
and exposes no documented digest override. Historical alpha.1 evidence is
retained but does not prove the alpha.2 adoption binary. No tag, hosted Release,
deployment, release/API/provider/runtime credential use, or Mail source change
occurred. The branch sync used only existing configured Git
authentication.
