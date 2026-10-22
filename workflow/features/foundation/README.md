# Foundation

State: `blocked`

Implements the bounded webhook alert extension, immutable artifact, and GOTTH
Mail lifecycle integration described by the governing PRD, architecture, and
implementation specification.

Implementation and real cross-process host integration are complete. Release
admission remains blocked on public GitHub creation, one-way mirror setup,
immutable tag/artifact publication, release parity, and deployed acceptance.

The exact clean candidate is
`3627c43213b82afa290372a6e2d9c36978a54758`. Forgejo `main` contains that
candidate plus later workflow evidence. The eventual release tag must resolve
exactly to the tested candidate, not to an evidence-only descendant. No
release tag exists; creating one before the public mirror and release channel
exist would violate the distribution contract.

## Superseding closeout status — 2026-10-21

The preceding no-mirror/no-tag statements are historical. The public GitHub
mirror now exists, and Forgejo and GitHub expose the same annotated
`v1.0.0-alpha.1` tag object
`a126253f637ee9cf8b60482640f5f99c128a2d24`, peeled to the tested candidate
`3627c43213b82afa290372a6e2d9c36978a54758`.

Foundation remains `blocked`. Forgejo and GitHub have no published Release
objects for this tag. The GitHub tag fallback URL is not publication. Fresh
exact-tag verification on 2026-10-21 reproduced the recorded archive exactly.
Later GOTTH Mail evidence binds that archive to real packaged delivery and
A→B→A administrator rollback on the development test equipment; this is
test-equipment evidence, not a live production change. The closeout candidate
adds the missing post-body transport-close regression; the complete fresh gate
passed on development.
Two current clean reviews plus Forgejo/GitHub prerelease-object and exact asset
publication remain open.

The bounded closeout check is:

```sh
./scripts/check-release-closeout.sh
```

Current-state mode verifies authenticated Forgejo absence and public GitHub
published-Release absence; GitHub draft state is unknown without authentication.
Complete mode requires matching non-draft prereleases and exact assets.
It does not, by itself, prove artifact identity, lifecycle behavior, or review.
