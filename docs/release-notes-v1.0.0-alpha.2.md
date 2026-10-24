# gotth-extension-webhook v1.0.0-alpha.2

Alpha.2 adopts the published GOTTH Extensions V1 foundation and Go SDK while
preserving the existing protobuf wire contract and webhook behavior.

## Changes from alpha.1

- Select `github.com/gotthboard/gotth-extensions v1.0.0` and
  `github.com/gotthboard/gotth-sdk v1.0.0` without a `replace` directive.
- Import `gotth.extensions.v1.ExtensionControl` generated messages and gRPC
  services from `github.com/gotthboard/gotth-sdk/gen/go/gotth/extensions/v1`.
- Remove the repository-local control schema and generated Go files.
- Retain foundation imports for manifest construction, validation, canonical
  digest binding, grant, and session contracts.
- Add exact published-module, source, descriptor, RPC-path, and no-local-copy
  regression gates plus a bounded real-process lifecycle gate.
- Move reproducible release-artifact verification to `1.0.0-alpha.2`.

## Exact upstream binding

- `gotth-extensions v1.0.0` module sum:
  `h1:lKXOydN/7+Mq9W2juqnWUXNJiirpYR6F3KPXYV55Hho=`
- `gotth-sdk v1.0.0` module sum:
  `h1:RmjRs/l+H3KMDmnimZs2OL8ExCchT/opVO5tDdjb4EA=`
- authoritative control source SHA-256:
  `846dbdd487e9688a338c7885478e1296456f994813330ecccdc3f4b225f10cfe`
- deterministic descriptor-set SHA-256:
  `62a549b29a030f0decb46b2b1c2f382bbbc75430ba0b3d7c729bd21528b505cf`

## Admission status

These are candidate notes, not a publication claim. The final alpha.2 commit,
tree, source archive, executable, manifest, release archive, and checksums do
not exist until an authorized commit and final clean build. No alpha.2 tag,
hosted Release, deployment, or owner acceptance exists. Revision-matched GOTTH
Mail lifecycle evidence and final admission reviews remain required. Alpha.1
release notes and immutable release history are unchanged.

The sealed uncommitted snapshot passed bounded local verification and the
designated development-host full, race, repeat, shuffled-race, hostile-input,
vet, build, process, reproducible two-build, vulnerability, and secret-scan
gates. That evidence does not replace the missing revision-matched GOTTH Mail
production-supervisor lifecycle gate.
