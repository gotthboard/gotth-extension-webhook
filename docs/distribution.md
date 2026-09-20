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
