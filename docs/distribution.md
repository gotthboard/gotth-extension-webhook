# Distribution and rollback

Forgejo `gotthboard/gotth-extension-webhook` is canonical. GitHub
`gotthboard/gotth-extension-webhook` is the public one-way distribution mirror.
Refs, annotated tags, release artifacts, and checksums must match exactly.

Every release archive is immutable and named
`gotth-extension-webhook-<version>-linux-amd64.tar.gz`. GOTTH Mail pins the
archive SHA-256 digest and installs it into a digest-named read-only artifact
directory. It never executes a floating branch, mutable `latest` name, or
in-place replacement.

The builder requires the exact release grammar, binds archive timestamps to
the source commit epoch, verifies the linked binary version, emits deterministic
USTAR plus timestamp-free gzip, and publishes with no-replace semantics plus
file and directory synchronization. Final release builds set
`GOTTH_EXTENSION_REQUIRE_CLEAN=1`; an existing artifact is never overwritten.

Update retains the complete preceding artifact pin and configuration contract.
Rollback first revokes routing and stops the current process, selects the
preceding immutable artifact, repeats handshake and health, and only then
restores routing. A failed rollback leaves routing closed and both artifacts
available for diagnosis.

No tag or release is admissible without reproducible archive proof, exact
binary/tag identity, Forgejo/GitHub parity, checksums, full/race/vet/build,
GOTTH Mail lifecycle integration, and two fresh clean reviews.

The next candidate version is `v1.0.0-alpha.2`. Its reproducibility gate and
release notes are alpha.2-specific. The sealed implementation checkpoint at
`7ac45df3` reproducibly builds a local archive with SHA-256
`0b09bd192e517a79d0088e9c639dd774a81d7089c26c655765f6095103e6a8a2`.
That local evidence is not publication: no alpha.2 tag, hosted Release,
published immutable archive, deployment, or owner acceptance exists. The
closeout script below remains deliberately bound to the already-published
alpha.1 identity; changing its immutable hashes would rewrite release history
rather than prepare alpha.2.

## Release closeout check

The following closeout status and command are the retained alpha.1 record.

The current state is exact tag parity plus a published Forgejo prerelease with
exact assets. No GitHub Release is publicly visible; GitHub draft state is
unknown without authentication:

```sh
FORGEJO_TOKEN=... EXPECT_STATE=incomplete sh scripts/check-release-closeout.sh
```

After authorized GitHub publication, require matching non-draft prereleases and
exact assets on both hosts:

```sh
FORGEJO_TOKEN=... EXPECT_STATE=complete sh scripts/check-release-closeout.sh
```

Supply credentials through protected operator environments. The check binds
canonical HTTP origins, validates Release metadata, downloads authenticated
Forgejo assets only from that origin, verifies exact checksums, and compares
cross-host bytes. It never moves a tag.
