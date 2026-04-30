# git-helpers — work plan

## Phase 1 — Correctness fixes
- [x] `(*Git).Exec` — stderr-on-success no longer treated as error; `-C WorkingDir` actually applied
- [x] Drop `fmt.Errorf(stderr.String())` non-constant format
- [x] `version.go` `commit[:10]` panic guard
- [x] Typo `abborting` → `aborting`

## Phase 2 — Regression suite (minimal baseline)
- [x] Bump `rogpeppe/go-internal` so the suite builds on modern Go
- [x] Extract `cmd.Run(args) int` so the CLI can run in-process
- [x] Wire CLI into testscript harness with a shared per-scenario `Setup`
- [x] Unit tests that lock down both Phase 1 bugs
- [x] CLI scenarios: `help`, `checkout-main` (master, missing), `rebase-all` (happy, unpublished-main)
- [x] Low-level Exec scenarios: stderr-on-success, exit-code propagation

## Phase 3 — Consolidate to one `Exec`
- [x] Delete the string-returning `(*Git).Exec`; rename `Exec2` → `Exec`
- [x] All callers go through the `Output{Stdout, Stderr, ExitCode}` struct; string callers use `out.TrimmedStdout()`
- [x] Error semantics in one place: exit code drives error, stderr is data, error message enriched with stderr at the package-level `Exec`

## Phase 3.5 — Code-review cleanup
- [x] Replace functional-options + generics in `git/scripts/git.go` with simple `Init` / `InitOnBranch`
- [x] Inline `executer.go` into a `run` helper using `sh -c` (no tempfile, no Params struct); delete the file
- [x] Rename receiver `git` → `g` throughout `git/git.go`
- [x] Trim doc comments that echoed signatures; drop speculative future-work prose; drop the vague `--abort` TODO
- [x] `Output.StdoutLines()` returns `nil` on empty stdout (was `[""]`); collapses `Cherry`'s manual filter
- [x] Inline `Checkout` (one external caller); use `ExecSilent` for the discarded `--abort`
- [x] `cmd/rebase-all.go`: include `err` in failure log, drop the `remove` helper for an inline `continue`, fix misleading `return err`
- [x] `cmd/install-git-aliases.go`: don't log "installed" after a failure
- [x] `cmd/version.go`: typed struct instead of `map[string]string`; nil-check `ReadBuildInfo`
- [x] `cmd/run.go`: lowercase `App()` → `app()` (internal-only)
- [x] Gate `TestWork` and `UpdateScripts` behind `TESTWORK=1` / `UPDATE_SCRIPTS=1` env vars in both testscript runners
- [x] `git/git_scripts_test.go`: print only on success; delete commented-out `Cmds:` block

## Phase 4 — Rewrite `rebase-all` on `git replay`
- [x] New flow per branch: `git replay --onto <main> <main>..<branch>` piped to `git update-ref --stdin` — never checks out anything
- [x] Remove the manual checkout / defer-restore dance (deleted `Git.Rebase`)
- [x] Branches checked out in another worktree just work (refs only)
- [x] Branches that conflict are reported as skipped; others continue
- [x] Cherry safety check dropped (no longer applies — Phase 6 will replace with auto-fetch)
- [x] Deleted `Git.Cherry` (only caller was the dropped check)
- [x] Scenarios:
  - [x] branch attached to a worktree
  - [x] one branch conflicts on replay; others still rebase
  - [x] starting from a feature branch
- [ ] Skipped: explicit version probe — git produces a clear "is not a git command" error itself if `replay` is missing

