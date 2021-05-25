# tfvm

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/ethanhassett/tfvm/blob/main/LICENSE)

> A Terraform Version Manager written in Go

## Table of Contents

- [tfvm](#tfvm)
  - [Table of Contents](#table-of-contents)
  - [Changelog](#changelog)
  - [How it Works](#how-it-works)
  - [Getting Started](#getting-started)
    - [Intallation](#installation)
    - [CLI Usage](#cli-usage)
  - [Contributing](#contributing)
    - [Development](#development)
  - [TODO](#todo)

## Changelog

See the [CHANGELOG](https://github.com/ethanhassett/tfvm/blob/main/CHANGELOG.md)

## How it Works

tfvm installs and manages different versions of terraform in the CLI.

## Getting Started
### Installation

Download the appropriate package from [GitHub](https://github.com/ethanhassett/tfvm/releases) and add it to PATH. A proper installation script is in the works.

Run `tfvm` to verify installation.

### CLI Usage

```
$ tfvm help

tfvm usage:
  help
    tfvm help - Shows this help text.
  install
    tfvm install [version] - Installs terraform. If no version is specified, the latest will be installed.
    tfvm install list - Lists the available terraform versions.
  list
    tfvm list - Lists all installed terraform versions. The current version is indicated with a *.
  remove
    tfvm remove <version> - Uninstalls the specified terraform version.
  select
    tfvm select <version> - Selects the specified terraform version to be used.
```

## Contributing

Contributions to this project are welcome and much appreciated!

### Development

1. Use Golang version `>= 1.16`
2. Fork [this repo](https://github.com/ethanhassett/tfvm)
3. Create a `feat:` branch
4. Commit and push your changes
5. Open a Pull Request, rebasing against `main` if needed.

Bugs, feaure requests, and comments are more than welcome in the [issues].

## TODO

- [ ] Add installation script
- [ ] Add pagination to `tfvm install list`