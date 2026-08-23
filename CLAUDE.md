# CLAUDE.md

## Project Overview

**valethra-game** is a fantasy RPG game engine set in the world of Valethra, written in Go 1.21 with zero external dependencies. Early-stage development focused on core combat mechanics and entity architecture.

## Tech Stack

- **Language:** Go 1.21
- **Dependencies:** Standard library only
- **Linting:** golangci-lint
- **Formatting:** goimports, gofumpt (via gopls)

## Project Structure

```
domain/          → Core business entities (zero external imports)
core/            → Use cases and application logic (entities, game master)
infrastructure/  → External systems (DB, APIs, file I/O) — not yet created
docs/            → Architecture docs, world lore, fight system rules
main.go          → Entry point
```

Dependencies point inward: `infrastructure → core → domain`. Domain has zero imports from other project packages.

## Build & Run Commands

```bash
go build ./...          # Build
go test ./... -v        # Run tests (verbose)
go run main.go          # Run
golangci-lint run       # Lint
goimports -local valethra-game ./...  # Format
```

## Coding Conventions

Defined in `coding-rules.md`. Key rules:

- **Naming:** Short, descriptive. Exported names clear without package prefix. Acronyms ALL CAPS (`ID`, `HTTP`). Single-method interfaces use `-er` suffix.
- **Errors:** Always handle explicitly. Return as last value. Wrap with context: `fmt.Errorf("doing X: %w", err)`.
- **Functions:** Short, single responsibility. Max 3 params (use struct if more). Early returns.
- **Receivers:** Pointer for mutation/large structs, value for small/immutable.
- **Packages:** Short, lowercase, singular. No generic names (`util`, `common`).
- **Testing:** `Test<Function>_<Scenario>_<ExpectedResult>`. Table-driven for multiple cases. Independent, order-independent.
- **Clean code:** No magic numbers. One primary type per file. Simplest solution first. Delete dead code.

## Architecture

Clean Architecture with layers: `domain → core → infrastructure → cmd`.

- Interfaces defined at consumption point, not implementation
- Keep interfaces small (1-2 methods)
- No business logic in main (dependency wiring only)
- Each layer translates data format for the next

## Key Entities

- **Character** — Player/NPC/Enemy with Attributes. Methods: `Attack(target)`, `Spell(target)`
- **Attributes** — Health, AttackPower, Armor, Agility, Intelligence, Willpower, Speed, Level
- **Quest/Event** — Quest has Events, Events have Enemies
- **GameMaster** — Orchestrates quests and game flow
- **Game/Round** — Game contains Rounds (placeholder)

## Current Status

- Combat system partially implemented (some tests failing on boundary conditions)
- Spell system, GameMaster orchestration, and Round/turn system are stubs
- 5 of 11 tests currently failing
