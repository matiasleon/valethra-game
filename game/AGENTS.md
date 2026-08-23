# Game adapter guidance

- This tree owns Ebitengine integration, scene flow, rendering, input, layout, animation, and embedded runtime assets.
- Keep game rules in `core/`; scenes should translate input into core operations and render their results.
- Preserve the `ebiten.Game` lifecycle: update state in `Update`, draw in `Draw`, and keep `Layout` free of mutations.
- Avoid tests that require a visible window. Extract pure layout or transition logic when it needs unit coverage.
- Verify asset paths against `game/assets.go`; embedded assets must exist at build time.
