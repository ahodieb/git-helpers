# Changelog

Project journey, oldest first. Each phase: what changed, then why.

## Phase 1 — Correctness fixes
Fixed two `(*Git).Exec` bugs (stderr-on-success treated as failure; `-C WorkingDir` silently dropped due to slice-mutation order). Guarded a `commit[:10]` panic in `version.go`. Bumped `rogpeppe/go-internal` to v1.14.1 so the suite builds on modern Go.

`rebase-all` logged every successful rebase as a failure — git uses stderr as a status channel even on exit 0. The test suite wouldn't build at all, blocking everything else.

## Phase 2 — Regression suite
Wired the CLI into a testscript harness with a per-scenario `Setup` hook (creates `./repo` with one initial commit). Added 5 CLI scenarios, 2 low-level `Exec` scenarios, and 2 unit tests covering the Phase 1 bugs. Extracted `cmd.Run(args) int` so the binary could be invoked in-process.

Needed a baseline before touching anything substantive — couldn't tell if later simplifications broke real behavior without scenarios that locked it in.

## Phase 3 — Consolidate `Exec`
Deleted the buggy string-returning `(*Git).Exec`; renamed `Exec2` → `Exec`. Error enrichment (stderr included in the message) moved to the package-level `Exec` so every caller benefits.

Two parallel exec paths had drifted apart. The older one was where the Phase 1 bugs lived. One canonical implementation, one place to think about error semantics.

## Phase 3.5 — Code review pass
Replaced functional-options + generics in `git/scripts/` with two plain functions. Inlined `executer.go`. Renamed receiver `git` → `g`. Trimmed signature-echoing doc comments. Fixed `StdoutLines` returning `[""]` on empty input (which had forced `Cherry` to filter it back out). Fixed `cmd/` footguns: install-aliases logging "installed" after a failure, `rebase-all` swallowing errors. Gated dev-only testscript flags (`TestWork`, `UpdateScripts`) behind env vars.

A pile of small over-engineering had accumulated. The functional-options-with-generics for a 2-field struct was the worst offender — abstraction larger than the thing being abstracted.

## Phase 4 — `git replay`
Replaced the checkout-and-rebase dance with `git replay --onto <base> <base>..<target> | git update-ref --stdin`. No checkout, no working-tree changes, no per-branch defer-restore. Dropped the cherry safety check, `Git.Cherry`, `Git.Rebase` and the unpublished-main scenario. Added scenarios for worktree, conflict, and starting from a feature branch.

The old `git rebase` flow failed when a branch was checked out in another worktree — the original motivating bug. `git replay` (≥ 2.44) updates refs without touching the working tree, which sidesteps the problem entirely.

## Phase 5 — Minimize
Collapsed `cmd/` into `package main` with stdlib `flag`. Dropped `urfave/cli` and its transitive deps. Deleted `checkout-main`, `install-git-aliases`, `version` subcommands and `version.json`. Reorganized scenarios. Rewrote README and CLAUDE.md.

With `git replay` doing the heavy lifting, the program was ~10 lines of logic wrapped in 384 lines of scaffolding. The honest comparison was a bash script plus the test suite; kept Go specifically for the test suite. From 9 source files / 384 LOC down to 3 source files / 231 LOC.

## Phase 6 — `--onto` flag
Added `--onto <ref>`. When unset, falls back to local `main` or `master`. No auto-fetch, no remote logic. Dropped now-dead `Git.FindMainBranch`, `Git.CurrentBranch`, `scripts.InitOnBranch`, and `scripts.Repository.Checkout` along with their tests.

Sometimes the rebase target should be something other than `main` (typically `origin/main` after a manual fetch). The user does their own fetching — keeps the tool predictable, no surprise side effects from a `--no-fetch` toggle.

Final size: **165 LOC across 2 source files**, plus 55 lines of test fixtures. 12 tests green (5 CLI scenarios, 4 low-level scenarios, 3 unit tests).

## Phase 7 — Second simplification pass
Collapsed `Output{Stdout, Stderr, ExitCode}` plus `TrimmedStdout()` and `StdoutLines()` methods into a plain `(string, error)` return on `Exec`. Folded the package-level `Exec` into the method; added a private `run(stdin, args...)` helper shared between `Exec` and the `update-ref` pipe in `Replay`. Inlined `resolveBase` into `rebaseAll`. Inlined the `git/scripts/` test fixture package as a `setupRepo` helper directly in `git_test.go`. Deleted the low-level `git/git_scripts_test.go` runner and its `testdata/scripts/` (3 of its 4 scenarios duplicated unit-test coverage; the 4th — exit-code propagation — became irrelevant once `Exec` returns just `(string, error)`, replaced by one new unit test asserting non-zero exit returns an error).

Reviewing the project as a glorified bash script — a loop over branches calling `git replay`. Anything beyond that loop and the test suite was dead weight. The `Output` struct existed to be `Output`; the package-level `Exec` existed because old testscripts had registered it as an in-process command; the `scripts/` subpackage existed because the helpers had been there before they were trimmed.

Final size: **124 LOC of production code across 2 source files** (main.go 59, git/git.go 65), plus 121 lines of tests across 2 files. 9 tests green (5 CLI scenarios, 4 unit tests). One Watch-dropped detail caught and fixed mid-pass: trimming Exec output broke `update-ref --stdin`, which is line-terminated — `Exec` now trims at the public boundary while the private `run` returns raw stdout for the pipe path.
