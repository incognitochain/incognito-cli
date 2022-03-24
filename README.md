[![Go Report Card](https://goreportcard.com/badge/github.com/incognitochain/incognito-cli)](https://goreportcard.com/report/github.com/incognitochain/incognito-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/incognitochain/incognito-cli/blob/main/LICENSE)

incognito-cli
=============
A command line tool for the Incognito network

<!-- toc -->
* [Usage](#usage)
* [Commands](./commands.md)
<!-- tocstop -->

# Usage
<!-- usage -->
## Installation
To install, for Linux and macOS users, try the following command:
```shell
$ make linux # or make macos
```
To install with a custom name, try:
```shell
$ bash ./scripts/install_unix.sh -n CUSTOM_APP_NAME -a
```

For Windows user, try:
```shell
$ go install
```
The first two commands will install the CLI application into your `$GOPATH` folder and also enable bash completion, make sure
that you have added `$GOPATH` to the global environment `$PATH`; while the last command will only install the CLI to your `$GOPATH$`.
Alternatively, you can build and install the binary file into a desired folder by the following command.
```shell
$ go build -o PATH/TO/YOUR/FOLDER/appName
```
If you have issues with these commands, try to clean the golang module cache first.
```shell
go clean --modcache
```

## Usage
See [Commands](./commands.md)