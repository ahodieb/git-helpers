# git-helpers

![](https://img.shields.io/badge/WARNING-rewrites%20local%20branch%20refs-red?style=for-the-badge)

Rebases every local branch onto `main` (or `master`) using `git replay` —
no checkout, no working-tree changes. Branches checked out in another
worktree are rebased without conflict; branches that can't replay
cleanly are reported and skipped.

Requires git ≥ 2.44 (when `git replay` shipped).

## Install

```bash
go install github.com/ahodieb/git-helpers@latest

# optional: wire it up as `git rebase-all`
git config --global alias.rebase-all '!git-helpers'
```

## Use

```bash
# from anywhere inside a repo
git-helpers

# or with the alias
git rebase-all

# from outside, point at a repo
git-helpers -C /path/to/repo
```
