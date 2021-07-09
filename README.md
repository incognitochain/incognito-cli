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
$ incognito-cli help
NAME:
   incognito-cli - A simple CLI application for the Incognito network

USAGE:
   incognito-cli [global options] command [command options] [arguments...]

VERSION:
   v0.0.2

DESCRIPTION:
   A simple CLI application for the Incognito network. With this tool, you can run some basic functions on your computer to interact with the Incognito network such as checking balances, transferring PRV or tokens, consolidating and converting your UTXOs, etc.

AUTHOR:
   Incognito Devs Team <support@incognito.org>

COMMANDS:
   help, h  Shows a list of commands or help for one command
   ACCOUNTS:
     balance                  Check the balance of an account.
     consolidate, csl         Consolidate UTXOs of an account.
     generateaccount, genacc  Generate a new Incognito account.
     history, hst             Retrieve the history of an account.
     keyinfo                  Print all related-keys of a private key.
     submitkey, sub           Submit an ota key to the full-node.
     utxo                     Print the UTXOs of an account.
   COMMITTEES:
     checkrewards    Get all rewards of a payment address.
     withdrawreward  Withdraw the reward of a privateKey w.r.t to a tokenID.
   TRANSACTIONS:
     convert  Convert UTXOs of an account w.r.t a tokenID.
     send     Send an amount of PRV or token from one wallet to another wallet.

GLOBAL OPTIONS:
   --host value                  custom full-node host
   --network value, --net value  network environment (mainnet, testnet, testnet1, devnet, local, custom) (default: "mainnet")
   --help, -h                    show help (default: false)
   --version, -v                 print the version (default: false)

COPYRIGHT:
   This tool is developed and maintained by the Incognito Devs Team. It is free for anyone. However, any commercial usages should be acknowledged by the Incognito Devs Team.
```
# Commands
<!-- commands -->
* [`ACCOUNTS`](#accounts)
  * [`balance`](#balance)
  * [`consolidate`](#consolidate)
  * [`generateaccount`](#generateaccount)
  * [`history`](#history)
  * [`keyinfo`](#keyinfo)
  * [`submitkey`](#submitkey)
  * [`utxo`](#utxo)
* [`COMMITTEES`](#committees)
  * [`checkrewards`](#checkrewards)
  * [`withdrawreward`](#withdrawreward)
* [`TRANSACTIONS`](#transactions)
  * [`convert`](#convert)
  * [`send`](#send)
## ACCOUNTS
### balance
Check the balance of an account.
```shell
$ incognito-cli help balance
NAME:
   incognito-cli balance - Check the balance of an account.

USAGE:
   balance --privateKey PRIVATE_KEY [--tokenID TOKEN_ID]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   ACCOUNTS

OPTIONS:
   --privateKey value, --prvKey value  a base58-encoded private key
   --tokenID value                     ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   
```

### consolidate
This function helps consolidate UTXOs of an account. It consolidates a version of UTXOs at a time, users need to specify which version they need to consolidate. Please note that this process is time-consuming and requires a considerable amount of CPU.
```shell
$ incognito-cli help consolidate
NAME:
   incognito-cli consolidate - Consolidate UTXOs of an account.

USAGE:
   consolidate --privateKey PRIVATE_KEY [--tokenID TOKEN_ID] [--version VERSION] [--numThreads NUM_THREADS] [--enableLog ENABLE_LOG] [--logFile LOG_FILE]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   ACCOUNTS

DESCRIPTION:
   This function helps consolidate UTXOs of an account. It consolidates a version of UTXOs at a time, users need to specify which version they need to consolidate. Please note that this process is time-consuming and requires a considerable amount of CPU.

OPTIONS:
   --privateKey value, --prvKey value  a base58-encoded private key
   --tokenID value                     ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   --version value, -v value           version of the transaction (1 or 2) (default: 1)
   --numThreads value                  number of threads used in this action (default: 4)
   --enableLog                         enable log for this action (default: false)
   --logFile value                     location of the log file (default: "os.Stdout")
   
```

### generateaccount
This function helps generate a new mnemonic phrase and its Incognito account.
```shell
$ incognito-cli help generateaccount
NAME:
   incognito-cli generateaccount - Generate a new Incognito account.

USAGE:
   generateaccount

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   ACCOUNTS

DESCRIPTION:
   This function helps generate a new mnemonic phrase and its Incognito account.
```

### history
This function helps retrieve the history of an account w.r.t a tokenID. Please note that this process is time-consuming and requires a considerable amount of CPU.
```shell
$ incognito-cli help history
NAME:
   incognito-cli history - Retrieve the history of an account.

USAGE:
   history --privateKey PRIVATE_KEY [--tokenID TOKEN_ID] [--numThreads NUM_THREADS] [--enableLog ENABLE_LOG] [--logFile LOG_FILE] [--csvFile CSV_FILE]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   ACCOUNTS

DESCRIPTION:
   This function helps retrieve the history of an account w.r.t a tokenID. Please note that this process is time-consuming and requires a considerable amount of CPU.

