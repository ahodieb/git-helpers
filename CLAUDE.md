# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

A single-binary Go tool that rebases every local branch onto `main`/`master` using `git replay` (git ≥ 2.44). No subcommands, no checkout — just refs being rewritten via `git replay | git update-ref --stdin`. The whole program is ~60 lines in `main.go`.

The README warns it rewrites local branch refs — that's the whole point, but be deliberate when changing the rebase loop.

## Common commands

```bash
go build ./...
go test ./...
go test . -run TestScripts/happy   # one CLI scenario
go test ./git -run TestListBranches

# testscript dev flags (off by default; on-in-CI is wrong)
TESTWORK=1 go test .               # keep WorkDir for inspection
UPDATE_SCRIPTS=1 go test .         # auto-update failing scripts
```

## Architecture

Two source files do all the work:

- **`main.go`** — `flag` parses `-C path` and `--onto ref`; `rebaseAll(g, onto)` resolves the base ref (explicit `--onto` or first of local `main`/`master`), lists branches, calls `g.Replay(base, branch)` per non-base branch, logs failures and continues. That's the whole program.
- **`git/git.go`** — `Git{WorkingDir}` wraps the `git` CLI. `(*Git).Exec` returns trimmed stdout for callers; the private `run(stdin, args...)` returns raw stdout and is shared between `Exec` and the `update-ref` pipe inside `Replay`. `Replay` runs `git replay --onto <base> <base>..<target>` and pipes the raw output to `git update-ref --stdin` (the trailing newline must survive — that's why `run` doesn't trim and `Exec` does). Errors are signalled by exit code only; git writes status to stderr on success ("Switched to branch 'foo'") so stderr alone isn't a failure signal. On failure the error is enriched with stderr.

## Tests

- **Unit tests** (`git/git_test.go`) — testify/assert against real git repos. The `setupRepo(t, branches...)` helper at the bottom of the file runs a small `sh -c` script in `t.TempDir()` to init a repo and fork branches. 4 tests: stderr-on-success, `-C` honored, non-zero exit propagates as error, ListBranches.
- **CLI testscript scenarios** (`testdata/scripts/*.txtar`) — driven by `main_test.go`. The harness registers `git-helpers` as an in-process command and pre-creates `./repo` with one initial commit per scenario via a `Setup` hook (see the `initRepo` shell snippet). 5 scenarios: `happy`, `feature-branch` (HEAD on a non-main branch), `worktree` (a branch attached to a worktree), `conflict` (one branch fails to replay; siblings succeed), `onto-explicit` (`--onto` ref).

`TestWork` and `UpdateScripts` are gated on env vars (`TESTWORK=1`, `UPDATE_SCRIPTS=1`) so CI fails on stale snapshots and dev can opt in when iterating.

## Hard requirement

Git ≥ 2.44 (when `git replay` shipped). No explicit version probe — git's own "'replay' is not a git command" error is clear enough if a user tries with an older binary.
