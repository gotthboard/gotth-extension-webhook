# Coverage map

| Subsystem | Required harnesses | Current evidence | Gaps |
| --- | --- | --- | --- |
| Configuration and secrets | unit, permission, symlink, duplicate/unknown input | `internal/webhook/config_test.go` | parent-directory swap by a same-UID hostile process is outside the owner-only runtime-directory boundary |
| Foundation control and RPC admission | unit and real Unix gRPC server/client | `internal/webhook/service_test.go`; `internal/webhook/grpc_admission_test.go` | malformed protobuf is handled by gRPC/protobuf runtime rather than hand-generated frames |
| Published V1 adoption | exact module/source/descriptor/RPC-path, no-local-copy, and authoritative session gates | `scripts/check-published-control.sh`; `internal/webhook/control_sdk_test.go`; negotiated-session regressions | final revision-matched heavy evidence pending |
| HTTPS alert delivery | TLS integration, HMAC, bounds, redirects, failures, ambiguity | `internal/webhook/service_test.go` | `internal/webhook/service_test.go` proves a receiver can consume the full body and close before response without a false delivery claim; full normal/race verification passed |
| Artifact and process | immutable-HEAD archive/SBOM/provenance/checksum reproducibility and local process lifecycle | `scripts/test-reproducible-artifact.sh`; `scripts/test-process.sh`; v1 release evidence | clean committed identity, two fresh clean reviews, tag, and Forgejo Release remain open; Mail integration follows all four extension releases |
| Stable supply chain | exact linked graph, source manifest, third-party licenses, SPDX SBOM, truthful SLSA provenance, and atomic bundle publication | `scripts/release-metadata.py`; `scripts/test-reproducible-artifact.sh`; `docs/dependencies.md` | local provenance is checksum-bound but not an additional cryptographic signature; signing is future work, not a v1.0.0 gate |
