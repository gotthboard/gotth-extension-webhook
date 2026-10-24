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

GOTTH Mail was not modified. No retained Mail record names this alpha.2 source
snapshot and artifact, so the exact remaining runtime gate is a fresh
production-supervisor lifecycle run against the final alpha.2 revision and
artifact. This blocks admission but does not invalidate the passed extension
repository gates.
