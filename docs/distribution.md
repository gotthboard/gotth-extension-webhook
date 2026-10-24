# Distribution and rollback

Forgejo gotthboard/gotth-extension-webhook is canonical. The v1.0.0 producer release is complete when the reviewed annotated tag and non-draft stable Forgejo Release expose the exact verified bundle. GitHub publication is Bryce's later distribution handoff; missing GitHub parity does not block canonical Forgejo publication.

Every release archive is immutable and named gotth-extension-webhook-<version>-linux-amd64.tar.gz. Consumers pin its SHA-256 and install it into a digest-named read-only directory. They never execute a floating branch, mutable latest name, or in-place replacement.

## Reproducible release build

scripts/build-artifact.sh resolves one HEAD commit, exports it through git archive, then executes the exported builder against exported source. The mutable checkout is never compiled or inventoried. Ambient GOENV, GOWORK, GOFLAGS, GOEXPERIMENT, GOFIPS140, workspace replacements, automatic toolchain selection, and network module resolution are disabled; readonly module mode, local toolchain, GOAMD64=v1, and the verified module cache are required.

The bundle contains:

- gotth-extension-webhook-1.0.0-linux-amd64.tar.gz;
- gotth-extension-webhook-1.0.0-linux-amd64.spdx.json;
- gotth-extension-webhook-1.0.0-linux-amd64.intoto.jsonl;
- SHA256SUMS.

SPDX SHA256 values hash the module ZIP named by downloadLocation. Go h1 directory hashes remain explicitly labelled h1 values in package comments and provenance dirHash1 entries. Deterministic provenance omits invocationId, startedOn, and finishedOn rather than inventing execution facts.

The builder resolves an ordinary relative output path against the caller's working directory before crossing into the immutable source export. It creates and synchronizes all four files in a hidden same-filesystem staging directory, then exposes the complete bundle with one directory rename to a previously absent output path. It refuses an existing target and uses a cooperative lock. Pre-rename failure may leave hidden staging/lock state but no output. Failure after rename may leave a complete but unaccepted output and/or lock. Inspect checksums and state before cleanup; command failure never implies acceptance.

Final release builds set GOTTH_EXTENSION_REQUIRE_CLEAN=1. The source commit epoch is the only accepted SOURCE_DATE_EPOCH.

## Producer release acceptance

Producer acceptance is defined by docs/v1.0.0-acceptance.md: exact released foundation/SDK pins, authority and admission tests, full verification/security matrix, deterministic immutable-source bundle, two fresh clean reviews, clean commit identity, annotated tag, and stable Forgejo Release. Publication remains a separately authorized action.

Mail A→B→A is post-release consumer integration after all four extension releases. This repository does not edit or bypass Mail to make publication possible. GitHub parity is likewise a post-Forgejo handoff owned by Bryce.

## Stable Forgejo closeout

Before publication, authenticated unpublished-state verification is:

    FORGEJO_TOKEN=... EXPECT_STATE=unpublished sh scripts/check-v1-release-closeout.sh

After authorized tag and Release publication, bind the exact identities and local bundle:

    FORGEJO_TOKEN=... EXPECT_STATE=published \
      EXPECTED_TAG_OBJECT=<annotated-tag-object> EXPECTED_COMMIT=<peeled-commit> \
      BUNDLE_DIR=<final-bundle-directory> sh scripts/check-v1-release-closeout.sh

The stable check validates the canonical Forgejo repository, annotated tag identity, stable non-draft Release metadata, exact four-asset set, authenticated download origins, checksums, and byte equality with the final local bundle. It never moves a tag or writes remote state.

scripts/check-release-closeout.sh remains the immutable alpha.1 cross-host historical check. Its hard-coded alpha identity is not reused for stable v1.0.0.

## Consumer rollback

Mail retains the preceding artifact pin and configuration contract. Rollback revokes routing, stops the current process, selects the preceding immutable artifact, repeats authoritative binding derivation, handshake, and health, and only then restores routing. Failure leaves routing closed and both artifacts available for diagnosis.
