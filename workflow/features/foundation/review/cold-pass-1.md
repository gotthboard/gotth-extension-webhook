# Cold review 1

## Decision

Narrow, fix, and retry.

## Findings

The runtime boundary was sound, but the release builder was not admissible.
It used a loose shell glob for versions, accepted arbitrary/default archive
epochs, wrote directly over the final path, could leave a partial artifact,
did not verify the linked binary version, and inherited host directory modes.
The first rebuilt archive proved the last defect concretely: its root carried
setgid/group-write bits from the development filesystem.

## Resolution

The builder now validates the exact release grammar, binds timestamps to the
source commit epoch, supports mandatory clean-tree release mode, verifies the
linked version, emits deterministic USTAR/gzip, publishes through a
no-replace hard link, synchronizes file and directory state, pins archive
member modes, and explicitly clears inherited setgid. Regression checks cover
every defect above and preserve an existing artifact byte-for-byte.
