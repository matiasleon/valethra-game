# Coding Rules

## Go Best Practices

### Naming
- Use short, descriptive names. Prefer `srv` over `myServer`, `cfg` over `configuration`.
- Exported names should be clear without the package name: `entities.Character`, not `entities.EntityCharacter`.
- Interfaces with a single method use the `-er` suffix: `Reader`, `Writer`, `Attacker`.
- Acronyms should be all caps: `ID`, `HTTP`, `URL`.

### Error Handling
- Always handle errors explicitly. Never discard them with `_`.
- Return errors as the last return value.
- Wrap errors with context using `fmt.Errorf("doing X: %w", err)`.
- Use custom error types when callers need to distinguish error cases.

### Functions
- Keep functions short and focused on a single responsibility.
- Limit function parameters. If more than 3, consider using a struct.
- Return early to reduce nesting. Prefer guard clauses over deep if/else chains.

### Structs and Methods
- Use pointer receivers when the method modifies state or the struct is large.
- Use value receivers for small, immutable structs.
- Group related fields together in struct definitions.

### Packages
- Package names should be short, lowercase, and singular: `entities`, `core`, `game`.
- Avoid generic package names like `util`, `common`, `helpers`.
- A package should have a single, clear purpose.

### Concurrency
- Don't start goroutines without a clear shutdown strategy.
- Use channels to communicate, not shared memory.
- Always handle context cancellation.

### Testing
- Name tests as `Test<Function>_<Scenario>_<ExpectedResult>`.
- Use table-driven tests when testing multiple cases of the same function.
- Keep test setup close to the assertion. Avoid shared mutable state between tests.
- Tests should be independent and runnable in any order.

---

## Clean Code

### Readability
- Code should read like prose. If a block needs a comment to explain what it does, consider refactoring it into a well-named function.
- Avoid magic numbers. Use named constants.
- Keep files focused. One primary type or concept per file.

### Simplicity
- Write the simplest solution that works. Refactor only when complexity is justified.
- Delete dead code. Don't comment it out "just in case" — that's what git is for.
- Avoid premature optimization. Profile first, optimize second.

### Dependencies
- Accept interfaces, return structs.
- Keep the dependency graph shallow. A package should not know about distant parts of the system.

---

## Clean Architecture

### Layer Structure
```
domain/          → Core business entities and rules (no external dependencies)
core/            → Use cases and application logic
infrastructure/  → External systems (DB, APIs, file I/O)
cmd/             → Entry points (main packages)
```

### Rules
1. **Dependencies point inward.** Outer layers depend on inner layers, never the reverse.
2. **Domain has zero imports** from other project packages. It defines the core types and business rules.
3. **Core orchestrates** business logic using domain types. It defines interfaces that infrastructure implements.
4. **Infrastructure implements** interfaces defined by core. Database, HTTP clients, file storage live here.
5. **No business logic in main.** The `main` function wires dependencies together and starts the application.

### Interfaces
- Define interfaces where they are consumed, not where they are implemented.
- Keep interfaces small. One or two methods is ideal.
- Use interfaces at layer boundaries to decouple packages.

### Data Flow
- External input → Infrastructure → Core → Domain → Core → Infrastructure → External output.
- Each layer translates data into the format needed by the next layer.
- Domain types should never contain framework-specific annotations unless absolutely necessary.
