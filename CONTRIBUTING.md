# Contributing to Valethra

Keep contributions focused, reviewable, and explicit about whether they change code, game design, or world canon.

## Workflow

1. Create a short-lived branch from `main`.
2. Describe the player-facing or worldbuilding outcome before implementation details.
3. Make the smallest coherent change and add focused tests for changed rules.
4. Run `make check`.
5. Open a pull request explaining scope, validation, screenshots for visual changes, and remaining risks.

Commits should be imperative and scoped, for example `Add dialogue action routing` or `Clarify Valethra magic canon`.

## Quality bar

- No new warnings from `go vet`.
- All tests and the build pass.
- Go files are formatted with `gofmt`.
- Documentation describes the current implementation accurately.
- Proposed features remain under `docs/architecture/proposals/` or `docs/plans/` until implemented.
- Canon changes are intentional and reviewed separately from mechanical refactors when possible.

## Dependencies and assets

Explain why a new dependency is needed and prefer stable releases. For new assets, include origin, author, license, modifications, and where attribution must appear before merging.
