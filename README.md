semantic
========

Create and Parse [SemVer 2.0](https://semver.org/) tags in Go

[![Go Report Card](https://goreportcard.com/badge/gophers.dev/pkgs/semantic)](https://goreportcard.com/report/gophers.dev/pkgs/semantic)
[![Build Status](https://travis-ci.com/shoenig/semantic.svg?branch=master)](https://travis-ci.com/shoenig/semantic)
[![GoDoc](https://godoc.org/gophers.dev/pkgs/semantic?status.svg)](https://godoc.org/gophers.dev/pkgs/semantic)
[![NetflixOSS Lifecycle](https://img.shields.io/osslifecycle/shoenig/semantic.svg)](OSSMETADATA)
[![GitHub](https://img.shields.io/github/license/shoenig/semantic.svg)](LICENSE)

# Project Overview

Module `gophers.dev/pkgs/semantic` can be used to parse / create / manipulate
[SemVer 2.0](https://semver.org/) tags.

# Getting Started

The `semantic` module can be installed by running
```
$ go get gophers.dev/pkgs/semantic
```

#### Example Usage
```golang
// create a tag for "v1.2.3"
tag := New(1, 2, 3)

// create a pre-release tag with -rc1
tag := New2(1, 2, 3, "rc1")
```

# Contributing

The `gophers.dev/pkgs/semantic` module is always improving with new features
and error corrections. For contributing bug fixes and new features please file an issue.

# License

The `gophers.dev/pkgs/semantic` module is open source under the [BSD-3-Clause](LICENSE) license.