OPTIONS:
   --privateKey value, --prvKey value  a base58-encoded private key
   --tokenID value                     ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   --numThreads value                  number of threads used in this action (default: 4)
   --enableLog                         enable log for this action (default: false)
   --logFile value                     location of the log file (default: "os.Stdout")
   --csvFile value, --csv value        the csv file location to store the history
   
```

### keyinfo
Print all related-keys of a private key.
```shell
$ incognito-cli help keyinfo
NAME:
   incognito-cli keyinfo - Print all related-keys of a private key.

USAGE:
   keyinfo --privateKey PRIVATE_KEY

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   ACCOUNTS

OPTIONS:
   --privateKey value, --prvKey value  a base58-encoded private key
   
```

### submitkey
This function submits an otaKey to the full-node to use the full-node's cache. If an access token is provided, it will submit the ota key in an authorized manner. See https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/accounts/submit_key.md for more details.
```shell
$ incognito-cli help submitkey
NAME:
   incognito-cli submitkey - Submit an ota key to the full-node.

USAGE:
   submitkey --otaKey OTA_KEY [--accessToken ACCESS_TOKEN] [--fromHeight FROM_HEIGHT] [--isReset IS_RESET]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   ACCOUNTS

DESCRIPTION:
   This function submits an otaKey to the full-node to use the full-node's cache. If an access token is provided, it will submit the ota key in an authorized manner. See https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/accounts/submit_key.md for more details.

OPTIONS:
   --otaKey value, --ota value  a base58-encoded ota key
   --accessToken value          a 64-character long hex-encoded authorized access token
   --fromHeight value           the beacon height at which the full-node will sync from (default: 0)
   --isReset                    whether the full-node should reset the cache for this ota key (default: false)
   
```

### utxo
Print the UTXOs of an account.
```shell
$ incognito-cli help utxo
NAME:
   incognito-cli utxo - Print the UTXOs of an account.

USAGE:
   utxo --privateKey PRIVATE_KEY [--tokenID TOKEN_ID]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   ACCOUNTS

OPTIONS:
   --privateKey value, --prvKey value  a base58-encoded private key
   --tokenID value                     ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   
```

## COMMITTEES
### checkrewards
Get all rewards of a payment address.
```shell
$ incognito-cli help checkrewards
NAME:
   incognito-cli checkrewards - Get all rewards of a payment address.

USAGE:
   checkrewards --address ADDRESS

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   COMMITTEES

OPTIONS:
   --address value, --addr value  a base58-encoded payment address
   
```

### withdrawreward
Withdraw the reward of a privateKey w.r.t to a tokenID.
```shell
$ incognito-cli help withdrawreward
NAME:
   incognito-cli withdrawreward - Withdraw the reward of a privateKey w.r.t to a tokenID.

USAGE:
   withdrawreward --privateKey PRIVATE_KEY [--address ADDRESS] [--tokenID TOKEN_ID] [--version VERSION]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   COMMITTEES

OPTIONS:
   --privateKey value, --prvKey value  a base58-encoded private key
   --address value, --addr value       the payment address of a candidate (default: the payment address of the privateKey)
   --tokenID value                     ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   --version value, -v value           version of the transaction (1 or 2) (default: 1)
   
```

## TRANSACTIONS
### convert
This function helps convert UTXOs v1 of a user to UTXO v2 w.r.t a tokenID. Please note that this process is time-consuming and requires a considerable amount of CPU.
```shell
$ incognito-cli help convert
NAME:
   incognito-cli convert - Convert UTXOs of an account w.r.t a tokenID.

USAGE:
   convert --privateKey PRIVATE_KEY [--tokenID TOKEN_ID] [--numThreads NUM_THREADS] [--enableLog ENABLE_LOG] [--logFile LOG_FILE]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   TRANSACTIONS

DESCRIPTION:
   This function helps convert UTXOs v1 of a user to UTXO v2 w.r.t a tokenID. Please note that this process is time-consuming and requires a considerable amount of CPU.

OPTIONS:
   --privateKey value, --prvKey value  a base58-encoded private key
   --tokenID value                     ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   --numThreads value                  number of threads used in this action (default: 4)
   --enableLog                         enable log for this action (default: false)
   --logFile value                     location of the log file (default: "os.Stdout")
   
```

### send
This function sends an amount of PRV or token from one wallet to another wallet. By default, it used 100 nano PRVs to pay the transaction fee.
```shell
$ incognito-cli help send
NAME:
   incognito-cli send - Send an amount of PRV or token from one wallet to another wallet.

USAGE:
   send --privateKey PRIVATE_KEY --address ADDRESS --amount AMOUNT [--tokenID TOKEN_ID] [--fee FEE] [--version VERSION]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   TRANSACTIONS

DESCRIPTION:
   This function sends an amount of PRV or token from one wallet to another wallet. By default, it used 100 nano PRVs to pay the transaction fee.

OPTIONS:
   --privateKey value, --prvKey value  a base58-encoded private key
   --address value, --addr value       a base58-encoded payment address
   --amount value, --amt value         the amount being transferred (default: 0)
   --tokenID value                     ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   --fee value                         the PRV amount for paying the transaction fee (default: 100)
   --version value, -v value           version of the transaction (1 or 2) (default: 1)
   
```

<!-- commandsstop -->
