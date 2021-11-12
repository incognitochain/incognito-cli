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
    * [`evm_retryshield`](#evm_retryshield)
    * [`evm_retryunshield`](#evm_retryunshield)
    * [`evm_shield`](#evm_shield)
    * [`evm_unshield`](#evm_unshield)
  * [`portal`](#portal)
    * [`portal_shield`](#portal_shield)
    * [`portal_shieldaddress`](#portal_shieldaddress)
    * [`portal_shieldstatus`](#portal_shieldstatus)
    * [`portal_unshield`](#portal_unshield)
    * [`portal_unshieldstatus`](#portal_unshieldstatus)
* [`COMMITTEES`](#committees)
  * [`checkrewards`](#checkrewards)
  * [`stake`](#stake)
  * [`unstake`](#unstake)
  * [`withdrawreward`](#withdrawreward)
* [`DEX`](#dex)
  * [`pdeaction`](#pdeaction)
    * [`pdeaction_addorder`](#pdeaction_addorder)
    * [`pdeaction_contribute`](#pdeaction_contribute)
    * [`pdeaction_mintnft`](#pdeaction_mintnft)
    * [`pdeaction_stake`](#pdeaction_stake)
    * [`pdeaction_trade`](#pdeaction_trade)
    * [`pdeaction_unstake`](#pdeaction_unstake)
    * [`pdeaction_withdraw`](#pdeaction_withdraw)
    * [`pdeaction_withdrawlpfee`](#pdeaction_withdrawlpfee)
    * [`pdeaction_withdraworder`](#pdeaction_withdraworder)
    * [`pdeaction_withdrawstakereward`](#pdeaction_withdrawstakereward)
  * [`pdeinfo`](#pdeinfo)
    * [`pdeinfo_checkprice`](#pdeinfo_checkprice)
    * [`pdeinfo_findpath`](#pdeinfo_findpath)
    * [`pdeinfo_lpvalue`](#pdeinfo_lpvalue)
    * [`pdeinfo_mynft`](#pdeinfo_mynft)
    * [`pdeinfo_share`](#pdeinfo_share)
    * [`pdeinfo_stakereward`](#pdeinfo_stakereward)
  * [`pdestatus`](#pdestatus)
    * [`pdestatus_addorder`](#pdestatus_addorder)
    * [`pdestatus_contribute`](#pdestatus_contribute)
    * [`pdestatus_mintnft`](#pdestatus_mintnft)
    * [`pdestatus_stake`](#pdestatus_stake)
    * [`pdestatus_trade`](#pdestatus_trade)
    * [`pdestatus_unstake`](#pdestatus_unstake)
    * [`pdestatus_withdraw`](#pdestatus_withdraw)
    * [`pdestatus_withdrawlpfee`](#pdestatus_withdrawlpfee)
    * [`pdestatus_withdraworder`](#pdestatus_withdraworder)
    * [`pdestatus_withdrawstakereward`](#pdestatus_withdrawstakereward)
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
$ incognito-cli help evm, retryshield
NAME:
   incognito-cli evm - Perform an EVM action (e.g, shield, unshield, etc.).

USAGE:
   evm

CATEGORY:
   BRIDGE

DESCRIPTION:
   This command helps perform an EVM action (e.g, shield, unshield, etc.).
```

#### evm_retryunshield
This command tries to un-shield an asset from an already-been-burned Incognito transaction in case of prior failure.
```shell
$ incognito-cli help evm, retryunshield
NAME:
   incognito-cli evm - Perform an EVM action (e.g, shield, unshield, etc.).

USAGE:
   evm

CATEGORY:
   BRIDGE

DESCRIPTION:
   This command helps perform an EVM action (e.g, shield, unshield, etc.).
```

#### evm_shield
This function helps shield an EVM (ETH/BNB/ERC20/BEP20) token into the Incognito network. It will ask for users' EVM PRIVATE KEY to proceed. The shielding process consists of the following operations.
1. Deposit the EVM asset into the corresponding smart contract.
1.1. In case the asset is an ERC20/BEP20 token, an approval transaction is performed (if needed) the before the actual deposit. For this operation, a prompt will be displayed to ask for user's approval.
2. Get the deposited EVM transaction, parse the depositing proof and submit it to the Incognito network. This step requires an Incognito private key with a sufficient amount of PRV to create an issuing transaction.

Note that EVM shielding is a complicated process, users MUST understand how the process works before using this function. We RECOMMEND users test the function with test networks BEFORE performing it on the live networks.
DO NOT USE THIS FUNCTION UNLESS YOU UNDERSTAND THE SHIELDING PROCESS.
```shell
$ incognito-cli help evm, shield
NAME:
   incognito-cli evm - Perform an EVM action (e.g, shield, unshield, etc.).

USAGE:
   evm

CATEGORY:
   BRIDGE

DESCRIPTION:
   This command helps perform an EVM action (e.g, shield, unshield, etc.).
```

#### evm_unshield
This function helps withdraw an EVM (ETH/BNB/ERC20/BEP20) token out of the Incognito network.The un-shielding process consists the following operations.
1. Users burn the token inside the Incognito chain.
2. After the burning is success, wait for 1-2 Incognito blocks and retrieve the corresponding burn proof from the Incognito chain.
3. After successfully retrieving the burn proof, users submit the burn proof to the smart contract to get back the corresponding public token. This step will ask for users' EVM PRIVATE KEY to proceed. Note that ONLY UNTIL this step, it is feasible to estimate the actual un-shielding fee (mainly is the fee interacting with the smart contract).

Please be aware that EVM un-shielding is a complicated process; and once burned, there is NO WAY to recover the asset inside the Incognito network. Therefore, use this function IF ADN ONLY IF you understand the way un-shielding works. Otherwise, use the un-shielding function from the Incognito app. We RECOMMEND users test the function with test networks BEFORE performing it on the live networks.
DO NOT USE THIS FUNCTION UNLESS YOU UNDERSTAND THE UN-SHIELDING PROCESS.
```shell
$ incognito-cli help evm, unshield
NAME:
   incognito-cli evm - Perform an EVM action (e.g, shield, unshield, etc.).

USAGE:
   evm

CATEGORY:
   BRIDGE

DESCRIPTION:
   This command helps perform an EVM action (e.g, shield, unshield, etc.).
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
$ incognito-cli help portal, shield
NAME:
   incognito-cli portal - Perform a portal action (e.g, shield, unshield, etc.).

USAGE:
   portal

CATEGORY:
   BRIDGE

DESCRIPTION:
   This command helps perform a portal action (e.g, shield, unshield, etc.).
```

#### portal_shieldaddress
This command helps generate the portal shielding address for a payment address and a tokenID.
```shell
$ incognito-cli help portal, shieldaddress
NAME:
   incognito-cli portal - Perform a portal action (e.g, shield, unshield, etc.).

USAGE:
   portal

CATEGORY:
   BRIDGE

DESCRIPTION:
   This command helps perform a portal action (e.g, shield, unshield, etc.).
```

#### portal_shieldstatus
This command helps retrieve the status of a portal shielding request.
Status should be understood as: 0 - rejected; 1 - accepted.
If you encounter an error, it might be because the request hasn't reached the beacon chain yet. Please try again a few minutes later.
```shell
$ incognito-cli help portal, shieldstatus
NAME:
   incognito-cli portal - Perform a portal action (e.g, shield, unshield, etc.).

USAGE:
   portal

CATEGORY:
   BRIDGE

DESCRIPTION:
   This command helps perform a portal action (e.g, shield, unshield, etc.).
```

#### portal_unshield
This command helps withdraw portal tokens (BTC) out of the Incognito network.
```shell
$ incognito-cli help portal, unshield
NAME:
   incognito-cli portal - Perform a portal action (e.g, shield, unshield, etc.).

USAGE:
   portal

CATEGORY:
   BRIDGE

DESCRIPTION:
   This command helps perform a portal action (e.g, shield, unshield, etc.).
```

#### portal_unshieldstatus
This command helps retrieve the status of a portal un-shielding request.
Status should be understood as: 0 - waiting; 1 - processed but not completed; 2 - completed; 3 - rejected.
If you encounter an error saying "unexpected end of JSON input", it might be because the request hasn't reached the beacon chain yet. Please try again a few minutes later.
```shell
$ incognito-cli help portal, unshieldstatus
NAME:
   incognito-cli portal - Perform a portal action (e.g, shield, unshield, etc.).

USAGE:
   portal

CATEGORY:
   BRIDGE

DESCRIPTION:
   This command helps perform a portal action (e.g, shield, unshield, etc.).
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
$ incognito-cli help pdeaction, addorder
NAME:
   incognito-cli pdeaction - Perform a pDEX action.

USAGE:
   pdeaction

CATEGORY:
   DEX

DESCRIPTION:
   This command helps perform a pDEX action.
```

#### pdeaction_contribute
This command creates a pDEX liquidity-contributing transaction. See more about this transaction: https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/pdex/contribute.md
```shell
$ incognito-cli help pdeaction, contribute
NAME:
   incognito-cli pdeaction - Perform a pDEX action.

USAGE:
   pdeaction

CATEGORY:
   DEX

DESCRIPTION:
   This command helps perform a pDEX action.
```

#### pdeaction_mintnft
This command creates and broadcasts a transaction that mints a new (pDEX) NFT for the pDEX.
```shell
$ incognito-cli help pdeaction, mintnft
NAME:
   incognito-cli pdeaction - Perform a pDEX action.

USAGE:
   pdeaction

CATEGORY:
   DEX

DESCRIPTION:
   This command helps perform a pDEX action.
```

#### pdeaction_stake
This command creates a transaction staking a token to the pDEX.
```shell
$ incognito-cli help pdeaction, stake
NAME:
   incognito-cli pdeaction - Perform a pDEX action.

USAGE:
   pdeaction

CATEGORY:
   DEX

DESCRIPTION:
   This command helps perform a pDEX action.
```

#### pdeaction_trade
This command creates a trade transaction on the pDEX.
```shell
$ incognito-cli help pdeaction, trade
NAME:
   incognito-cli pdeaction - Perform a pDEX action.

USAGE:
   pdeaction

CATEGORY:
   DEX

DESCRIPTION:
   This command helps perform a pDEX action.
```

#### pdeaction_unstake
This command creates a transaction un-staking a token from the pDEX.
```shell
$ incognito-cli help pdeaction, unstake
NAME:
   incognito-cli pdeaction - Perform a pDEX action.

USAGE:
   pdeaction

CATEGORY:
   DEX

DESCRIPTION:
   This command helps perform a pDEX action.
```

#### pdeaction_withdraw
This command creates a transaction withdrawing an amount of `share` from the pDEX. See more about this transaction: https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/pdex/withdrawal.md
```shell
$ incognito-cli help pdeaction, withdraw
NAME:
   incognito-cli pdeaction - Perform a pDEX action.

USAGE:
   pdeaction

CATEGORY:
   DEX

DESCRIPTION:
   This command helps perform a pDEX action.
```

#### pdeaction_withdrawlpfee
This command creates a transaction withdrawing LP fees from the pDEX.
```shell
$ incognito-cli help pdeaction, withdrawlpfee
NAME:
   incognito-cli pdeaction - Perform a pDEX action.

USAGE:
   pdeaction

CATEGORY:
   DEX

DESCRIPTION:
   This command helps perform a pDEX action.
```

#### pdeaction_withdraworder
This command creates a transaction withdrawing an order to the pDEX.
```shell
$ incognito-cli help pdeaction, withdraworder
NAME:
   incognito-cli pdeaction - Perform a pDEX action.

USAGE:
   pdeaction

CATEGORY:
   DEX

DESCRIPTION:
   This command helps perform a pDEX action.
```

#### pdeaction_withdrawstakereward
This command creates a transaction withdrawing staking rewards from the pDEX.
```shell
$ incognito-cli help pdeaction, withdrawstakereward
NAME:
   incognito-cli pdeaction - Perform a pDEX action.

USAGE:
   pdeaction

CATEGORY:
   DEX

DESCRIPTION:
   This command helps perform a pDEX action.
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
$ incognito-cli help pdeinfo, checkprice
NAME:
   incognito-cli pdeinfo - Retrieve pDEX information.

USAGE:
   pdeinfo

CATEGORY:
   DEX

DESCRIPTION:
   This command helps retrieve some information of the pDEX.
```

#### pdeinfo_findpath
This command helps find a good trading path for a trade.
```shell
$ incognito-cli help pdeinfo, findpath
NAME:
   incognito-cli pdeinfo - Retrieve pDEX information.

USAGE:
   pdeinfo

CATEGORY:
   DEX

DESCRIPTION:
   This command helps retrieve some information of the pDEX.
```

#### pdeinfo_lpvalue
This command retrieves the information about the value of an LP in a given pool.
```shell
$ incognito-cli help pdeinfo, lpvalue
NAME:
   incognito-cli pdeinfo - Retrieve pDEX information.

USAGE:
   pdeinfo

CATEGORY:
   DEX

DESCRIPTION:
   This command helps retrieve some information of the pDEX.
```

#### pdeinfo_mynft
This command returns the list of NFTs for a given private key.
```shell
$ incognito-cli help pdeinfo, mynft
NAME:
   incognito-cli pdeinfo - Retrieve pDEX information.

USAGE:
   pdeinfo

CATEGORY:
   DEX

DESCRIPTION:
   This command helps retrieve some information of the pDEX.
```

#### pdeinfo_share
This command returns the share amount of an nftID within a pDEX poolID.
```shell
$ incognito-cli help pdeinfo, share
NAME:
   incognito-cli pdeinfo - Retrieve pDEX information.

USAGE:
   pdeinfo

CATEGORY:
   DEX

DESCRIPTION:
   This command helps retrieve some information of the pDEX.
```

#### pdeinfo_stakereward
This command returns the estimated pDEX staking rewards of an nftID within a pDEX staking pool.
```shell
$ incognito-cli help pdeinfo, stakereward
NAME:
   incognito-cli pdeinfo - Retrieve pDEX information.

USAGE:
   pdeinfo

CATEGORY:
   DEX

DESCRIPTION:
   This command helps retrieve some information of the pDEX.
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
$ incognito-cli help pdestatus, addorder
NAME:
   incognito-cli pdestatus - Retrieve the status of a pDEX action.

USAGE:
   pdestatus

CATEGORY:
   DEX

DESCRIPTION:
   This command helps retrieve the status of a pDEX action given its hash. If an error is thrown, it is mainly because the transaction has not yet reached the beacon chain or the txHash is invalid.
```

#### pdestatus_contribute
Check the status of a pDEX liquidity contribution.
```shell
$ incognito-cli help pdestatus, contribute
NAME:
   incognito-cli pdestatus - Retrieve the status of a pDEX action.

USAGE:
   pdestatus

CATEGORY:
   DEX

DESCRIPTION:
   This command helps retrieve the status of a pDEX action given its hash. If an error is thrown, it is mainly because the transaction has not yet reached the beacon chain or the txHash is invalid.
```

#### pdestatus_mintnft
Check the status of a (pDEX) NFT minting transaction.
```shell
$ incognito-cli help pdestatus, mintnft
NAME:
   incognito-cli pdestatus - Retrieve the status of a pDEX action.

USAGE:
   pdestatus

CATEGORY:
   DEX

DESCRIPTION:
   This command helps retrieve the status of a pDEX action given its hash. If an error is thrown, it is mainly because the transaction has not yet reached the beacon chain or the txHash is invalid.
```

#### pdestatus_stake
Check the status of a pDEX staking transaction.
```shell
$ incognito-cli help pdestatus, stake
NAME:
   incognito-cli pdestatus - Retrieve the status of a pDEX action.

USAGE:
   pdestatus

CATEGORY:
   DEX

DESCRIPTION:
   This command helps retrieve the status of a pDEX action given its hash. If an error is thrown, it is mainly because the transaction has not yet reached the beacon chain or the txHash is invalid.
```

#### pdestatus_trade
Check the status of a pDEX trade.
```shell
$ incognito-cli help pdestatus, trade
NAME:
   incognito-cli pdestatus - Retrieve the status of a pDEX action.

USAGE:
   pdestatus

CATEGORY:
   DEX

DESCRIPTION:
   This command helps retrieve the status of a pDEX action given its hash. If an error is thrown, it is mainly because the transaction has not yet reached the beacon chain or the txHash is invalid.
```

#### pdestatus_unstake
Check the status of a pDEX un-staking transaction.
```shell
$ incognito-cli help pdestatus, unstake
NAME:
   incognito-cli pdestatus - Retrieve the status of a pDEX action.

USAGE:
   pdestatus

CATEGORY:
   DEX

DESCRIPTION:
   This command helps retrieve the status of a pDEX action given its hash. If an error is thrown, it is mainly because the transaction has not yet reached the beacon chain or the txHash is invalid.
```

#### pdestatus_withdraw
Check the status of a pDEX liquidity withdrawal.
```shell
$ incognito-cli help pdestatus, withdraw
NAME:
   incognito-cli pdestatus - Retrieve the status of a pDEX action.

USAGE:
   pdestatus

CATEGORY:
   DEX

DESCRIPTION:
   This command helps retrieve the status of a pDEX action given its hash. If an error is thrown, it is mainly because the transaction has not yet reached the beacon chain or the txHash is invalid.
```

#### pdestatus_withdrawlpfee
Check the status of a pDEX LP fee withdrawal transaction.
```shell
$ incognito-cli help pdestatus, withdrawlpfee
NAME:
   incognito-cli pdestatus - Retrieve the status of a pDEX action.

USAGE:
   pdestatus

CATEGORY:
   DEX

DESCRIPTION:
   This command helps retrieve the status of a pDEX action given its hash. If an error is thrown, it is mainly because the transaction has not yet reached the beacon chain or the txHash is invalid.
```

#### pdestatus_withdraworder
Check the status of a pDEX order-withdrawal transaction.
```shell
$ incognito-cli help pdestatus, withdraworder
NAME:
   incognito-cli pdestatus - Retrieve the status of a pDEX action.

USAGE:
   pdestatus

CATEGORY:
   DEX

DESCRIPTION:
   This command helps retrieve the status of a pDEX action given its hash. If an error is thrown, it is mainly because the transaction has not yet reached the beacon chain or the txHash is invalid.
```

#### pdestatus_withdrawstakereward
Check the status of a pDEX staking reward withdrawal transaction.
```shell
$ incognito-cli help pdestatus, withdrawstakereward
NAME:
   incognito-cli pdestatus - Retrieve the status of a pDEX action.

USAGE:
   pdestatus

CATEGORY:
   DEX

DESCRIPTION:
   This command helps retrieve the status of a pDEX action given its hash. If an error is thrown, it is mainly because the transaction has not yet reached the beacon chain or the txHash is invalid.
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
