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
Install to the `$GOPATH` folder.
```shell
$ go install
```
This command will install the CLI application into your `GOPATH` folder. Alternatively, you can build and install the binary file
into a desired folder by the following command.
```shell
$ go build -o PATH/TO/YOUR/FOLDER/appName
```
If you have issues with these commands, try to clean the golang module cache first.
```shell
go clean --modcache
```


## Initializing environments

An environment is a set of variables you need in your commands. Some important variables:
- `host`: Custom full-node host. This flag is combined with the `network` flag to initialize the environment in which the custom host points to.
- `network`: Network environment (mainnet, testnet, testnet1, local)
- `cache`: Whether to use the UTXO cache (0 - disabled, <> 0 - enabled). See https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/accounts/utxo_cache.md for more information.
- `debug`: Whether to enable the debug mode (0 - disabled, <> 0 - enabled)

To initialize the environment that you want to use for commands, you can use [some flags](#enviroment-flags) with every command at [Commands](./commands.md) or use a CLI [configuration file](#configuration-file).


### Enviroment flags

- `host`: `--host`, e.g. `incognito-cli --host 127.0.0.1:9334`
- `network`: `--network value`, `--net value`, e.g. `incognito-cli --network mainnet`
- `cache`: `--utxoCache value`, `-c value`, `--cache value`, e.g. `incognito-cli --utxoCache 1`
- `debug`: `--debug value`, `-d value`, e.g. `incognito-cli --debug 1`

### Configuration file
The basic syntax of the configuration file is simple:
```
network: "testnet"
host: "http://51.83.36.184:9334"
selfcache: 1
debug: 1
maxthreads: 20
```
This file should be a YAML file, which has the `.yml` extension. To use a configuration file, you can create an `app.yml` file in the current working directory, or use the `fconfig` flag to specify the path where the configuration file resides, e.g. `incognito-cli --fconfig ./app.yml command`



## Usage


See [Commands](./commands.md)
