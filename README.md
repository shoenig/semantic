semantic
========

Create and Parse [SemVer 2.0](https://semver.org/) tags in Go

[![GitHub](https://img.shields.io/github/license/shoenig/semantic.svg)](LICENSE)

# Project Overview

Module `github.com/shoenig/semantic` can be used to parse / create / manipulate
[SemVer 2.0](https://semver.org/) tags.

# Getting Started

The `semantic` module can be installed by running
```
$ go get github.com/shoenig/semantic
```

#### Example Create Tag
```golang
// create a tag for "v1.2.3"
tag := semantic.New(1, 2, 3)

// create a pre-release tag with -rc1
tag := semantic.New2(1, 2, 3, "rc1")

// create a tag with build-metadata
tag := semantic.New4(1, 2, 3, "exp.sha.5114f86")

// create a pre-release tag with build-metadata
tag := semantic.New3(1, 2, 3, "rc1", "exp.sha.5114f86")
```

#### Example Parse Tag
```golang
tag, ok := semantic.Parse("v1.2.3-rc1+linux")
```

#### Example Sort Tags
```golang
tags := []semantic.Tag{
    semantic.New2(1, 1, 1, "alpha.1"),
    semantic.New(1, 1, 1),
    semantic.New2(1, 1, 1, "alpha"),
    semantic.New2(1, 1, 1, "beta.11"),
    semantic.New2(1, 1, 1, "beta.2"),
    semantic.New2(1, 1, 1, "rc.1"),
    semantic.New2(1, 1, 1, "beta"),
    semantic.New2(1, 1, 1, "alpha.beta"),
}
sort.Sort(semantic.BySemver(tags))
```

# Contributing

The `github.com/shoenig/semantic` module is always improving with new features
and error corrections. For contributing bug fixes and new features please file an issue.

# License

The `github.com/shoenig/semantic` module is open source under the [BSD-3-Clause](LICENSE) license.
