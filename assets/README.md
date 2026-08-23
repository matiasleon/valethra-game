# Source assets

This directory contains editable source art and upstream sprite sheets. It is intentionally a nested Go module so Go tooling does not interpret asset directories with spaces or parentheses as package import paths.

Runtime-ready assets are copied deliberately to `game/assets/` and embedded by `game/assets.go`. Before distributing the game, record origin, author, license, modifications, and attribution requirements for every retained asset.
