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
