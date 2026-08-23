# Core domain guidance

- Keep this tree deterministic and independent from rendering, input, filesystem, and Ebitengine.
- Preserve domain invariants: health and armor never become negative; defeated characters cannot attack; combat outcomes remain reproducible unless a task explicitly adds randomness.
- Define interfaces at the point of consumption and keep them narrow.
- Prefer table-driven tests for rule boundaries such as zero health, zero armor, and damage overflow.
- Run `go test -race ./core/...` after changing this tree.
