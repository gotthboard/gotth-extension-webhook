# Verification

`make verify` runs unit, hostile-input, TLS/HMAC delivery, race, vet, build,
format, published-control binding, real process, and reproducible alpha.2
artifact checks. `scripts/check-published-control.sh` requires exact
`gotth-extensions v1.0.0` and `gotth-sdk v1.0.0` selection without replacement,
compares their authoritative control sources, checks the SDK source binding and
foundation descriptor pin, and rejects any repository-local control copy. A Go
regression hashes the descriptor actually linked from the SDK and pins the two
full RPC paths. The GOTTH Mail integration gate also
builds this repository's exact candidate binary, computes its canonical
manifest digest, and proves real process start, authenticated handshake,
health, route admission, bounded delivery failure, route revocation, stop, and
runtime-secret cleanup through Mail's production supervisor.

Release evidence records exact commits, dependency pins, commands, artifact
digests, relevant coverage, explicit coverage gaps, mirror parity, rollback,
and review results. Green unit tests alone are not release evidence.

Alpha.1 Mail lifecycle records are retained history. They may be reused for
unchanged lifecycle mechanics only when their Mail and extension revisions and
artifact identity match the alpha.2 candidate; otherwise the remaining gate is
a fresh Mail production-supervisor lifecycle run. Mail source is outside this
adoption patch.

## Alpha.2 adoption result — 2026-10-23

Bounded local `GOMAXPROCS=2 make verify` passed. On the designated development
host, fresh full and race tests, 20 normal repeats, 20 shuffled race repeats,
20 focused hostile-input repeats, vet, build, real-process lifecycle,
reproducible two-build artifact comparison, `govulncheck`, and `gitleaks` all
passed. `govulncheck` reported no vulnerabilities and `gitleaks` reported no
leaks. The sealed logs, complete binary patch, source-file checksums, and
candidate Git tree are under
`.artifacts/gotth-extension-webhook-v1-adoption-20261023` in the parent
workspace.

The verified implementation checkpoint is
`7ac45df3d9712978868f15d14bc8c402e0d1b0ed`, tree
`c8d62f9ad598409071dce76bdfc89ae2d9a110b3`. Its two independently built
archives match at SHA-256
`0b09bd192e517a79d0088e9c639dd774a81d7089c26c655765f6095103e6a8a2`.

GOTTH Mail was not modified. The unmodified Mail harness at commit
`4a83dee139f3735a4b2f925c54c339cc1afd19df`, tree
`d85c22620f6e2d72b5fcfb49f34ffd16de0d583e`, documents an optional fourth
argument for an alpha.2 archive path. Its loader nevertheless admits only the
compiled SHA-256
`d3b79e577bc68aedcb83b8c3cc378341cf8610db514c751be3ff331df9a8937a`;
there is no documented flag or environment input for an alternate digest.
Passing the exact adoption archive would fail the loader's digest check before
staging B. Changing the constant, using an overlay, or bypassing validation
would not be revision-matched Mail proof, so the lifecycle gate was not run.
This exact blocker keeps alpha.2 unpublished without invalidating the passed
extension-repository gates.
