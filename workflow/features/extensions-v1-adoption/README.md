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

Admission remains blocked until the final extension revision and artifact have
revision-matched GOTTH Mail production-supervisor lifecycle evidence. Historical
alpha.1 evidence is retained but does not prove the alpha.2 binary. No commit,
tag, hosted Release, deployment, or Mail source change is authorized here.
