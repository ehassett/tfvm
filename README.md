# tfvm

[![Version](https://img.shields.io/github/v/release/ehassett/tfvm)](https://github.com/ehassett/tfvm/releases)
[![MIT License](https://img.shields.io/github/license/ehassett/tfvm)](https://github.com/ehassett/tfvm/blob/main/LICENSE)
[![ci](https://github.com/ehassett/tfvm/actions/workflows/release.yaml/badge.svg)](https://github.com/ehassett/tfvm/actions/workflows/release.yaml)

> A Terraform Version Manager written in Go

## Table of Contents

- [tfvm](#tfvm)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [How it Works](#how-it-works)
  - [Getting Started](#getting-started)
    - [Installation](#installation)
      - [Homebrew (for Mac and Linux)](#homebrew-for-mac-and-linux)
      - [Scoop (for Windows)](#scoop-for-windows)
      - [Script (for Mac and Linux)](#script-for-mac-and-linux)
      - [Go users](#go-users)
    - [CLI Usage](#cli-usage)
  - [Contributing](#contributing)
    - [Development](#development)

## Features

- Easily manage multiple terraform versions to use across projects.
- Run `tfvm use` with no version argument to switch to the version specified in the current directory's `.tfversion` file.
- Works on Linux, Mac, and Windows.

## How it Works

tfvm installs and manages different versions of terraform in the CLI.

## Getting Started
### Installation
#### Homebrew (for Mac and Linux)
Install via [Homebrew](https://brew.sh):
```bash
brew tap ehassett/tfvm
brew install tfvm
```

#### Scoop (for Windows)
Install via [Scoop](https://scoop.sh):
```PowerShell
scoop bucket add tfvm https://github.com/ehassett/tfvm
scoop install tfvm
```

#### Script (for Mac and Linux)
Install via the [install script](install.sh) (requires both curl and wget):
```bash
wget -q -O - https://raw.githubusercontent.com/ehassett/tfvm/master/install.sh | bash
```

#### Go users
Install latest with `go install` (or substitute a version):
```bash
go install github.com/ehassett/tfvm@latest
```

Run `tfvm --version` to verify installation.

| :warning: Important Note                                        |
| :-------------------------------------------------------------- |
| You may need to add `$HOME/.tfvm` to PATH after installing tfvm |

### CLI Usage

```
$ tfvm --help

Usage: tfvm [--version] [--help] <command> [<args>]

Available commands are:
    install    Install a version of Terraform
    list       List all installed versions of Terraform
    remove     Remove a specific version of Terraform
    use        Select a version of Terraform to use
```

## Contributing

Contributions to this project are welcome and much appreciated!

### Development

1. Use Golang version `1.16`.
2. Fork [this repo](https://github.com/ehassett/tfvm).
3. Commit and push your changes, using proper [conventional commit format](https://www.conventionalcommits.org/en/v1.0.0/).
4. Open a Pull Request, rebasing against `master` if needed.

Bugs, feature requests, and comments are more than welcome in the [issues](https://github.com/ehassett/tfvm/issues).
