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
   v0.0.3

DESCRIPTION:
   A simple CLI application for the Incognito network. With this tool, you can run some basic functions on your computer to interact with the Incognito network such as checking balances, transferring PRV or tokens, consolidating and converting your UTXOs, transferring tokens, manipulating with the pDEX, shielding or un-shielding ETH/BNB/ERC20/BEP20, etc.

AUTHOR:
   Incognito Devs Team

COMMANDS:
   help, h  Shows a list of commands or help for one command
   ACCOUNTS:
     balance                  Check the balance of an account for a tokenID.
     consolidate, csl         Consolidate UTXOs of an account.
     generateaccount, genacc  Generate a new Incognito account.
     history, hst             Retrieve the history of an account.
     importeaccount, import   Import a mnemonic of 12 words.
     keyinfo                  Print all related-keys of a private key.
     outcoin                  Print the output coins of an account.
     submitkey, sub           Submit an ota key to the full-node.
     utxo                     Print the UTXOs of an account.
   BRIDGE:
     evm     Perform an EVM action (e.g, shield, unshield, etc.).
     portal  Perform a portal action (e.g, shield, unshield, etc.).
   COMMITTEES:
     checkrewards    Get all rewards of a payment address.
     stake           Create a staking transaction (https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/staking/stake.md).
     unstake         Create an un-staking transaction (https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/staking/unstake.md).
     withdrawreward  Withdraw the reward of a privateKey w.r.t to a tokenID.
   DEX:
     pdeaction  Perform a pDEX action.
     pdeinfo    Retrieve pDEX information.
     pdestatus  Retrieve the status of a pDEX action.
   TRANSACTIONS:
     checkreceiver  Check if an OTA key is a receiver of a transaction.
     convert        Convert UTXOs of an account w.r.t a tokenID.
     convertall     Convert UTXOs of an account for all assets.
     send           Send an amount of PRV or token from one wallet to another wallet.

