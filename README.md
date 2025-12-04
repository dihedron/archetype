# archetype

[![Go Report Card](https://goreportcard.com/badge/github.com/dihedron/archetype)](https://goreportcard.com/report/github.com/dihedron/archetype)

Initialise a new repository in the current directory from a template stored in a Git repository.

## How to prepare an archetype repository

The repository should have a `.archetype` directory at its root, containing the archetype metadata.
All the other repository files will be used to initialise the new repository.

The archetype metadata should contain the `parameters.yml` file, containing the parameters to be used to initialise the new repository.



If a reposiotry files

## How to test

An example repository is available at git@github.com:go-git/go-git.git; you can check the v1.0.0 tag:

```bash
$> archetype init -r=https://github.com/go-git/go-git.git -t=v1.0.0 -p=@_test/parameters.yml
```

the repository HEAD (latest):

```bash
archetype init -r=git@github.com:go-git/go-git.git
```

or a specific commit (either by long or short hash):

```bash
archetype init -r=git@github.com:go-git/go-git.git -t=663f81a
```

They will each have different contents.

## How to see the logs

In order to enable the logs, export or set the ARCHETYPE_LOG_LEVEL=d environment variable.
