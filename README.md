# git-helpers

![](https://img.shields.io/badge/version-v0.0.3%20(beta)-green?style=for-the-badge)
![](https://img.shields.io/badge/WARNING-can%20damage%20your%20repository%20structure-red?style=for-the-badge)

Provides a collection of helper functionality that I used every day

## Installation

```bash
# install cli
go install github.com/ahodieb/git-helpers

# setup git aliases for a quicker flow
git-helpers install-git-aliases
```

## Helpers

### `checkout-main`

Checkout the main branch.
Many projects have migrated from `master` to `main` as the main branch name, so this switch to whatever is the main
branch in the current repository

```bash
git-helpers checkout-main

# with git aliases setup
git main
```

### `rebase-all`

Rebase all the other branches on top of the `main` branch.
When I work on multiple feature branches locally I like to rebase all of them onto main frequently.

```bash
git-helpers rebase-all

# with git aliases setup
git rebase-all
```

