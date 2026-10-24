# Coverage map

| Subsystem | Required harnesses | Current evidence | Gaps |
| --- | --- | --- | --- |
| Configuration and secrets | unit, permission, symlink, duplicate/unknown input | `internal/webhook/config_test.go` | parent-directory swap by a same-UID hostile process is outside the owner-only runtime-directory boundary |
| Foundation control | unit and real Unix gRPC process | `internal/webhook/service_test.go`; GOTTH Mail `internal/extensionsruntime` cross-process test | malformed protobuf is handled by gRPC/protobuf runtime rather than hand-generated frames |
| Published V1 adoption | exact module/source/descriptor/RPC-path and no-local-copy gates | `scripts/check-published-control.sh`; `internal/webhook/control_sdk_test.go` | final revision-matched GOTTH Mail lifecycle proof remains open |
| HTTPS alert delivery | TLS integration, HMAC, bounds, redirects, failures, ambiguity | `internal/webhook/service_test.go` | `internal/webhook/service_test.go` proves a receiver can consume the full body and close before response without a false delivery claim; full normal/race verification passed |
| Artifact and supervisor | alpha.2 reproducibility, local process lifecycle, and GOTTH Mail administrator lifecycle | `scripts/test-reproducible-artifact.sh`; `scripts/test-process.sh`; historical GOTTH Mail `internal/extensionsruntime` alpha.1 evidence | alpha.2 revision-matched Mail lifecycle, immutable candidate identity, reviews, tag, hosted releases, and cross-host asset parity remain open |
