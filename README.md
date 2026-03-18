# gghq

[![CI](https://github.com/kotaoue/gghq/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/kotaoue/gghq/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/kotaoue/gghq/branch/main/graph/badge.svg)](https://codecov.io/gh/kotaoue/gghq)
[![Go Report Card](https://goreportcard.com/badge/github.com/kotaoue/gghq)](https://goreportcard.com/report/github.com/kotaoue/gghq)
[![License](https://img.shields.io/github/license/kotaoue/gghq)](https://github.com/kotaoue/gghq/blob/main/LICENSE)

Extract the local path from ghq get output.

## Installation

```sh
go install github.com/kotaoue/gghq@latest
```

## Usage

```sh
gghq <repository>
```

`gghq` runs `ghq get` on the given repository, then prints the full local path where it was cloned.

### Examples

```sh
$ gghq https://github.com/example/repo
/home/user/ghq/github.com/example/repo

$ gghq git@github.com:example/repo.git
/home/user/ghq/github.com/example/repo

$ gghq example/repo
/home/user/ghq/github.com/example/repo
```

You can use the output directly to change into the cloned directory:

```sh
cd $(gghq example/repo)
```
