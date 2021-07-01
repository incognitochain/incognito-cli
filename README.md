[![Go Report Card](https://goreportcard.com/badge/github.com/incognitochain/incognito-cli)](https://goreportcard.com/report/github.com/incognitochain/incognito-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/incognitochain/incognito-cli/blob/main/LICENSE)

incognito-cli
=============
A command line tool for the Incognito network

<!-- toc -->
* [Usage](#usage)
* [Commands](#commands)
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

## Usage
There are two options for you to run the Incognito CLI by:
1. Downloading the pre-compiled executable binary file, you can find it in the [releases](https://github.com/incognitochain/incognito-cli/releases).
2. Compiling your own executable binary file from source as in the Installation instruction above.

Then execute the binary file with the following commands.

```shell
$ incognito-cli [global options] command [command options] [arguments...]
some command is running here...
$ incognito-cli --help
NAME:
   incognito-cli - A new cli application

USAGE:
   incognito-cli [global options] command [command options] [arguments...]

VERSION:
   v0.0.1

DESCRIPTION:
   A simple CLI application for the Incognito network

COMMANDS:
   help, h  Shows a list of commands or help for one command
   account:
     balance           check the balance of an account
     consolidate, csl  consolidate UTXOs of an account
     history, hst      retrieve the history of an account
     keyinfo           print all related-keys of a private key
     utxo              print the UTXOs of an account

GLOBAL OPTIONS:
   --host value                  custom full-node host
   --network value, --net value  network environment (mainnet, testnet, testnet1, devnet, local, custom) (default: "mainnet")
   --help, -h                    show help (default: false)
   --version, -v                 print the version (default: false)

```
<!-- usagestop -->

# Commands
<!-- commands -->
* [`Accounts`](#accounts)
    * [`balance`](#balance)
    * [`keyinfo`](#keyinfo)
    * [`utxo`](#utxo)
    * [`consolidate`](#consolidate)
    * [`history`](#history)
## Accounts
### balance
Display the balance of a private key w.r.t a tokenIDStr.
```shell
$ incognito-cli help balance
NAME:
   incognito-cli balance - check the balance of an account

USAGE:
   balance --privateKey PRIVATE_KEY --tokenID TOKEN_ID

CATEGORY:
   account

OPTIONS:
   --privateKey value, --prvKey value  a base58-encoded private key
   --tokenID value                     ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
```

### keyinfo
Display the all related-key of a private key.
```shell
$ incognito-cli help keyinfo
NAME:
   incognito-cli keyinfo - print all related-keys of a private key

USAGE:
   keyinfo --privateKey PRIVATE_KEY

CATEGORY:
   account

OPTIONS:
   --privateKey value, --prvKey value  a base58-encoded private key
```

### utxo
Print all UTXOs of a private key w.r.t a tokenID.
```shell
$ incognito-cli help utxo
NAME:
   incognito-cli utxo - print the UTXOs of an account

USAGE:
   utxo --privateKey PRIVATE_KEY --tokenID TOKEN_ID

CATEGORY:
   account

OPTIONS:
   --privateKey value, --prvKey value  a base58-encoded private key
   --tokenID value                     ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
```

### consolidate
Consolidate UTXOs of an account.
```shell
$ incognito-cli help consolidate
NAME:
   incognito-cli consolidate - consolidate UTXOs of an account

USAGE:
   consolidate --privateKey PRIVATE_KEY --tokenID TOKEN_ID --version VERSION --numThreads NUM_THREADS --enableLog ENABLE_LOG --logFile LOG_FILE

CATEGORY:
   account

DESCRIPTION:
   This function helps consolidate UTXOs of an account. It consolidates a version of UTXOs at a time, users need to specify which version they need to consolidate. Please note that this process is time-consuming and requires a considerable amount of CPU.

OPTIONS:
   --privateKey value, --prvKey value  a base58-encoded private key
   --tokenID value                     ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   --version value                     version of the UTXOs being converted (1, 2) (default: 1)
   --numThreads value                  number of threads used in this action (default: 4)
   --enableLog                         enable log for this action (default: false)
```

### history
Retrieve the history of an account w.r.t a tokenID.
```shell
$ incognito-cli help history
NAME:
   incognito-cli history - retrieve the history of an account

USAGE:
   history --privateKey PRIVATE_KEY --tokenID TOKEN_ID --numThreads NUM_THREADS --enableLog ENABLE_LOG --logFile LOG_FILE --csvFile CSV_FILE

CATEGORY:
   account

DESCRIPTION:
   This function helps retrieve the history of an account w.r.t a tokenID. Please note that this process is time-consuming and requires a considerable amount of CPU.

OPTIONS:
   --privateKey value, --prvKey value  a base58-encoded private key
   --tokenID value                     ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   --numThreads value                  number of threads used in this action (default: 4)
   --enableLog                         enable log for this action (default: false)
   --logFile value                     location of the log file (default: "os.Stdout")
   --csvFile value                     the csv file location to store the history
```
<!-- commandsstop -->