GLOBAL OPTIONS:
   --debug value, -d value                     Whether to enable the debug mode (0 - disabled, <> 0 - enabled) (default: 0)
   --host network                              Custom full-node host. This flag is combined with the network flag to initialize the environment in which the custom host points to.
   --network value, --net value                Network environment (mainnet, testnet, testnet1, local) (default: "mainnet")
   --utxoCache value, -c value, --cache value  Whether to use the UTXO cache (0 - disabled, <> 0 - enabled). See https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/accounts/utxo_cache.md for more information. (default: 0)
   --help, -h                                  show help (default: false)
   --version, -v                               print the version (default: false)

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
  * [`importeaccount`](#importeaccount)
  * [`keyinfo`](#keyinfo)
  * [`outcoin`](#outcoin)
  * [`submitkey`](#submitkey)
  * [`utxo`](#utxo)
* [`BRIDGE`](#bridge)
  * [`evm`](#evm)
    * [`evm retryshield`](#evm_retryshield)
    * [`evm retryunshield`](#evm_retryunshield)
    * [`evm shield`](#evm_shield)
    * [`evm unshield`](#evm_unshield)
  * [`portal`](#portal)
    * [`portal shield`](#portal_shield)
    * [`portal shieldaddress`](#portal_shieldaddress)
    * [`portal shieldstatus`](#portal_shieldstatus)
    * [`portal unshield`](#portal_unshield)
    * [`portal unshieldstatus`](#portal_unshieldstatus)
* [`COMMITTEES`](#committees)
  * [`checkrewards`](#checkrewards)
  * [`stake`](#stake)
  * [`unstake`](#unstake)
  * [`withdrawreward`](#withdrawreward)
* [`DEX`](#dex)
  * [`pdeaction`](#pdeaction)
    * [`pdeaction addorder`](#pdeaction_addorder)
    * [`pdeaction contribute`](#pdeaction_contribute)
    * [`pdeaction mintnft`](#pdeaction_mintnft)
    * [`pdeaction stake`](#pdeaction_stake)
    * [`pdeaction trade`](#pdeaction_trade)
    * [`pdeaction unstake`](#pdeaction_unstake)
    * [`pdeaction withdraw`](#pdeaction_withdraw)
    * [`pdeaction withdrawlpfee`](#pdeaction_withdrawlpfee)
    * [`pdeaction withdraworder`](#pdeaction_withdraworder)
    * [`pdeaction withdrawstakereward`](#pdeaction_withdrawstakereward)
  * [`pdeinfo`](#pdeinfo)
    * [`pdeinfo checkprice`](#pdeinfo_checkprice)
    * [`pdeinfo findpath`](#pdeinfo_findpath)
    * [`pdeinfo lpvalue`](#pdeinfo_lpvalue)
    * [`pdeinfo mynft`](#pdeinfo_mynft)
    * [`pdeinfo share`](#pdeinfo_share)
    * [`pdeinfo stakereward`](#pdeinfo_stakereward)
  * [`pdestatus`](#pdestatus)
    * [`pdestatus addorder`](#pdestatus_addorder)
    * [`pdestatus contribute`](#pdestatus_contribute)
    * [`pdestatus mintnft`](#pdestatus_mintnft)
    * [`pdestatus stake`](#pdestatus_stake)
    * [`pdestatus trade`](#pdestatus_trade)
    * [`pdestatus unstake`](#pdestatus_unstake)
    * [`pdestatus withdraw`](#pdestatus_withdraw)
    * [`pdestatus withdrawlpfee`](#pdestatus_withdrawlpfee)
    * [`pdestatus withdraworder`](#pdestatus_withdraworder)
    * [`pdestatus withdrawstakereward`](#pdestatus_withdrawstakereward)
* [`TRANSACTIONS`](#transactions)
  * [`checkreceiver`](#checkreceiver)
  * [`convert`](#convert)
  * [`convertall`](#convertall)
  * [`send`](#send)
## ACCOUNTS
### balance
Check the balance of an account for a tokenID.
```shell
$ incognito-cli help balance
NAME:
   incognito-cli balance - Check the balance of an account for a tokenID.

USAGE:
   balance --privateKey PRIVATE_KEY [--tokenID TOKEN_ID]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   ACCOUNTS

OPTIONS:
   --privateKey value, -p value, --prvKey value  A base58-encoded Incognito private key
   --tokenID value, --id value, --ID value       The Incognito ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   
```

### consolidate
This command helps consolidate UTXOs of an account. It consolidates a version of UTXOs at a time, users need to specify which version they need to consolidate. Please note that this process is time-consuming and requires a considerable amount of CPU.
```shell
$ incognito-cli help consolidate
NAME:
   incognito-cli consolidate - Consolidate UTXOs of an account.

USAGE:
   consolidate --privateKey PRIVATE_KEY [--tokenID TOKEN_ID] [--version VERSION] [--numThreads NUM_THREADS]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   ACCOUNTS

DESCRIPTION:
   This command helps consolidate UTXOs of an account. It consolidates a version of UTXOs at a time, users need to specify which version they need to consolidate. Please note that this process is time-consuming and requires a considerable amount of CPU.

OPTIONS:
   --privateKey value, -p value, --prvKey value  A base58-encoded Incognito private key
   --tokenID value, --id value, --ID value       The Incognito ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   --version value, -v value                     Version of the transaction (1 or 2) (default: 2)
   --numThreads value                            Number of threads used in this action (default: 4)
   
```

### generateaccount
This command helps generate a new mnemonic phrase and its Incognito accounts.
```shell
$ incognito-cli help generateaccount
NAME:
   incognito-cli generateaccount - Generate a new Incognito account.

USAGE:
   generateaccount [--numShards NUM_SHARDS]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   ACCOUNTS

DESCRIPTION:
   This command helps generate a new mnemonic phrase and its Incognito accounts.

OPTIONS:
   --numShards value  The number of shard (default: 8)
   
```

### history
This command helps retrieve the history of an account w.r.t a tokenID. Please note that this process is time-consuming and requires a considerable amount of CPU.
```shell
$ incognito-cli help history
NAME:
   incognito-cli history - Retrieve the history of an account.

USAGE:
   history --privateKey PRIVATE_KEY [--tokenID TOKEN_ID] [--numThreads NUM_THREADS] [--csvFile CSV_FILE]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   ACCOUNTS

DESCRIPTION:
   This command helps retrieve the history of an account w.r.t a tokenID. Please note that this process is time-consuming and requires a considerable amount of CPU.

OPTIONS:
   --privateKey value, -p value, --prvKey value  A base58-encoded Incognito private key
   --tokenID value                               ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   --numThreads value                            Number of threads used in this action (default: 4)
   --csvFile value, --csv value                  The csv file location to store the history
   
```

### importeaccount
This command helps generate Incognito accounts given a mnemonic.
```shell
$ incognito-cli help importeaccount
NAME:
   incognito-cli importeaccount - Import a mnemonic of 12 words.

USAGE:
   importeaccount --mnemonic MNEMONIC [--numShards NUM_SHARDS]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   ACCOUNTS

DESCRIPTION:
   This command helps generate Incognito accounts given a mnemonic.

OPTIONS:
   --mnemonic value, -m value  A 12-word mnemonic phrase, words are separated by a "-" (Example: artist-decline-pepper-spend-good-enemy-caught-sister-sure-opinion-hundred-lake).
   --numShards value           The number of shard (default: 8)
   
```

### keyinfo
Print all related-keys of a private key.
```shell
$ incognito-cli help keyinfo
NAME:
   incognito-cli keyinfo - Print all related-keys of a private key.

USAGE:
   keyinfo --privateKey PRIVATE_KEY

CATEGORY:
   ACCOUNTS

OPTIONS:
   --privateKey value, -p value, --prvKey value  a base58-encoded private key
   
```

### outcoin
Print the output coins of an account.
```shell
$ incognito-cli help outcoin
NAME:
   incognito-cli outcoin - Print the output coins of an account.

USAGE:
   outcoin --address ADDRESS --otaKey OTA_KEY [--readonlyKey READONLY_KEY] [--tokenID TOKEN_ID]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   ACCOUNTS

OPTIONS:
   --address value, --addr value            A base58-encoded payment address
   --otaKey value, --ota value              A base58-encoded ota key
   --readonlyKey value, --ro value          A base58-encoded read-only key
   --tokenID value, --id value, --ID value  The Incognito ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   
```

### submitkey
This command submits an otaKey to the full-node to use the full-node's cache. If an access token is provided, it will submit the ota key in an authorized manner. See https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/accounts/submit_key.md for more details.
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
   This command submits an otaKey to the full-node to use the full-node's cache. If an access token is provided, it will submit the ota key in an authorized manner. See https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/accounts/submit_key.md for more details.

OPTIONS:
   --otaKey value, --ota value  A base58-encoded ota key
   --accessToken value          A 64-character long hex-encoded authorized access token
   --fromHeight value           The beacon height at which the full-node will sync from (default: 0)
   --isReset                    Whether the full-node should reset the cache for this ota key (default: false)
   
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
   --privateKey value, -p value, --prvKey value  A base58-encoded Incognito private key
   --tokenID value, --id value, --ID value       The Incognito ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   
```

## BRIDGE
### evm
This command helps perform an EVM action (e.g, shield, unshield, etc.).
```shell
$ incognito-cli help evm
NAME:
   incognito-cli evm - Perform an EVM action (e.g, shield, unshield, etc.).

USAGE:
   evm

CATEGORY:
   BRIDGE

DESCRIPTION:
   This command helps perform an EVM action (e.g, shield, unshield, etc.).
```

#### evm_retryshield
This command re-shields an already-been-deposited-to-sc transaction in case of prior failure.
```shell
$ incognito-cli evm help retryshield
NAME:
   incognito-cli evm retryshield - Retry a shield from the given already-been-deposited-to-sc EVM transaction.

USAGE:
   evm retryshield --privateKey PRIVATE_KEY --externalTxHash EXTERNAL_TX_HASH [--evm EVM] [--externalTokenAddress EXTERNAL_TOKEN_ADDRESS]

   OPTIONAL flags are denoted by a [] bracket.

DESCRIPTION:
   This command re-shields an already-been-deposited-to-sc transaction in case of prior failure.

OPTIONS:
   --privateKey value, -p value, --prvKey value  A base58-encoded Incognito private key
   --externalTxHash value, --eTxID value         The external transaction hash
   --evm value                                   The EVM network (ETH or BSC) (default: "ETH")
   --externalTokenAddress value                  ID of the token on ETH/BSC networks (default: "0x0000000000000000000000000000000000000000")
   
```

#### evm_retryunshield
This command tries to un-shield an asset from an already-been-burned Incognito transaction in case of prior failure.
```shell
$ incognito-cli evm help retryunshield
NAME:
   incognito-cli evm retryunshield - Retry an un-shielding request from the given already-been-burned Incognito transaction.

USAGE:
   evm retryunshield --txHash TX_HASH [--evm EVM]

   OPTIONAL flags are denoted by a [] bracket.

DESCRIPTION:
   This command tries to un-shield an asset from an already-been-burned Incognito transaction in case of prior failure.

OPTIONS:
   --txHash value, --iTxID value  An Incognito transaction hash
   --evm value                    The EVM network (ETH or BSC) (default: "ETH")
   
```

#### evm_shield
This function helps shield an EVM (ETH/BNB/ERC20/BEP20) token into the Incognito network. It will ask for users' EVM PRIVATE KEY to proceed. The shielding process consists of the following operations.
1. Deposit the EVM asset into the corresponding smart contract.
1.1. In case the asset is an ERC20/BEP20 token, an approval transaction is performed (if needed) the before the actual deposit. For this operation, a prompt will be displayed to ask for user's approval.
2. Get the deposited EVM transaction, parse the depositing proof and submit it to the Incognito network. This step requires an Incognito private key with a sufficient amount of PRV to create an issuing transaction.

Note that EVM shielding is a complicated process, users MUST understand how the process works before using this function. We RECOMMEND users test the function with test networks BEFORE performing it on the live networks.
DO NOT USE THIS FUNCTION UNLESS YOU UNDERSTAND THE SHIELDING PROCESS.
```shell
$ incognito-cli evm help shield
NAME:
   incognito-cli evm shield - Shield an EVM (ETH/BNB/ERC20/BEP20) token into the Incognito network.

USAGE:
   evm shield --privateKey PRIVATE_KEY --shieldAmount SHIELD_AMOUNT [--evm EVM] [--externalTokenAddress EXTERNAL_TOKEN_ADDRESS] [--address ADDRESS]

   OPTIONAL flags are denoted by a [] bracket.

DESCRIPTION:
   This function helps shield an EVM (ETH/BNB/ERC20/BEP20) token into the Incognito network. It will ask for users' EVM PRIVATE KEY to proceed. The shielding process consists of the following operations.
      1. Deposit the EVM asset into the corresponding smart contract.
        1.1. In case the asset is an ERC20/BEP20 token, an approval transaction is performed (if needed) the before the actual deposit. For this operation, a prompt will be displayed to ask for user's approval.
      2. Get the deposited EVM transaction, parse the depositing proof and submit it to the Incognito network. This step requires an Incognito private key with a sufficient amount of PRV to create an issuing transaction.
   
   Note that EVM shielding is a complicated process, users MUST understand how the process works before using this function. We RECOMMEND users test the function with test networks BEFORE performing it on the live networks.
   DO NOT USE THIS FUNCTION UNLESS YOU UNDERSTAND THE SHIELDING PROCESS.

OPTIONS:
   --privateKey value, -p value, --prvKey value  A base58-encoded Incognito private key
   --shieldAmount value, --amt value             The shielding amount measured in token unit (e.g, 10, 1, 0.1, 0.01) (default: 0)
   --evm value                                   The EVM network (ETH or BSC) (default: "ETH")
   --externalTokenAddress value                  ID of the token on ETH/BSC networks (default: "0x0000000000000000000000000000000000000000")
   --address value, --addr value                 The Incognito payment address to receive the shielding asset (default: the payment address of the privateKey)
   
```

#### evm_unshield
This function helps withdraw an EVM (ETH/BNB/ERC20/BEP20) token out of the Incognito network.The un-shielding process consists the following operations.
1. Users burn the token inside the Incognito chain.
2. After the burning is success, wait for 1-2 Incognito blocks and retrieve the corresponding burn proof from the Incognito chain.
3. After successfully retrieving the burn proof, users submit the burn proof to the smart contract to get back the corresponding public token. This step will ask for users' EVM PRIVATE KEY to proceed. Note that ONLY UNTIL this step, it is feasible to estimate the actual un-shielding fee (mainly is the fee interacting with the smart contract).

Please be aware that EVM un-shielding is a complicated process; and once burned, there is NO WAY to recover the asset inside the Incognito network. Therefore, use this function IF ADN ONLY IF you understand the way un-shielding works. Otherwise, use the un-shielding function from the Incognito app. We RECOMMEND users test the function with test networks BEFORE performing it on the live networks.
DO NOT USE THIS FUNCTION UNLESS YOU UNDERSTAND THE UN-SHIELDING PROCESS.
```shell
$ incognito-cli evm help unshield
NAME:
   incognito-cli evm unshield - Withdraw an EVM (ETH/BNB/ERC20/BEP20) token from the Incognito network.

USAGE:
   evm unshield --privateKey PRIVATE_KEY --tokenID TOKEN_ID --amount AMOUNT

DESCRIPTION:
   This function helps withdraw an EVM (ETH/BNB/ERC20/BEP20) token out of the Incognito network.The un-shielding process consists the following operations.
      1. Users burn the token inside the Incognito chain.
      2. After the burning is success, wait for 1-2 Incognito blocks and retrieve the corresponding burn proof from the Incognito chain.
      3. After successfully retrieving the burn proof, users submit the burn proof to the smart contract to get back the corresponding public token. This step will ask for users' EVM PRIVATE KEY to proceed. Note that ONLY UNTIL this step, it is feasible to estimate the actual un-shielding fee (mainly is the fee interacting with the smart contract).
   
   Please be aware that EVM un-shielding is a complicated process; and once burned, there is NO WAY to recover the asset inside the Incognito network. Therefore, use this function IF ADN ONLY IF you understand the way un-shielding works. Otherwise, use the un-shielding function from the Incognito app. We RECOMMEND users test the function with test networks BEFORE performing it on the live networks.
   DO NOT USE THIS FUNCTION UNLESS YOU UNDERSTAND THE UN-SHIELDING PROCESS.

OPTIONS:
   --privateKey value, -p value, --prvKey value  A base58-encoded Incognito private key
   --tokenID value, --id value, --ID value       The Incognito tokenID of the un-shielding asset
   --amount value, --amt value                   The Incognito amount of the action (default: 0)
   
```

### portal
This command helps perform a portal action (e.g, shield, unshield, etc.).
```shell
$ incognito-cli help portal
NAME:
   incognito-cli portal - Perform a portal action (e.g, shield, unshield, etc.).

USAGE:
   portal

CATEGORY:
   BRIDGE

DESCRIPTION:
   This command helps perform a portal action (e.g, shield, unshield, etc.).
```

#### portal_shield
This command helps shield a portal token into the Incognito network after the fund has been transferred to the depositing address (generated by `portalshieldaddress`).
```shell
$ incognito-cli portal help shield
NAME:
   incognito-cli portal shield - Shield a portal token (e.g, BTC) into the Incognito network.

USAGE:
   portal shield --privateKey PRIVATE_KEY --externalTxHash EXTERNAL_TX_HASH [--tokenID TOKEN_ID] [--address ADDRESS]

   OPTIONAL flags are denoted by a [] bracket.

DESCRIPTION:
   This command helps shield a portal token into the Incognito network after the fund has been transferred to the depositing address (generated by `portalshieldaddress`).

OPTIONS:
   --privateKey value, -p value, --prvKey value  A base58-encoded Incognito private key
   --externalTxHash value, --eTxID value         The external transaction hash
   --tokenID value, --id value, --ID value       The Incognito tokenID of the shielding asset (default: "b832e5d3b1f01a4f0623f7fe91d6673461e1f5d37d91fe78c5c2e6183ff39696")
   --address value, --addr value                 The Incognito payment address to receive the shielding asset (default: the payment address of the privateKey)
   
```

#### portal_shieldaddress
This command helps generate the portal shielding address for a payment address and a tokenID.
```shell
$ incognito-cli portal help shieldaddress
NAME:
   incognito-cli portal shieldaddress - Generate a portal shielding address.

USAGE:
   portal shieldaddress --address ADDRESS [--tokenID TOKEN_ID]

   OPTIONAL flags are denoted by a [] bracket.

DESCRIPTION:
   This command helps generate the portal shielding address for a payment address and a tokenID.

OPTIONS:
   --address value, --addr value            A base58-encoded payment address
   --tokenID value, --id value, --ID value  The Incognito tokenID of the shielding asset (default: "b832e5d3b1f01a4f0623f7fe91d6673461e1f5d37d91fe78c5c2e6183ff39696")
   
```

#### portal_shieldstatus
This command helps retrieve the status of a portal shielding request.
Status should be understood as: 0 - rejected; 1 - accepted.
If you encounter an error, it might be because the request hasn't reached the beacon chain yet. Please try again a few minutes later.
```shell
$ incognito-cli portal help shieldstatus
NAME:
   incognito-cli portal shieldstatus - Get the status of a portal shielding request.

USAGE:
   portal shieldstatus --txHash TX_HASH

DESCRIPTION:
   This command helps retrieve the status of a portal shielding request.
   Status should be understood as: 0 - rejected; 1 - accepted.
   If you encounter an error, it might be because the request hasn't reached the beacon chain yet. Please try again a few minutes later.

OPTIONS:
   --txHash value, --iTxID value  An Incognito transaction hash
   
```

#### portal_unshield
This command helps withdraw portal tokens (BTC) out of the Incognito network.
```shell
$ incognito-cli portal help unshield
NAME:
   incognito-cli portal unshield - Withdraw portal tokens (BTC) from the Incognito network.

USAGE:
   portal unshield --privateKey PRIVATE_KEY --externalAddress EXTERNAL_ADDRESS --amount AMOUNT [--tokenID TOKEN_ID]

   OPTIONAL flags are denoted by a [] bracket.

DESCRIPTION:
   This command helps withdraw portal tokens (BTC) out of the Incognito network.

OPTIONS:
   --privateKey value, -p value, --prvKey value  A base58-encoded Incognito private key
   --externalAddress value, --eAddr value        A valid remote address for the currently-processed tokenID. User MUST make sure this address is valid to avoid the loss of money.
   --amount value, --amt value                   The Incognito amount of the action (default: 0)
   --tokenID value, --id value, --ID value       The Incognito tokenID of the un-shielding asset (default: "b832e5d3b1f01a4f0623f7fe91d6673461e1f5d37d91fe78c5c2e6183ff39696")
   
```

#### portal_unshieldstatus
This command helps retrieve the status of a portal un-shielding request.
Status should be understood as: 0 - waiting; 1 - processed but not completed; 2 - completed; 3 - rejected.
If you encounter an error saying "unexpected end of JSON input", it might be because the request hasn't reached the beacon chain yet. Please try again a few minutes later.
```shell
$ incognito-cli portal help unshieldstatus
NAME:
   incognito-cli portal unshieldstatus - Get the status of a portal un-shielding request.

USAGE:
   portal unshieldstatus --txHash TX_HASH

DESCRIPTION:
   This command helps retrieve the status of a portal un-shielding request.
   Status should be understood as: 0 - waiting; 1 - processed but not completed; 2 - completed; 3 - rejected.
   If you encounter an error saying "unexpected end of JSON input", it might be because the request hasn't reached the beacon chain yet. Please try again a few minutes later.

OPTIONS:
   --txHash value, --iTxID value  An Incognito transaction hash
   
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

CATEGORY:
   COMMITTEES

OPTIONS:
   --address value, --addr value  A base58-encoded payment address
   
```

### stake
Create a staking transaction (https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/staking/stake.md).
```shell
$ incognito-cli help stake
NAME:
   incognito-cli stake - Create a staking transaction (https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/staking/stake.md).

USAGE:
   stake --privateKey PRIVATE_KEY [--miningKey MINING_KEY] [--candidateAddress CANDIDATE_ADDRESS] [--rewardAddress REWARD_ADDRESS] [--autoReStake AUTO_RE_STAKE]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   COMMITTEES

OPTIONS:
   --privateKey value, -p value, --prvKey value   A base58-encoded Incognito private key
   --miningKey value, --mKey value, --vKey value  An Incognito mining key of the committee candidate (default: the mining key associated with the privateKey)
   --candidateAddress value, --canAddr value      The Incognito payment address of the committee candidate (default: the payment address of the privateKey)
   --rewardAddress value, --rwdAddr value         The Incognito payment address of the reward receiver (default: the payment address of the privateKey)
   --autoReStake value, --reStake value           Whether or not to automatically re-stake (0 - false, <> 0 - true) (default: 1)
   
```

### unstake
Create an un-staking transaction (https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/staking/unstake.md).
```shell
$ incognito-cli help unstake
NAME:
   incognito-cli unstake - Create an un-staking transaction (https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/staking/unstake.md).

USAGE:
   unstake --privateKey PRIVATE_KEY [--miningKey MINING_KEY] [--candidateAddress CANDIDATE_ADDRESS]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   COMMITTEES

OPTIONS:
   --privateKey value, -p value, --prvKey value   A base58-encoded Incognito private key
   --miningKey value, --mKey value, --vKey value  An Incognito mining key of the committee candidate (default: the mining key associated with the privateKey)
   --candidateAddress value, --canAddr value      The Incognito payment address of the committee candidate (default: the payment address of the privateKey)
   
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
   --privateKey value, -p value, --prvKey value  A base58-encoded Incognito private key
   --address value, --addr value                 the payment address of a candidate (default: the payment address of the privateKey)
   --tokenID value, --id value, --ID value       The Incognito ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   --version value, -v value                     Version of the transaction (1 or 2) (default: 2)
   
```

## DEX
### pdeaction
This command helps perform a pDEX action.
```shell
$ incognito-cli help pdeaction
NAME:
   incognito-cli pdeaction - Perform a pDEX action.

USAGE:
   pdeaction

CATEGORY:
   DEX

DESCRIPTION:
   This command helps perform a pDEX action.
```

#### pdeaction_addorder
This command creates a transaction adding an order to the pDEX.
```shell
$ incognito-cli pdeaction help addorder
NAME:
   incognito-cli pdeaction addorder - Add an order book to the pDEX.

USAGE:
   pdeaction addorder --privateKey PRIVATE_KEY --pairID PAIR_ID --nftID NFT_ID --sellTokenID SELL_TOKEN_ID --buyTokenID BUY_TOKEN_ID --sellingAmount SELLING_AMOUNT [--minAcceptAmount MIN_ACCEPT_AMOUNT]

   OPTIONAL flags are denoted by a [] bracket.

DESCRIPTION:
   This command creates a transaction adding an order to the pDEX.

OPTIONS:
   --privateKey value, -p value, --prvKey value         A base58-encoded Incognito private key
   --pairID value, --pairId value                       The ID of the target pool pair
   --nftID value, --nftId value                         A pDEX NFT generated by the nft minting command
   --sellTokenID value, --sellID value, --sellId value  ID of the token to sell
   --buyTokenID value, --buyID value, --buyId value     ID of the token to buy
   --sellingAmount value, --sellAmt value               The amount of sellTokenID wished to sell (default: 0)
   --minAcceptAmount value, --minAmt value              The minimum acceptable amount of buyTokenID wished to receive (default: 0)
   
```

#### pdeaction_contribute
This command creates a pDEX liquidity-contributing transaction. See more about this transaction: https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/pdex/contribute.md
```shell
$ incognito-cli pdeaction help contribute
NAME:
   incognito-cli pdeaction contribute - Create a pDEX liquidity-contributing transaction.

USAGE:
   pdeaction contribute --privateKey PRIVATE_KEY --nftID NFT_ID --pairHash PAIR_HASH --amount AMOUNT --amplifier AMPLIFIER [--tokenID TOKEN_ID] [--pairID PAIR_ID]

   OPTIONAL flags are denoted by a [] bracket.

DESCRIPTION:
   This command creates a pDEX liquidity-contributing transaction. See more about this transaction: https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/pdex/contribute.md

OPTIONS:
   --privateKey value, -p value, --prvKey value  A base58-encoded Incognito private key
   --nftID value, --nftId value                  A pDEX NFT generated by the nft minting command
   --pairHash value                              A unique string representing the contributing pair
   --amount value, --amt value                   The Incognito amount of the action (default: 0)
   --amplifier value, --amp value                The amplifier for the target contributing pool (default: 0)
   --tokenID value, --id value, --ID value       The Incognito ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   --pairID value                                The ID of the contributing pool pair. For pool-initializing transactions (e.g, first contribution in the pool), it should be left empty.
   
```

#### pdeaction_mintnft
This command creates and broadcasts a transaction that mints a new (pDEX) NFT for the pDEX.
```shell
$ incognito-cli pdeaction help mintnft
NAME:
   incognito-cli pdeaction mintnft - Create a (pDEX) NFT minting transaction.

USAGE:
   pdeaction mintnft --privateKey PRIVATE_KEY

DESCRIPTION:
   This command creates and broadcasts a transaction that mints a new (pDEX) NFT for the pDEX.

OPTIONS:
   --privateKey value, -p value, --prvKey value  A base58-encoded Incognito private key
   
```

#### pdeaction_stake
This command creates a transaction staking a token to the pDEX.
```shell
$ incognito-cli pdeaction help stake
NAME:
   incognito-cli pdeaction stake - Stake a token to the pDEX.

USAGE:
   pdeaction stake --privateKey PRIVATE_KEY --nftID NFT_ID --amount AMOUNT [--tokenID TOKEN_ID]

   OPTIONAL flags are denoted by a [] bracket.

DESCRIPTION:
   This command creates a transaction staking a token to the pDEX.

OPTIONS:
   --privateKey value, -p value, --prvKey value  A base58-encoded Incognito private key
   --nftID value, --nftId value                  A pDEX NFT generated by the nft minting command
   --amount value, --amt value                   The Incognito amount of the action (default: 0)
   --tokenID value                               The ID of the target staking pool ID (or token ID) (default: "0000000000000000000000000000000000000000000000000000000000000004")
   
```

#### pdeaction_trade
This command creates a trade transaction on the pDEX.
```shell
$ incognito-cli pdeaction help trade
NAME:
   incognito-cli pdeaction trade - Create a trade transaction.

USAGE:
   pdeaction trade --privateKey PRIVATE_KEY --sellTokenID SELL_TOKEN_ID --buyTokenID BUY_TOKEN_ID --sellingAmount SELLING_AMOUNT --tradingFee TRADING_FEE [--minAcceptAmount MIN_ACCEPT_AMOUNT] [--tradingPath TRADING_PATH] [--prvFee PRV_FEE] [--maxPaths MAX_PATHS]

   OPTIONAL flags are denoted by a [] bracket.

DESCRIPTION:
   This command creates a trade transaction on the pDEX.

OPTIONS:
   --privateKey value, -p value, --prvKey value         A base58-encoded Incognito private key
   --sellTokenID value, --sellID value, --sellId value  ID of the token to sell
   --buyTokenID value, --buyID value, --buyId value     ID of the token to buy
   --sellingAmount value, --sellAmt value               The amount of sellTokenID wished to sell (default: 0)
   --tradingFee value                                   The trading fee (default: 0)
   --minAcceptAmount value, --minAmt value              The minimum acceptable amount of buyTokenID wished to receive (default: 0)
   --tradingPath pairID1,pairID2                        A list of trading pair IDs seperated by a comma (Example: pairID1,pairID2). If none is given, the tool will automatically find a suitable path.
   --prvFee value                                       Whether or not to pay fee in PRV (0 - no, <> 0 - yes) (default: 1)
   --maxPaths value                                     The maximum length of the trading path. (default: 5)
   
```

#### pdeaction_unstake
This command creates a transaction un-staking a token from the pDEX.
```shell
$ incognito-cli pdeaction help unstake
NAME:
   incognito-cli pdeaction unstake - Un-stake a token from the pDEX.

USAGE:
   pdeaction unstake --privateKey PRIVATE_KEY --nftID NFT_ID --amount AMOUNT [--tokenID TOKEN_ID]

   OPTIONAL flags are denoted by a [] bracket.

DESCRIPTION:
   This command creates a transaction un-staking a token from the pDEX.

OPTIONS:
   --privateKey value, -p value, --prvKey value  A base58-encoded Incognito private key
   --nftID value, --nftId value                  A pDEX NFT generated by the nft minting command
   --amount value, --amt value                   The Incognito amount of the action (default: 0)
   --tokenID value                               The ID of the target staking pool ID (or token ID) (default: "0000000000000000000000000000000000000000000000000000000000000004")
   
```

#### pdeaction_withdraw
This command creates a transaction withdrawing an amount of `share` from the pDEX. See more about this transaction: https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/pdex/withdrawal.md
```shell
$ incognito-cli pdeaction help withdraw
NAME:
   incognito-cli pdeaction withdraw - Create a pDEX liquidity-withdrawal transaction.

USAGE:
   pdeaction withdraw --privateKey PRIVATE_KEY --pairID PAIR_ID --nftID NFT_ID --tokenID1 TOKEN_ID_1 [--tokenID2 TOKEN_ID_2] [--amount AMOUNT]

   OPTIONAL flags are denoted by a [] bracket.

DESCRIPTION:
   This command creates a transaction withdrawing an amount of `share` from the pDEX. See more about this transaction: https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/pdex/withdrawal.md

OPTIONS:
   --privateKey value, -p value, --prvKey value  A base58-encoded Incognito private key
   --pairID value                                The ID of the contributed pool pair
   --nftID value, --nftId value                  A pDEX NFT generated by the nft minting command
   --tokenID1 value, --id1 value, --ID1 value    ID of the first token
   --tokenID2 value, --id2 value, --ID2 value    ID of the second token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   --amount value, --amt value                   The amount of share wished to withdraw. If set to 0, it will withdraw all of the share. (default: 0)
   
```

#### pdeaction_withdrawlpfee
This command creates a transaction withdrawing LP fees from the pDEX.
```shell
$ incognito-cli pdeaction help withdrawlpfee
NAME:
   incognito-cli pdeaction withdrawlpfee - Withdraw LP fees from the pDEX.

USAGE:
   pdeaction withdrawlpfee --privateKey PRIVATE_KEY --pairID PAIR_ID --nftID NFT_ID

DESCRIPTION:
   This command creates a transaction withdrawing LP fees from the pDEX.

OPTIONS:
   --privateKey value, -p value, --prvKey value  A base58-encoded Incognito private key
   --pairID value, --pairId value                The ID of the target pool pair
   --nftID value, --nftId value                  A pDEX NFT generated by the nft minting command
   
```

#### pdeaction_withdraworder
This command creates a transaction withdrawing an order to the pDEX.
```shell
$ incognito-cli pdeaction help withdraworder
NAME:
   incognito-cli pdeaction withdraworder - Withdraw an order from the pDEX.

USAGE:
   pdeaction withdraworder --privateKey PRIVATE_KEY --orderID ORDER_ID --pairID PAIR_ID --nftID NFT_ID --amount AMOUNT --tokenID1 TOKEN_ID_1 [--tokenID2 TOKEN_ID_2]

   OPTIONAL flags are denoted by a [] bracket.

DESCRIPTION:
   This command creates a transaction withdrawing an order to the pDEX.

OPTIONS:
   --privateKey value, -p value, --prvKey value  A base58-encoded Incognito private key
   --orderID value, --orderId value              The ID of the order.
   --pairID value, --pairId value                The ID of the target pool pair
   --nftID value, --nftId value                  A pDEX NFT generated by the nft minting command
   --amount value, --amt value                   The Incognito amount of the action (default: 0)
   --tokenID1 value, --id1 value, --ID1 value    ID of the first token
   --tokenID2 value, --id2 value, --ID2 value    ID of the second token (if have). In the case of withdrawing a single token, leave it empty.
   
```

#### pdeaction_withdrawstakereward
This command creates a transaction withdrawing staking rewards from the pDEX.
```shell
$ incognito-cli pdeaction help withdrawstakereward
NAME:
   incognito-cli pdeaction withdrawstakereward - Withdraw staking rewards from the pDEX.

USAGE:
   pdeaction withdrawstakereward --privateKey PRIVATE_KEY --nftID NFT_ID [--tokenID TOKEN_ID]

   OPTIONAL flags are denoted by a [] bracket.

DESCRIPTION:
   This command creates a transaction withdrawing staking rewards from the pDEX.

OPTIONS:
   --privateKey value, -p value, --prvKey value  A base58-encoded Incognito private key
   --nftID value, --nftId value                  A pDEX NFT generated by the nft minting command
   --tokenID value                               The ID of the target staking pool ID (or token ID) (default: "0000000000000000000000000000000000000000000000000000000000000004")
   
```

### pdeinfo
This command helps retrieve some information of the pDEX.
```shell
$ incognito-cli help pdeinfo
NAME:
   incognito-cli pdeinfo - Retrieve pDEX information.

USAGE:
   pdeinfo

CATEGORY:
   DEX

DESCRIPTION:
   This command helps retrieve some information of the pDEX.
```

#### pdeinfo_checkprice
This command checks the price of a pair of tokenIds. It must be supplied with the selling amount since the pDEX uses the AMM algorithm.
```shell
$ incognito-cli pdeinfo help checkprice
NAME:
   incognito-cli pdeinfo checkprice - Check the price between two tokenIDs.

USAGE:
   pdeinfo checkprice --sellTokenID SELL_TOKEN_ID --buyTokenID BUY_TOKEN_ID --sellingAmount SELLING_AMOUNT --pairID PAIR_ID

DESCRIPTION:
   This command checks the price of a pair of tokenIds. It must be supplied with the selling amount since the pDEX uses the AMM algorithm.

OPTIONS:
   --sellTokenID value, --sellID value, --sellId value  ID of the token to sell
   --buyTokenID value, --buyID value, --buyId value     ID of the token to buy
   --sellingAmount value, --sellAmt value               The amount of sellTokenID wished to sell (default: 0)
   --pairID value                                       The ID of the target pool pair
   
```

#### pdeinfo_findpath
This command helps find a good trading path for a trade.
```shell
$ incognito-cli pdeinfo help findpath
NAME:
   incognito-cli pdeinfo findpath - Find a `good` trading path for a trade.

USAGE:
   pdeinfo findpath --sellTokenID SELL_TOKEN_ID --buyTokenID BUY_TOKEN_ID --sellingAmount SELLING_AMOUNT [--maxPaths MAX_PATHS]

   OPTIONAL flags are denoted by a [] bracket.

DESCRIPTION:
   This command helps find a good trading path for a trade.

OPTIONS:
   --sellTokenID value, --sellID value, --sellId value  ID of the token to sell
   --buyTokenID value, --buyID value, --buyId value     ID of the token to buy
   --sellingAmount value, --sellAmt value               The amount of sellTokenID wished to sell (default: 0)
   --maxPaths value                                     The maximum length of the trading path. (default: 5)
   
```

#### pdeinfo_lpvalue
This command retrieves the information about the value of an LP in a given pool.
```shell
$ incognito-cli pdeinfo help lpvalue
NAME:
   incognito-cli pdeinfo lpvalue - Check the estimated LP value in a given pool.

USAGE:
   pdeinfo lpvalue --pairID PAIR_ID --nftID NFT_ID

DESCRIPTION:
   This command retrieves the information about the value of an LP in a given pool.

OPTIONS:
   --pairID value, --pairId value  The ID of the target pool pair
   --nftID value, --nftId value    A pDEX NFT generated by the nft minting command
   
```

#### pdeinfo_mynft
This command returns the list of NFTs for a given private key.
```shell
$ incognito-cli pdeinfo help mynft
NAME:
   incognito-cli pdeinfo mynft - Retrieve the list of NFTs for a given private key.

USAGE:
   pdeinfo mynft --privateKey PRIVATE_KEY

DESCRIPTION:
   This command returns the list of NFTs for a given private key.

OPTIONS:
   --privateKey value, -p value, --prvKey value  A base58-encoded Incognito private key
   
```

#### pdeinfo_share
This command returns the share amount of an nftID within a pDEX poolID.
```shell
$ incognito-cli pdeinfo help share
NAME:
   incognito-cli pdeinfo share - Retrieve the share amount of a pDEX poolID given an nftID.

USAGE:
   pdeinfo share --pairID PAIR_ID --nftID NFT_ID

DESCRIPTION:
   This command returns the share amount of an nftID within a pDEX poolID.

OPTIONS:
   --pairID value                The ID of the target pool pair
   --nftID value, --nftId value  A pDEX NFT generated by the nft minting command
   
```

#### pdeinfo_stakereward
This command returns the estimated pDEX staking rewards of an nftID within a pDEX staking pool.
```shell
$ incognito-cli pdeinfo help stakereward
NAME:
   incognito-cli pdeinfo stakereward - Retrieve the estimated pDEX staking rewards.

USAGE:
   pdeinfo stakereward --nftID NFT_ID [--tokenID TOKEN_ID]

   OPTIONAL flags are denoted by a [] bracket.

DESCRIPTION:
   This command returns the estimated pDEX staking rewards of an nftID within a pDEX staking pool.

OPTIONS:
   --nftID value, --nftId value  A pDEX NFT generated by the nft minting command
   --tokenID value               The ID of the target staking pool ID (or token ID) (default: "0000000000000000000000000000000000000000000000000000000000000004")
   
```

### pdestatus
This command helps retrieve the status of a pDEX action given its hash. If an error is thrown, it is mainly because the transaction has not yet reached the beacon chain or the txHash is invalid.
```shell
$ incognito-cli help pdestatus
NAME:
   incognito-cli pdestatus - Retrieve the status of a pDEX action.

USAGE:
   pdestatus

CATEGORY:
   DEX

DESCRIPTION:
   This command helps retrieve the status of a pDEX action given its hash. If an error is thrown, it is mainly because the transaction has not yet reached the beacon chain or the txHash is invalid.
```

#### pdestatus_addorder
Check the status of a pDEX order-adding withdrawal.
```shell
$ incognito-cli pdestatus help addorder
NAME:
   incognito-cli pdestatus addorder - Check the status of a pDEX order-adding withdrawal.

USAGE:
   pdestatus addorder --txHash TX_HASH

OPTIONS:
   --txHash value, --iTxID value  An Incognito transaction hash
   
```

#### pdestatus_contribute
Check the status of a pDEX liquidity contribution.
```shell
$ incognito-cli pdestatus help contribute
NAME:
   incognito-cli pdestatus contribute - Check the status of a pDEX liquidity contribution.

USAGE:
   pdestatus contribute --txHash TX_HASH

OPTIONS:
   --txHash value, --iTxID value  An Incognito transaction hash
   
```

#### pdestatus_mintnft
Check the status of a (pDEX) NFT minting transaction.
```shell
$ incognito-cli pdestatus help mintnft
NAME:
   incognito-cli pdestatus mintnft - Check the status of a (pDEX) NFT minting transaction.

USAGE:
   pdestatus mintnft --txHash TX_HASH

OPTIONS:
   --txHash value, --iTxID value  An Incognito transaction hash
   
```

#### pdestatus_stake
Check the status of a pDEX staking transaction.
```shell
$ incognito-cli pdestatus help stake
NAME:
   incognito-cli pdestatus stake - Check the status of a pDEX staking transaction.

USAGE:
   pdestatus stake --txHash TX_HASH

OPTIONS:
   --txHash value, --iTxID value  An Incognito transaction hash
   
```

#### pdestatus_trade
Check the status of a pDEX trade.
```shell
$ incognito-cli pdestatus help trade
NAME:
   incognito-cli pdestatus trade - Check the status of a pDEX trade.

USAGE:
   pdestatus trade --txHash TX_HASH

OPTIONS:
   --txHash value, --iTxID value  An Incognito transaction hash
   
```

#### pdestatus_unstake
Check the status of a pDEX un-staking transaction.
```shell
$ incognito-cli pdestatus help unstake
NAME:
   incognito-cli pdestatus unstake - Check the status of a pDEX un-staking transaction.

USAGE:
   pdestatus unstake --txHash TX_HASH

OPTIONS:
   --txHash value, --iTxID value  An Incognito transaction hash
   
```

#### pdestatus_withdraw
Check the status of a pDEX liquidity withdrawal.
```shell
$ incognito-cli pdestatus help withdraw
NAME:
   incognito-cli pdestatus withdraw - Check the status of a pDEX liquidity withdrawal.

USAGE:
   pdestatus withdraw --txHash TX_HASH

OPTIONS:
   --txHash value, --iTxID value  An Incognito transaction hash
   
```

#### pdestatus_withdrawlpfee
Check the status of a pDEX LP fee withdrawal transaction.
```shell
$ incognito-cli pdestatus help withdrawlpfee
NAME:
   incognito-cli pdestatus withdrawlpfee - Check the status of a pDEX LP fee withdrawal transaction.

USAGE:
   pdestatus withdrawlpfee --txHash TX_HASH

OPTIONS:
   --txHash value, --iTxID value  An Incognito transaction hash
   
```

#### pdestatus_withdraworder
Check the status of a pDEX order-withdrawal transaction.
```shell
$ incognito-cli pdestatus help withdraworder
NAME:
   incognito-cli pdestatus withdraworder - Check the status of a pDEX order-withdrawal transaction.

USAGE:
   pdestatus withdraworder --txHash TX_HASH

OPTIONS:
   --txHash value, --iTxID value  An Incognito transaction hash
   
```

#### pdestatus_withdrawstakereward
Check the status of a pDEX staking reward withdrawal transaction.
```shell
$ incognito-cli pdestatus help withdrawstakereward
NAME:
   incognito-cli pdestatus withdrawstakereward - Check the status of a pDEX staking reward withdrawal transaction.

USAGE:
   pdestatus withdrawstakereward --txHash TX_HASH

OPTIONS:
   --txHash value, --iTxID value  An Incognito transaction hash
   
```

## TRANSACTIONS
### checkreceiver
This command checks if an OTA key is a receiver of a transaction. If so, it will try to decrypt the received outputs and return the receiving info.
```shell
$ incognito-cli help checkreceiver
NAME:
   incognito-cli checkreceiver - Check if an OTA key is a receiver of a transaction.

USAGE:
   checkreceiver --txHash TX_HASH --otaKey OTA_KEY [--readonlyKey READONLY_KEY]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   TRANSACTIONS

DESCRIPTION:
   This command checks if an OTA key is a receiver of a transaction. If so, it will try to decrypt the received outputs and return the receiving info.

OPTIONS:
   --txHash value, --iTxID value    An Incognito transaction hash
   --otaKey value, --ota value      A base58-encoded ota key
   --readonlyKey value, --ro value  A base58-encoded read-only key
   
```

### convert
This command helps convert UTXOs v1 of a user to UTXO v2 w.r.t a tokenID. Please note that this process is time-consuming and requires a considerable amount of CPU.
```shell
$ incognito-cli help convert
NAME:
   incognito-cli convert - Convert UTXOs of an account w.r.t a tokenID.

USAGE:
   convert --privateKey PRIVATE_KEY [--tokenID TOKEN_ID] [--numThreads NUM_THREADS]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   TRANSACTIONS

DESCRIPTION:
   This command helps convert UTXOs v1 of a user to UTXO v2 w.r.t a tokenID. Please note that this process is time-consuming and requires a considerable amount of CPU.

OPTIONS:
   --privateKey value, -p value, --prvKey value  A base58-encoded Incognito private key
   --tokenID value, --id value, --ID value       The Incognito ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   --numThreads value                            Number of threads used in this action (default: 4)
   
```

### convertall
This command helps convert UTXOs v1 of a user to UTXO v2 for all assets. It will automatically check for all UTXOs v1 of all tokens and convert them. Please note that this process is time-consuming and requires a considerable amount of CPU.
```shell
$ incognito-cli help convertall
NAME:
   incognito-cli convertall - Convert UTXOs of an account for all assets.

USAGE:
   convertall --privateKey PRIVATE_KEY [--numThreads NUM_THREADS]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   TRANSACTIONS

DESCRIPTION:
   This command helps convert UTXOs v1 of a user to UTXO v2 for all assets. It will automatically check for all UTXOs v1 of all tokens and convert them. Please note that this process is time-consuming and requires a considerable amount of CPU.

OPTIONS:
   --privateKey value, -p value, --prvKey value  A base58-encoded Incognito private key
   --numThreads value                            Number of threads used in this action (default: 4)
   
```

### send
This command sends an amount of PRV or token from one wallet to another wallet. By default, it used 100 nano PRVs to pay the transaction fee.
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
   This command sends an amount of PRV or token from one wallet to another wallet. By default, it used 100 nano PRVs to pay the transaction fee.

OPTIONS:
   --privateKey value, -p value, --prvKey value  A base58-encoded Incognito private key
   --address value, --addr value                 A base58-encoded payment address
   --amount value, --amt value                   The Incognito amount of the action (default: 0)
   --tokenID value, --id value, --ID value       The Incognito ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   --fee value                                   The PRV amount for paying the transaction fee (default: 100)
   --version value, -v value                     Version of the transaction (1 or 2) (default: 2)
   
```

<!-- commandsstop -->
