# Verification

`make verify` runs unit, hostile-input, TLS/HMAC delivery, race, vet, build,
format, and reproducible-artifact checks. The GOTTH Mail integration gate also
builds this repository's exact candidate binary, computes its canonical
manifest digest, and proves real process start, authenticated handshake,
health, route admission, bounded delivery failure, route revocation, stop, and
runtime-secret cleanup through Mail's production supervisor.

Release evidence records exact commits, dependency pins, commands, artifact
digests, relevant coverage, explicit coverage gaps, mirror parity, rollback,
and review results. Green unit tests alone are not release evidence.
