# Repository guide for coding agents

## Mission

Evolve the playable Valethra prototype and its world without presenting proposals as implemented features. Preserve existing behavior and established lore unless the task explicitly changes them.

## Start here

1. Read `README.md` and `docs/README.md`.
2. Inspect `git status --short`; preserve unrelated work.
3. Read the nearest nested `AGENTS.md` before editing under `core/`, `game/`, or `docs/world/`.
4. Prefer the smallest coherent change. Do not introduce frameworks or abstractions for hypothetical needs.

## Reproducible commands

- Setup: `make setup`
- Format: `make fmt`
- Unit tests: `make test`
- Static checks: `make lint`
- Build: `make build`
- Full local gate: `make check`
- Run the game: `make run`

Run `make check` after code, module, tooling, or asset-embedding changes. For documentation-only changes, verify links, paths, canon status, and run `make docs-check`.

## Boundaries

- `core/` must not import Ebitengine or `game/`.
- `game/` adapts core behavior to input, rendering, scenes, and embedded assets.
- Keep `main.go` limited to composition, initial content, and application startup.
- Treat `docs/world/canon.md` as established canon. Put uncertain ideas in `docs/world/ideas.md`.
- Treat `docs/architecture/proposals/` and `docs/plans/` as non-implemented work.
- Do not silently change combat rules or lore to make documentation match code; surface the mismatch and resolve it explicitly.
- Do not add dependencies unless the standard library or current stack is insufficient and the value is clear.

## Definition of done

- Changed behavior has focused tests where practical.
- `make check` passes, or the handoff states exactly what could not run and why.
- User-facing behavior and durable architectural decisions are documented.
- No generated binaries, local editor state, secrets, or temporary files are committed.
- The final handoff summarizes behavior, validation, and remaining risk.
