# Cold review 2

## Decision

Code and artifact: CLEAN. Release publication: BLOCKED.

## Independent pass

The final candidate was reread across protected file loading, strict JSON,
manifest/grant/session binding, service authentication, capability surface,
bounded deterministic alert encoding, HTTPS/HMAC transport, redirect and
response-data rejection, process/socket lifecycle, version identity, archive
construction, checksums, reproducibility, and Mail supervisor integration.

The exact archive is reproducible from a clean source commit, has fixed modes
and timestamps, contains only the five declared members, reports the linked
version, and passes all embedded checksums. Full, race, vet, build, shuffled
repeat, and real cross-process lifecycle gates pass. No runtime or artifact
defect remains.

The distribution contract still requires a public GitHub mirror plus exact
tag and release parity. No GitHub credential is available on either execution
host and the public repository cannot be created or populated anonymously.
Creating only the Forgejo tag would manufacture a split release. Rejected.
