# repo

Initialise a new repository in the current directory from a template stored in a Git repository.

## How to test

An example repository is available at git@github.com:go-git/go-git.git; you can check the v1.0.0 tag:

```bash
repo -r=git@github.com:go-git/go-git.git -t=v1.0.0
```

the repository HEAD:

```bash
repo -r=git@github.com:go-git/go-git.git
```

or tag v2.2.0:

```bash
repo -r=git@github.com:go-git/go-git.git -t=v2.2.0
```

They will each have different contents.