## Phase 5 — Minimize
With `git replay` doing the work, the Go scaffolding stopped earning its keep. Aggressively simplify and keep only what supports the test suite (the project's main remaining value).

- [x] Drop `cmd/` entirely — collapse into `package main` with stdlib `flag`
- [x] Drop `urfave/cli` dependency
- [x] Delete `checkout-main`, `install-git-aliases`, `version` subcommands and `version.json`
- [x] Delete dead methods: `Git.Checkout`, `Git.AddAlias`, `Git.ExecSilent` (no callers)
- [x] Move scenarios from `cmd/testdata/` to `testdata/`, drop the now-redundant `rebase-all.` prefix
- [x] Drop `help`, `checkout-main.master`, `checkout-main.missing` scenarios (subjects gone)
- [x] Update README + CLAUDE.md to reflect the single-binary shape
- [ ] Bump `.tool-versions` to current Go + drop the golangci-lint pin if unused

## Phase 6 — Future feature: configurable rebase target (deferred)
The user is considering whether to keep the project at all. If kept and extended:
- [ ] `--onto <ref>`, default `origin/main` (resolve remote-side main/master)
- [ ] Auto-`git fetch <remote>` by default, with `--no-fetch` opt-out
- [ ] Add a "fetch succeeded" precondition check (replaces the dropped cherry check)
- [ ] Document the worktree workflow as the primary use case

---

## Status

### Phase 1 — done
Both `(*Git).Exec` bugs were one underlying problem (the older path duplicated and got it wrong); the fix routes it through the corrected package-level `Exec`/`Output`, which also cleared the `fmt.Errorf` lint. `version.go` and the typo were one-line touches. Side effect: had to bump `rogpeppe/go-internal` to v1.14.1 because v1.10.0 didn't build on Go 1.24 — `go.mod` auto-bumped to `go 1.23` + `toolchain go1.24.2`.

### Phase 2 — done
Suite is `go test ./...`. Three layers:
- **Unit tests** (`git/git_test.go`): `TestExecDoesNotErrorOnStderrWhenCommandSucceeds`, `TestExecRespectsWorkingDir` — both written failing-first against the Phase 1 bugs, now passing.
- **CLI scenarios** (`cmd/testdata/scripts/`): driven by a new `cmd/cmd_test.go` that registers `git-helpers` as an in-process testscript command and pre-creates `./repo` per scenario via a `Setup` hook.
- **Low-level Exec scenarios** (`git/testdata/scripts/`): added two alongside the pre-existing ones.

Surfaced — not fixed — while writing scenarios: `rebase-all` only works cleanly when started from `main` because `Cherry()` runs `git cherry -v` against HEAD's upstream rather than main's. Captured in Phase 5.

### Phase 3 — done
One canonical `(*Git).Exec(args) (Output, error)` and one canonical package-level `Exec`. The error-context behavior the deleted wrapper used to provide (enriching errors with stderr) was lifted into the package-level `Exec` so every caller benefits without paying for it. `Cherry` switched from string-splitting to `out.StdoutLines()` with whitespace filtering. Everything else (`Checkout`, `AddAlias`, `Rebase`) already went through `ExecSilent`, so no caller changes there.
### Phase 3.5 — done
Worst single offender (`scripts.Init` functional-options pattern + generic `apply` for a 2-field struct, only one option ever defined) deleted. `executer.go` inlined and the tempfile-then-exec dance replaced with `sh -c`. `git.go` is shorter, receiver is `g`, doc comments don't echo signatures, `StdoutLines` no longer ships `[""]` on empty input. `cmd/` footguns fixed (failure-then-success log on alias install, discarded errors, misleading `return err`). `TestWork`/`UpdateScripts` are env-gated so CI can't auto-update its way past failures. `go test ./...` and `go vet ./...` clean.

### Phase 4 — done
`Git.Replay(base, target)` runs `git replay --onto <base> <base>..<target>` and pipes the resulting `update refs/heads/...` lines to `git update-ref --stdin`. No checkout, no working-tree state, no per-branch defer dance. `cmd/rebase-all.go` is now a flat loop: list branches, skip main, replay each, log-and-continue on failure. Cherry check + `Git.Cherry` + `Git.Rebase` all deleted as dead code. Old `rebase-all.unpublished-main.txtar` scenario removed (it locked in the cherry semantics we just dropped). Three new scenarios prove the wins: worktree-attached branch still rebases, conflicting branch is reported and skipped while siblings succeed, running from a feature branch works without ceremony. `rebase-all.happy.txtar` simplified — no longer needs a bare origin since cherry's gone.

Hard requirement: git ≥ 2.44 (when `replay` shipped). No explicit probe — git's own "'replay' is not a git command" error is clear enough.

### Phase 5 — done (mostly)
Down from 384 LOC / 9 Go files to **231 LOC / 3 Go files**. Scaffolding gone, behavior unchanged. `urfave/cli` dropped from `go.mod` along with its transitive deps. `main.go` is ~50 lines of plain stdlib + the rebase loop. Only `git/git.go` (124 lines, the `git` CLI wrapper + `Replay`) and `git/scripts/git.go` (55 lines, test fixtures) remain alongside it. All 9 testscript scenarios + 5 unit tests still green.

Open: `.tool-versions` still pins golang 1.19.3 / golangci-lint 1.50.1 — neither is used now.

### Phase 6 — deferred pending archive decision
