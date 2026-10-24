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

These are candidate notes, not a publication claim. The adoption implementation
checkpoint is commit `7ac45df3d9712978868f15d14bc8c402e0d1b0ed`, tree
`c8d62f9ad598409071dce76bdfc89ae2d9a110b3`. Its reproducible local
Linux/amd64 archive SHA-256 is
`0b09bd192e517a79d0088e9c639dd774a81d7089c26c655765f6095103e6a8a2`.
No alpha.2 tag, hosted Release, deployment, or owner acceptance exists.
Alpha.1 release notes and immutable release history are unchanged.

The exact checkpoint passed bounded local verification and the designated
development-host full, race, repeat, shuffled-race, hostile-input, vet, build,
process, reproducible two-build, vulnerability, and secret-scan gates.
Revision-matched Mail lifecycle admission is blocked rather than passed: Mail
commit `4a83dee139f3735a4b2f925c54c339cc1afd19df` pins alpha.2 archive
`d3b79e577bc68aedcb83b8c3cc378341cf8610db514c751be3ff331df9a8937a`.
The documented fourth harness argument supplies only the candidate path; no
documented input overrides that compiled pin. The exact adoption archive was
therefore not run through a modified or bypassed harness.
