# impacca

[![Release](https://img.shields.io/github/release/ShogunPanda/impacca.svg)](https://github.com/ShogunPanda/impacca/releases/latest)
[![GoDoc](https://godoc.org/github.com/ShogunPanda/impacca?status.svg)](https://godoc.org/github.com/ShogunPanda/impacca)
[![Go Report Card](https://goreportcard.com/badge/github.com/ShogunPanda/impacca)](https://goreportcard.com/report/github.com/ShogunPanda/impacca)
[![License](https://img.shields.io/github/license/ShogunPanda/impacca.svg)](https://github.com/ShogunPanda/impacca/blob/master/LICENSE.md)

Package releasing made easy.

http://sw.cowtech.it/impacca

## Installation

Type the following inside a fish shell and you're done!

```bash
curl -sL http://sw.cowtech.it/impacca/installer | sudo bash
```

impicca will be installed in `/usr/local/bin`.

## Usage

impacca is able to help you in maintaining the CHANGELOG.md file and the version for npm modules, Ruby gems and standalone git repository.

It is strongly opinionated, but it should work for most common use cases.

To see all the possible commands, simple run:

```bash
impacca -h
```

## Configuration

impacca tries to find a `.impacca.json` in the current working directory and all its parents and in your home directory.

Here's a list of supported configuration fine (with their default):

```json
{
  "commitMessages": {
    "changelog": "Updated CHANGELOG.md.", // Message used to commit changelog updates.
    "versioning": "Version %s." // Message used to commit version updates. %s will be replaced with the NEW version
  }
}
```

When releasing a new version in plain GIT repository, impacca will also look for `Impaccafile` executable script.
This script will be executed prior commiting the changes and will receive the NEW version as first argument and the OLD version as second argument.
You can find an example of a `Impaccafile` in this repository (which uses this feature).

## Contributing to impacca

- Check out the latest master to make sure the feature hasn't been implemented or the bug hasn't been fixed yet.
- Check out the issue tracker to make sure someone already hasn't requested it and/or contributed it.
- Fork the project.
- Start a feature/bugfix branch.
- Commit and push until you are happy with your contribution.
- Make sure to add tests for it. This is important so I don't break it in a future version unintentionally.

## Copyright

Copyright (C) 2018 and above Shogun (shogun@cowtech.it).

Licensed under the MIT license, which can be found at https://choosealicense.com/licenses/mit.
