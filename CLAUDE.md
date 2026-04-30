# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

A single-binary Go tool that rebases every local branch onto `main`/`master` using `git replay` (git ≥ 2.44). No subcommands, no checkout — just refs being rewritten via `git replay | git update-ref --stdin`. The whole program is ~50 lines in `main.go`.

The README warns it rewrites local branch refs — that's the whole point, but be deliberate when changing the rebase loop.

## Common commands

```bash
go build ./...
go test ./...
go test . -run TestScripts/happy        # one CLI scenario
go test ./git -run TestFindMainBranch   # one unit test

# testscript dev flags (off by default, on in CI is wrong)
TESTWORK=1 go test .                    # keep WorkDir for inspection
UPDATE_SCRIPTS=1 go test .              # auto-update failing scripts
```

## Architecture

Three files matter:

- **`main.go`** — `flag` parses `-C path`; `rebaseAll(g)` finds main, lists branches, calls `g.Replay(main, branch)` per non-main branch, logs failures and continues. That's the whole program.
- **`git/git.go`** — wrapper around the `git` CLI. `(*Git).Exec` prepends `-C WorkingDir` and delegates to the package-level `Exec`, which is the only place that touches `os/exec`. `Replay` is the one non-trivial method: it runs `git replay --onto <base> <base>..<target>`, then pipes the stdout to `git update-ref --stdin` (the second call uses raw `exec.Command` because it needs `Stdin`). Errors are signalled by exit code only — git writes status to stderr on success ("Switched to branch 'foo'") so stderr alone isn't a failure signal. On failure the error message is enriched with stderr.
- **`git/scripts/git.go`** — test fixture helpers: `Init`, `InitOnBranch`, `CreateBranch`, `Checkout`. Each runs a tiny `sh -c` script against a real `git` binary in `t.TempDir()`. Used by unit tests.

## Tests

Three styles:

- **Unit tests** (`git/git_test.go`) — testify/assert against real git repos created via `scripts.Init(t)`.
- **CLI testscript scenarios** (`testdata/scripts/*.txtar`) — driven by `main_test.go`. The harness registers `git-helpers` as an in-process command and pre-creates `./repo` with one initial commit per scenario via a `Setup` hook (see the `initRepo` shell snippet). Scenarios cover: `happy` (basic flow), `feature-branch` (HEAD on a non-main branch — no longer needs main checked out), `worktree` (a branch attached to a worktree still gets rebased), `conflict` (one branch fails to replay; siblings still succeed).
- **Low-level `git.Exec` testscript scenarios** (`git/testdata/scripts/*.txt`) — driven by `git/git_scripts_test.go`. The `TestMain` registers `git.Exec`, `git.CurrentBranch`, `git.ListBranches` so testscripts can call them directly. Scenarios lock down stderr-on-success and exit-code propagation.

`TestWork` and `UpdateScripts` are gated on env vars (`TESTWORK=1`, `UPDATE_SCRIPTS=1`) so CI fails on stale snapshots and dev can opt in when iterating.

## Hard requirement

Git ≥ 2.44 (when `git replay` shipped). No explicit version probe — git's own "'replay' is not a git command" error is clear enough if a user tries with an older binary.
