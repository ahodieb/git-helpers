# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

`git-helpers` is a small Go CLI that wraps the `git` command with a few personal helpers — most importantly `rebase-all`, which rebases every local branch on top of `main`/`master`. The CLI also installs git aliases (`git main`, `git rebase-all`) so the helpers feel like native git subcommands.

The README warns it can damage repository structure — `rebase-all` rewrites history on every local branch in one go, so be deliberate when changing its behavior.

## Common commands

```bash
# build / install
go install ./...

# run without installing
go run . rebase-all
go run . -C /path/to/repo checkout-main

# tests
go test ./...
go test ./git -run TestFindMainBranch          # single test
go test ./git -run TestScripts/git.ListBranches # single testscript file

# lint (matches .tool-versions: golangci-lint 1.50.1)
golangci-lint run
```

Tool versions are pinned in `.tool-versions` (golang 1.19.3 / golangci-lint 1.50.1). `go.mod` declares `go 1.19`.

## Architecture

Three layers, each in its own package:

- **`main.go`** — wires up `urfave/cli/v2`. Each subcommand maps to a function in `cmd/`. The global `-C/--working-dir` flag is read by every command via `cmd.FlagWorkingDir`.
- **`cmd/`** — one file per subcommand (`checkout_main.go`, `rebase-all.go`, `install-git-aliases.go`, `version.go`). Each builds a `git.Git{WorkingDir: ...}` and orchestrates calls. `version.go` embeds `cmd/version.json` (the source of truth for the version string read by the README badge) and reads VCS info from `runtime/debug.ReadBuildInfo` — so `go install` produces a version with the embedded string + the commit hash.
- **`git/git.go`** — thin wrapper that shells out to the `git` binary via `os/exec`. The `Git` struct just holds `WorkingDir`. Higher-level methods (`FindMainBranch`, `Cherry`, `Rebase`, `ListBranches`, etc.) compose the lower-level `Exec`.

`rebase-all` flow: find main → run `git cherry -v` to confirm main has no unpublished commits → list branches → for each non-main branch, `git rebase <main> <branch>`. After a rebase, HEAD lands on the rebased branch, so `Git.Rebase` saves the original branch up front and `defer`s a checkout back to it. On rebase failure it runs `git rebase --abort`.

### Two parallel Exec implementations

`git/git.go` currently has two execution paths that exist side by side:

- `(*Git).Exec(arg...) (string, error)` — older, returns trimmed stdout as a string.
- `(*Git).Exec2(arg...) (Output, error)` and the package-level `Exec(arg...) (Output, error)` — newer, return an `Output{Stdout, Stderr, ExitCode}` struct. This is what the testscript harness uses.

The newer path (`Exec2` / package `Exec` / `Output`) is the direction the code is moving in — the most recent commits introduced it. New code should use it; treat the old `(*Git).Exec` as legacy.

## Tests

Two styles, both in the `git` package:

- **Standard Go tests** (`git_test.go`) using `testify/assert`. They build temp repos via the `git/scripts` helper package — `scripts.Init(t)` creates a real git repo in `t.TempDir()` by running a small shell script (`scripts/executer.go` writes the script to disk and execs `/bin/sh`).
- **`testscript` tests** (`git_scripts_test.go`) using `rogpeppe/go-internal/testscript`. `TestMain` registers `git.Exec`, `git.CurrentBranch`, and `git.ListBranches` as in-process commands; the `.txt`/`.txtar` files under `git/testdata/scripts/` are the actual test scenarios. `UpdateScripts: true` is on, so a failing testscript can be auto-updated by re-running.

When adding a new helper to `git/git.go`, the convention is: write a unit test in `git_test.go` for the Go API, and (if the helper is exposed as a script command) add a `testdata/scripts/<name>.txt` exercising it end-to-end against a real git binary.
