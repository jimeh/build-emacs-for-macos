# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with
code in this repository.

## Project Overview

macOS Emacs building system with Ruby and Go components:

- **Ruby script** (`build-emacs-for-macos`): Main build script for creating
  self-contained Emacs.app bundles
- **Go CLI tool** (`emacs-builder`): Packaging, signing, notarization, and
  release management
- **Dual dependency management**: Nix (preferred) or Homebrew

## Common Commands

```bash
# Environment setup (Nix preferred)
nix develop                       # Default macOS 11 SDK
nix develop .#macos{11-15,26}     # Target specific SDK version

# Go development
make build                        # Build emacs-builder CLI
make test                         # Run tests with race detection
make lint                         # Run golangci-lint
make format                       # Format with gofumpt

# Ruby development
bundle exec rubocop               # Lint (with development group)

# Build Emacs
./build-emacs-for-macos           # Build from master
./build-emacs-for-macos emacs-29.4
```

## Architecture

### Ruby Build Script (`build-emacs-for-macos`)

Ruby build script (~2500 lines), with small supporting classes in `lib/`, that:

- Downloads source tarballs from emacs-mirror/emacs on GitHub
- Configures and compiles Emacs with native-comp, tree-sitter support
- Creates self-contained .app bundles by embedding/relinking dependencies
- Uses `ruby-macho` gem for Mach-O binary manipulation (RPATH handling)
- Adds the privacy descriptions defined in `lib/privacy_usage_descriptions.rb`
  before self-signing. Keep this step before all signing, and keep
  legacy/current key pairs together so the default macOS 11.3+ deployment
  range remains covered.

### Go CLI (`cmd/emacs-builder/`)

Uses `urfave/cli/v2` framework. Key packages in `pkg/`:

- `cli/`: Commands (plan, sign, sign-files, notarize, package, release, cask)
- `sign/`: macOS code signing via `codesign`
- `notarize/`: Apple notarization workflow via `notarytool`
- `release/`: GitHub release management
- `dmgbuild/`: DMG creation using Python dmgbuild
- `plan/`: Build plan JSON parsing and management

### Nix Environment (`flake.nix`)

- Multi-SDK support: macOS 11-15, 26 via `.#macos{11,12,13,14,15,26}`
- Excludes ncurses intentionally (links against system version for TUI)
- Sets `MACOSX_DEPLOYMENT_TARGET`, `DEVELOPER_DIR`, `NIX_LIBGCCJIT_*`
- Keep the full host Command Line Tools SDK out of the Nix build-time
  `LIBRARY_PATH`; expose only the sanitized system ncurses stub. Nix supplies
  the remaining libraries matched to its linker, while newer host SDK stubs
  may be incompatible with that toolchain.

## Testing

```bash
make test                         # All Go tests
make test-ruby                    # Privacy usage-description Ruby tests
go test ./pkg/release/...         # Single package
go test -run TestName ./pkg/...   # Single test
```

Tests use `_test.go` suffix alongside source files.

## Working Directories

- `sources/`: Downloaded/extracted Emacs source (gitignored)
- `builds/`: Build outputs and .app bundles (gitignored)
- `patches/`: Emacs patches applied during build
- `bin/`: Built Go binaries (gitignored)
