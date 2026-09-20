# Contributing

Changes follow the repository's `workflow.toml` sequence: PRD, architecture,
implementation specification, feature scope, implementation, verification,
evidence, and review. Keep this repository limited to the one webhook
extension identity. Shared foundation changes belong in `gotth-extensions`;
Mail policy and administration belong in `gotth-mail`.

Run `make verify` before proposing admission. Do not add generic invocation,
arbitrary headers or payloads, prompt/mutation authority, hidden retries,
secret values in configuration metadata, or additional provider mechanisms.
