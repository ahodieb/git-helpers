# Rebaser 

![](https://img.shields.io/badge/version-v0.0.1%20(beta)-green?style=for-the-badge)
![](https://img.shields.io/badge/WARNING-can%20damage%20your%20git%20history-red?style=for-the-badge)

is a git utility cli that rebases all branches onto a a specified branch. 

It addresses a common use case i have during my daily work, when i have multiple feature branches locally i like to rebase all of them onto main frequently, previously i had some bash scripts for it 


### Usage

```bash
# rebases all branches onto `main or master` default main branch

git rebase-all 
```

```bash
# rebases all branches onto `<branch-name>` 

git rebase-all <branch-name>
```


### Installation

```
go get github.com/ahodieb/rebaser
git config --global alias.rebase-all '!rebaser rebase-all'
```