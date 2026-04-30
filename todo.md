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

## Phase 4 — Rewrite `rebase-all` on `git replay`
- [ ] Probe `git --version` ≥ 2.44 at startup; fail with clear message otherwise
- [ ] New flow per branch: `git replay --onto <main> <merge-base>..<branch>` piped to `git update-ref --stdin` — never checks out anything
- [ ] Remove the manual checkout / defer-restore dance in `Git.Rebase`
- [ ] Branches checked out in another worktree just work (refs only)
- [ ] Branches that conflict are reported as skipped; others continue
- [ ] New scenarios:
  - [ ] branch attached to a worktree (currently broken — this is what fixes it)
  - [ ] one branch conflicts on replay; others still rebase
  - [ ] starting from a feature branch (no longer needs `main` to be checked out)

## Phase 5 — Cleanup
- [ ] Retire `checkout-main` subcommand and the `git main` alias (no longer earns its keep)
- [ ] Cherry safety check should target `main`, not HEAD's upstream — currently only works when run from `main`. Fold into Phase 4 redesign or address standalone.
- [ ] Bump `.tool-versions` to current Go + golangci-lint (still pins 1.19.3 / 1.50.1, out of sync with `go.mod` go 1.23 / toolchain 1.24.2)
- [ ] Update README: drop `checkout-main`, note git 2.44+ requirement
- [ ] Bump `cmd/version.json` on the next release

## Phase 6 — Future feature: configurable rebase target
- [ ] `rebase-all --onto <ref>`, default `origin/main` (resolve remote-side main/master)
- [ ] Auto-`git fetch <remote>` by default, with `--no-fetch` opt-out
- [ ] Replace the cherry safety check (no longer meaningful when target is a remote-tracking ref) with "fetch succeeded"
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
### Phase 4 — not started
### Phase 5 — not started
### Phase 6 — not started (depends on Phase 4)
