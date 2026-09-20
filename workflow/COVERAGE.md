# Coverage map

| Subsystem | Required harnesses | Current evidence | Gaps |
| --- | --- | --- | --- |
| Configuration and secrets | unit, permission, symlink, duplicate/unknown input | `internal/webhook/config_test.go` | parent-directory swap by a same-UID hostile process is outside the owner-only runtime-directory boundary |
| Foundation control | unit and real Unix gRPC process | `internal/webhook/service_test.go`; GOTTH Mail `internal/extensionsruntime` cross-process test | malformed protobuf is handled by gRPC/protobuf runtime rather than hand-generated frames |
| HTTPS alert delivery | TLS integration, HMAC, bounds, redirects, failures, ambiguity | `internal/webhook/service_test.go` | network-injected post-body TCP ambiguity needs the deployed acceptance harness |
| Artifact and supervisor | reproducibility and GOTTH Mail administrator lifecycle | `scripts/test-reproducible-artifact.sh`; GOTTH Mail `internal/extensionsruntime` | Forgejo/GitHub release parity and deployed rollback remain release gates |
