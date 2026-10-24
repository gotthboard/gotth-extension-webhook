# Coverage map

| Subsystem | Required harnesses | Current evidence | Gaps |
| --- | --- | --- | --- |
| Configuration and secrets | unit, permission, symlink, duplicate/unknown input | `internal/webhook/config_test.go` | parent-directory swap by a same-UID hostile process is outside the owner-only runtime-directory boundary |
| Foundation control | unit and real Unix gRPC process | `internal/webhook/service_test.go`; GOTTH Mail `internal/extensionsruntime` cross-process test | malformed protobuf is handled by gRPC/protobuf runtime rather than hand-generated frames |
| Published V1 adoption | exact module/source/descriptor/RPC-path and no-local-copy gates | `scripts/check-published-control.sh`; `internal/webhook/control_sdk_test.go`; exact `7ac45df3` development-host evidence | Mail `4a83dee1` pins a different alpha.2 digest and has no documented override; revision-matched lifecycle proof remains blocked |
| HTTPS alert delivery | TLS integration, HMAC, bounds, redirects, failures, ambiguity | `internal/webhook/service_test.go` | `internal/webhook/service_test.go` proves a receiver can consume the full body and close before response without a false delivery claim; full normal/race verification passed |
| Artifact and supervisor | alpha.2 reproducibility, local process lifecycle, and GOTTH Mail administrator lifecycle | `scripts/test-reproducible-artifact.sh`; `scripts/test-process.sh`; exact archive `0b09bd19…`; historical GOTTH Mail `internal/extensionsruntime` alpha.1 evidence | alpha.2 revision-matched Mail lifecycle is blocked on the Mail-owned digest pin; reviews, tag, hosted releases, and cross-host asset parity remain open |
