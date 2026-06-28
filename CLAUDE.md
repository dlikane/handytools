# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build the main CLI
go build ./cmd/img/

# Run without building
go run ./cmd/img/ <subcommand> [flags] [args]

# Run all tests
go test ./...

# Run tests for a single package
go test ./pkg/common/
go test ./pkg/assemble/

# Run a single test
go test ./pkg/common/ -run TestExpandWildcards_SinglePattern_MatchesFiles
```

## Architecture

The main binary is `cmd/img/img.go` — a Cobra CLI that wires together subcommands, each living in `internal/<name>/`.

**Subcommand structure** — every `internal/` package follows the same two-file pattern:
- `cmd.go` — defines a `Config` struct, a package-level `config` var, a `Cmd` cobra command, and an `init()` that binds flags to `config`
- `worker.go` — contains the business logic, called by the command's `Run` func

**Subcommands:**
- `optimise` — resizes/recompresses images to named profiles (`insta`=1350px, `small`=1440, `med`=1920, `large`=2560, `origin`=no resize). Defaults to dry-run; requires `--apply` to overwrite files. Preserves original `mtime` on apply. `--stat` mode profiles all sizes to a temp file without writing.
- `collage` — assembles images into a fixed grid (rows × cols), scaling each cell uniformly. Supports `4x5` aspect ratio override.
- `gallery` — assembles images into a justified, paginated layout using `pkg/assemble` + `pkg/layout`. Sources: Pinterest board URL (via headless Chrome), directory, text file list, or positional args.
- `rename` / `batchrename` — file renaming utilities.
- `grab` — reads text files matching glob/`...` wildcard patterns and dumps their content to stdout (a code-context tool, not an image tool).

**Shared packages:**
- `pkg/common` — global logrus logger (`GetLogger()`), `ExpandWildcards()` for glob patterns, and `SetDryRunMode(true)` which hooks all log output to prepend `DRYRUN:`.
- `pkg/assemble` — builds a continuous canvas from images using justified row layout, then crops it into pages. Canvas is always 1080px wide, pages break at 1920px height.
- `pkg/layout` — justified layout algorithm (`JustifyWithPageSplits`): packs images into rows respecting `TargetHeight`, `Spacing`, and `Tolerance`, inserting page-break markers when the canvas exceeds `maxHeight`.

**Standalone tools (not subcommands of `img`):**
- `cmd/effect/effects.go` — standalone `main` that applies beat-synced grayscale/invert effects to video frames. Run from `data/` after extracting frames with ffmpeg; see `scripts/README.md` for the full pipeline.
- `cmd/onceoff/` — one-off scripts, not wired into the main CLI.

## Key conventions

- All image ops use `github.com/disintegration/imaging` with Lanczos resampling; JPEG output is always quality 85 (`jpegQuality` const in `optimise`).
- Logging uses the shared logrus singleton from `pkg/common`. Format is `level: message` (no timestamps). Use `logger.WithError(err).Errorf(...)` for errors with context.
- Dry-run is the safe default for destructive commands (`optimise`). Always check `cfg.Apply` before writing files.
- Vendor directory is committed; update with `go mod vendor` after changing `go.mod`.
