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
   COMMITTEES:
     checkrewards    Get all rewards of a payment address.
     stake           Create a staking transaction (https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/staking/stake.md).
     unstake         Create an un-staking transaction (https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/staking/unstake.md).
     withdrawreward  Withdraw the reward of a privateKey w.r.t to a tokenID.
   EVMBRIDGE:
     evmretryshield    Retry a shield from the given already-been-deposited-to-sc EVM transaction.
     evmretryunshield  Retry an un-shielding request from the given already-been-burned Incognito transaction.
     evmshield         Shield an EVM (ETH/BNB/ERC20/BEP20) token into the Incognito network.
     evmunshield       Withdraw an EVM (ETH/BNB/ERC20/BEP20) token from the Incognito network.
   PDEX:
     pdecheckprice   Check the price between two tokenIDs.
     pdecontribute   Create a pDEX contributing transaction.
     pdeshare        Retrieve the share amount of a pDEX pai.r
     pdetrade        Create a trade transaction.
     pdetradestatus  Get the status of a trade.
     pdewithdraw     Create a pDEX withdrawal transaction.
   PORTAL:
     portalshield          Shield a portal token (e.g, BTC) into the Incognito network.
     portalshieldaddress   Generate a portal shielding address.
     portalshieldstatus    Get the status of a portal shielding request.
     portalunshield        Withdraw portal tokens (BTC) from the Incognito network.
     portalunshieldstatus  Get the status of a portal un-shielding request.
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
* [`COMMITTEES`](#committees)
  * [`checkrewards`](#checkrewards)
  * [`stake`](#stake)
  * [`unstake`](#unstake)
  * [`withdrawreward`](#withdrawreward)
* [`EVMBRIDGE`](#evmbridge)
  * [`evmretryshield`](#evmretryshield)
  * [`evmretryunshield`](#evmretryunshield)
  * [`evmshield`](#evmshield)
  * [`evmunshield`](#evmunshield)
* [`PDEX`](#pdex)
  * [`pdecheckprice`](#pdecheckprice)
  * [`pdecontribute`](#pdecontribute)
  * [`pdeshare`](#pdeshare)
  * [`pdetrade`](#pdetrade)
  * [`pdetradestatus`](#pdetradestatus)
  * [`pdewithdraw`](#pdewithdraw)
* [`PORTAL`](#portal)
  * [`portalshield`](#portalshield)
  * [`portalshieldaddress`](#portalshieldaddress)
  * [`portalshieldstatus`](#portalshieldstatus)
  * [`portalunshield`](#portalunshield)
  * [`portalunshieldstatus`](#portalunshieldstatus)
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
   --privateKey value, -p value  A base58-encoded Incognito private key
   --tokenID value, --id value   The Incognito ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   
```

### consolidate
This function helps consolidate UTXOs of an account. It consolidates a version of UTXOs at a time, users need to specify which version they need to consolidate. Please note that this process is time-consuming and requires a considerable amount of CPU.
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
   This function helps consolidate UTXOs of an account. It consolidates a version of UTXOs at a time, users need to specify which version they need to consolidate. Please note that this process is time-consuming and requires a considerable amount of CPU.

OPTIONS:
   --privateKey value, -p value  A base58-encoded Incognito private key
   --tokenID value, --id value   The Incognito ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   --version value, -v value     Version of the transaction (1 or 2) (default: 2)
   --numThreads value            Number of threads used in this action (default: 4)
   
```

### generateaccount
This function helps generate a new mnemonic phrase and its Incognito accounts.
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
   This function helps generate a new mnemonic phrase and its Incognito accounts.

OPTIONS:
   --numShards value  The number of shard (default: 8)
   
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
   --privateKey value, -p value  A base58-encoded Incognito private key
   --tokenID value               ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   --numThreads value            Number of threads used in this action (default: 4)
   --enableLog                   Enable log for this action (default: false)
   --logFile value               Location of the log file (default: "os.Stdout")
   --csvFile value, --csv value  The csv file location to store the history
   
```

### importeaccount
This function helps generate Incognito accounts given a mnemonic.
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
   This function helps generate Incognito accounts given a mnemonic.

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

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   ACCOUNTS

OPTIONS:
   --privateKey value, -p value  a base58-encoded private key
   
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
   --address value, --addr value    A base58-encoded payment address
   --otaKey value, --ota value      A base58-encoded ota key
   --readonlyKey value, --ro value  A base58-encoded read-only key
   --tokenID value, --id value      The Incognito ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   
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
   --privateKey value, -p value  A base58-encoded Incognito private key
   --tokenID value, --id value   The Incognito ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   
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
   --address value, --addr value  A base58-encoded payment address
   
```

### stake
Create a staking transaction (https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/staking/stake.md).
```shell
$ incognito-cli help stake
NAME:
   incognito-cli stake - Create a staking transaction (https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/staking/stake.md).

USAGE:
   stake --privateKey PRIVATE_KEY [--candidateAddress CANDIDATE_ADDRESS] [--rewardAddress REWARD_ADDRESS] [--autoReStake AUTO_RE_STAKE]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   COMMITTEES

OPTIONS:
   --privateKey value, -p value               A base58-encoded Incognito private key
   --candidateAddress value, --canAddr value  The Incognito payment address of the committee candidate (default: the payment address of the privateKey)
   --rewardAddress value, --rwdAddr value     The Incognito payment address of the reward receiver (default: the payment address of the privateKey)
   --autoReStake value                        Whether or not to automatically re-stake (0 - false, <> 0 - true) (default: 1)
   
```

### unstake
Create an un-staking transaction (https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/staking/unstake.md).
```shell
$ incognito-cli help unstake
NAME:
   incognito-cli unstake - Create an un-staking transaction (https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/staking/unstake.md).

USAGE:
   unstake --privateKey PRIVATE_KEY [--candidateAddress CANDIDATE_ADDRESS]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   COMMITTEES

OPTIONS:
   --privateKey value, -p value               A base58-encoded Incognito private key
   --candidateAddress value, --canAddr value  The Incognito payment address of the committee candidate (default: the payment address of the privateKey)
   
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
   --privateKey value, -p value   A base58-encoded Incognito private key
   --address value, --addr value  the payment address of a candidate (default: the payment address of the privateKey)
   --tokenID value, --id value    The Incognito ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   --version value, -v value      Version of the transaction (1 or 2) (default: 2)
   
```

## EVMBRIDGE
### evmretryshield
This function re-shields an already-been-deposited-to-sc transaction in case of prior failure.
```shell
$ incognito-cli help evmretryshield
NAME:
   incognito-cli evmretryshield - Retry a shield from the given already-been-deposited-to-sc EVM transaction.

USAGE:
   evmretryshield --privateKey PRIVATE_KEY --externalTxHash EXTERNAL_TX_HASH [--evm EVM] [--externalTokenAddress EXTERNAL_TOKEN_ADDRESS]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   EVMBRIDGE

DESCRIPTION:
   This function re-shields an already-been-deposited-to-sc transaction in case of prior failure.

OPTIONS:
   --privateKey value, -p value           A base58-encoded Incognito private key
   --externalTxHash value, --eTxID value  The external transaction hash
   --evm value                            The EVM network (ETH or BSC) (default: "ETH")
   --externalTokenAddress value           ID of the token on ETH/BSC networks (default: "0x0000000000000000000000000000000000000000")
   
```

### evmretryunshield
This function tries to un-shield an asset from an already-been-burned Incognito transaction in case of prior failure.
```shell
$ incognito-cli help evmretryunshield
NAME:
   incognito-cli evmretryunshield - Retry an un-shielding request from the given already-been-burned Incognito transaction.

USAGE:
   evmretryunshield --txHash TX_HASH [--evm EVM]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   EVMBRIDGE

DESCRIPTION:
   This function tries to un-shield an asset from an already-been-burned Incognito transaction in case of prior failure.

OPTIONS:
   --txHash value, --iTxID value  An Incognito transaction hash
   --evm value                    The EVM network (ETH or BSC) (default: "ETH")
   
```

### evmshield
This function helps shield an EVM (ETH/BNB/ERC20/BEP20) token into the Incognito network. It will ask for users' EVM PRIVATE KEY to proceed. The shielding process consists of the following operations.
1. Deposit the EVM asset into the corresponding smart contract.
1.1. In case the asset is an ERC20/BEP20 token, an approval transaction is performed (if needed) the before the actual deposit. For this operation, a prompt will be displayed to ask for user's approval.
2. Get the deposited EVM transaction, parse the depositing proof and submit it to the Incognito network. This step requires an Incognito private key with a sufficient amount of PRV to create an issuing transaction.

Note that EVM shielding is a complicated process, users MUST understand how the process works before using this function. We RECOMMEND users test the function with test networks BEFORE performing it on the live networks.
DO NOT USE THIS FUNCTION UNLESS YOU UNDERSTAND THE SHIELDING PROCESS.
```shell
$ incognito-cli help evmshield
NAME:
   incognito-cli evmshield - Shield an EVM (ETH/BNB/ERC20/BEP20) token into the Incognito network.

USAGE:
   evmshield --privateKey PRIVATE_KEY --shieldAmount SHIELD_AMOUNT [--evm EVM] [--externalTokenAddress EXTERNAL_TOKEN_ADDRESS] [--address ADDRESS]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   EVMBRIDGE

DESCRIPTION:
   This function helps shield an EVM (ETH/BNB/ERC20/BEP20) token into the Incognito network. It will ask for users' EVM PRIVATE KEY to proceed. The shielding process consists of the following operations.
      1. Deposit the EVM asset into the corresponding smart contract.
        1.1. In case the asset is an ERC20/BEP20 token, an approval transaction is performed (if needed) the before the actual deposit. For this operation, a prompt will be displayed to ask for user's approval.
      2. Get the deposited EVM transaction, parse the depositing proof and submit it to the Incognito network. This step requires an Incognito private key with a sufficient amount of PRV to create an issuing transaction.
   
   Note that EVM shielding is a complicated process, users MUST understand how the process works before using this function. We RECOMMEND users test the function with test networks BEFORE performing it on the live networks.
   DO NOT USE THIS FUNCTION UNLESS YOU UNDERSTAND THE SHIELDING PROCESS.

OPTIONS:
   --privateKey value, -p value       A base58-encoded Incognito private key
   --shieldAmount value, --amt value  The shielding amount measured in token unit (e.g, 10, 1, 0.1, 0.01) (default: 0)
   --evm value                        The EVM network (ETH or BSC) (default: "ETH")
   --externalTokenAddress value       ID of the token on ETH/BSC networks (default: "0x0000000000000000000000000000000000000000")
   --address value, --addr value      The Incognito payment address to receive the shielding asset (default: the payment address of the privateKey)
   
```

### evmunshield
This function helps withdraw an EVM (ETH/BNB/ERC20/BEP20) token out of the Incognito network.The un-shielding process consists the following operations.
1. Users burn the token inside the Incognito chain.
2. After the burning is success, wait for 1-2 Incognito blocks and retrieve the corresponding burn proof from the Incognito chain.
3. After successfully retrieving the burn proof, users submit the burn proof to the smart contract to get back the corresponding public token. This step will ask for users' EVM PRIVATE KEY to proceed. Note that ONLY UNTIL this step, it is feasible to estimate the actual un-shielding fee (mainly is the fee interacting with the smart contract).

Please be aware that EVM un-shielding is a complicated process; and once burned, there is NO WAY to recover the asset inside the Incognito network. Therefore, use this function IF ADN ONLY IF you understand the way un-shielding works. Otherwise, use the un-shielding function from the Incognito app. We RECOMMEND users test the function with test networks BEFORE performing it on the live networks.
DO NOT USE THIS FUNCTION UNLESS YOU UNDERSTAND THE UN-SHIELDING PROCESS.
```shell
$ incognito-cli help evmunshield
NAME:
   incognito-cli evmunshield - Withdraw an EVM (ETH/BNB/ERC20/BEP20) token from the Incognito network.

USAGE:
   evmunshield --privateKey PRIVATE_KEY --tokenID TOKEN_ID --amount AMOUNT

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   EVMBRIDGE

DESCRIPTION:
   This function helps withdraw an EVM (ETH/BNB/ERC20/BEP20) token out of the Incognito network.The un-shielding process consists the following operations.
      1. Users burn the token inside the Incognito chain.
      2. After the burning is success, wait for 1-2 Incognito blocks and retrieve the corresponding burn proof from the Incognito chain.
      3. After successfully retrieving the burn proof, users submit the burn proof to the smart contract to get back the corresponding public token. This step will ask for users' EVM PRIVATE KEY to proceed. Note that ONLY UNTIL this step, it is feasible to estimate the actual un-shielding fee (mainly is the fee interacting with the smart contract).
   
   Please be aware that EVM un-shielding is a complicated process; and once burned, there is NO WAY to recover the asset inside the Incognito network. Therefore, use this function IF ADN ONLY IF you understand the way un-shielding works. Otherwise, use the un-shielding function from the Incognito app. We RECOMMEND users test the function with test networks BEFORE performing it on the live networks.
   DO NOT USE THIS FUNCTION UNLESS YOU UNDERSTAND THE UN-SHIELDING PROCESS.

OPTIONS:
   --privateKey value, -p value  A base58-encoded Incognito private key
   --tokenID value, --id value   The Incognito tokenID of the un-shielding asset
   --amount value, --amt value   The Incognito amount of the action (default: 0)
   
```

## PDEX
### pdecheckprice
This function checks the price of a pair of tokenIds. It must be supplied with the selling amount since the pDEX uses the AMM algorithm.
```shell
$ incognito-cli help pdecheckprice
NAME:
   incognito-cli pdecheckprice - Check the price between two tokenIDs.

USAGE:
   pdecheckprice --sellTokenID SELL_TOKEN_ID --buyTokenID BUY_TOKEN_ID --sellingAmount SELLING_AMOUNT

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   PDEX

DESCRIPTION:
   This function checks the price of a pair of tokenIds. It must be supplied with the selling amount since the pDEX uses the AMM algorithm.

OPTIONS:
   --sellTokenID value    ID of the token to sell
   --buyTokenID value     ID of the token to buy
   --sellingAmount value  The amount of sellTokenID wished to sell (default: 0)
   
```

### pdecontribute
This function creates a pDEX contributing transaction. See more about this transaction: https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/pdex/contribute.md
```shell
$ incognito-cli help pdecontribute
NAME:
   incognito-cli pdecontribute - Create a pDEX contributing transaction.

USAGE:
   pdecontribute --privateKey PRIVATE_KEY --pairId PAIR_ID [--tokenID TOKEN_ID] --amount AMOUNT [--version VERSION]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   PDEX

DESCRIPTION:
   This function creates a pDEX contributing transaction. See more about this transaction: https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/pdex/contribute.md

OPTIONS:
   --privateKey value, -p value  A base58-encoded Incognito private key
   --pairId value                The ID of the contributing pair (see https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/pdex/contribute.md)
   --tokenID value, --id value   The Incognito ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   --amount value, --amt value   The Incognito amount of the action (default: 0)
   --version value, -v value     Version of the transaction (1 or 2) (default: 2)
   
```

### pdeshare
This function returns the share amount of a user within a pDEX pair.
```shell
$ incognito-cli help pdeshare
NAME:
   incognito-cli pdeshare - Retrieve the share amount of a pDEX pai.r

USAGE:
   pdeshare --address ADDRESS --tokenID1 TOKEN_ID_1 [--tokenID2 TOKEN_ID_2]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   PDEX

DESCRIPTION:
   This function returns the share amount of a user within a pDEX pair.

OPTIONS:
   --address value, --addr value  A base58-encoded payment address
   --tokenID1 value               ID of the first token
   --tokenID2 value               ID of the second token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   
```

### pdetrade
This function creates a trade transaction on the pDEX.
```shell
$ incognito-cli help pdetrade
NAME:
   incognito-cli pdetrade - Create a trade transaction.

USAGE:
   pdetrade --privateKey PRIVATE_KEY --sellTokenID SELL_TOKEN_ID --buyTokenID BUY_TOKEN_ID --sellingAmount SELLING_AMOUNT [--minAcceptAmount MIN_ACCEPT_AMOUNT] [--tradingFee TRADING_FEE]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   PDEX

DESCRIPTION:
   This function creates a trade transaction on the pDEX.

OPTIONS:
   --privateKey value, -p value  A base58-encoded Incognito private key
   --sellTokenID value           ID of the token to sell
   --buyTokenID value            ID of the token to buy
   --sellingAmount value         The amount of sellTokenID wished to sell (default: 0)
   --minAcceptAmount value       The minimum acceptable amount of buyTokenID wished to receive (default: 0)
   --tradingFee value            The trading fee (measured in nano PRV) (default: 0)
   
```

### pdetradestatus
This function returns the status of a trade (1: successful, 2: failed). If a `not found` error occurs, it means that the trade has not been acknowledged by the beacon chain. Just wait and check again later.
```shell
$ incognito-cli help pdetradestatus
NAME:
   incognito-cli pdetradestatus - Get the status of a trade.

USAGE:
   pdetradestatus --txHash TX_HASH

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   PDEX

DESCRIPTION:
   This function returns the status of a trade (1: successful, 2: failed). If a `not found` error occurs, it means that the trade has not been acknowledged by the beacon chain. Just wait and check again later.

OPTIONS:
   --txHash value, --iTxID value  An Incognito transaction hash
   
```

### pdewithdraw
This function creates a transaction withdrawing an amount of `shared` from the pDEX. See more about this transaction: https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/pdex/withdrawal.md
```shell
$ incognito-cli help pdewithdraw
NAME:
   incognito-cli pdewithdraw - Create a pDEX withdrawal transaction.

USAGE:
   pdewithdraw --privateKey PRIVATE_KEY --amount AMOUNT --tokenID1 TOKEN_ID_1 [--tokenID2 TOKEN_ID_2] [--version VERSION]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   PDEX

DESCRIPTION:
   This function creates a transaction withdrawing an amount of `shared` from the pDEX. See more about this transaction: https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/pdex/withdrawal.md

OPTIONS:
   --privateKey value, -p value  A base58-encoded Incognito private key
   --amount value, --amt value   The Incognito amount of the action (default: 0)
   --tokenID1 value              ID of the first token
   --tokenID2 value              ID of the second token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   --version value, -v value     Version of the transaction (1 or 2) (default: 2)
   
```

## PORTAL
### portalshield
This function helps shield a portal token into the Incognito network after the fund has been transferred to the depositing address (generated by `portalshieldaddress`).
```shell
$ incognito-cli help portalshield
NAME:
   incognito-cli portalshield - Shield a portal token (e.g, BTC) into the Incognito network.

USAGE:
   portalshield --privateKey PRIVATE_KEY --externalTxHash EXTERNAL_TX_HASH [--tokenID TOKEN_ID] [--address ADDRESS]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   PORTAL

DESCRIPTION:
   This function helps shield a portal token into the Incognito network after the fund has been transferred to the depositing address (generated by `portalshieldaddress`).

OPTIONS:
   --privateKey value, -p value           A base58-encoded Incognito private key
   --externalTxHash value, --eTxID value  The external transaction hash
   --tokenID value, --id value            The Incognito tokenID of the shielding asset (default: "b832e5d3b1f01a4f0623f7fe91d6673461e1f5d37d91fe78c5c2e6183ff39696")
   --address value, --addr value          The Incognito payment address to receive the shielding asset (default: the payment address of the privateKey)
   
```

### portalshieldaddress
This function helps generate the portal shielding address for a payment address and a tokenID.
```shell
$ incognito-cli help portalshieldaddress
NAME:
   incognito-cli portalshieldaddress - Generate a portal shielding address.

USAGE:
   portalshieldaddress --address ADDRESS [--tokenID TOKEN_ID]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   PORTAL

DESCRIPTION:
   This function helps generate the portal shielding address for a payment address and a tokenID.

OPTIONS:
   --address value, --addr value  A base58-encoded payment address
   --tokenID value, --id value    The Incognito tokenID of the shielding asset (default: "b832e5d3b1f01a4f0623f7fe91d6673461e1f5d37d91fe78c5c2e6183ff39696")
   
```

### portalshieldstatus
This function helps retrieve the status of a portal shielding request.
Status should be understood as: 0 - rejected; 1 - accepted.
If you encounter an error, it might be because the request hasn't reached the beacon chain yet. Please try again a few minutes later.
```shell
$ incognito-cli help portalshieldstatus
NAME:
   incognito-cli portalshieldstatus - Get the status of a portal shielding request.

USAGE:
   portalshieldstatus --txHash TX_HASH

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   PORTAL

DESCRIPTION:
   This function helps retrieve the status of a portal shielding request.
   Status should be understood as: 0 - rejected; 1 - accepted.
   If you encounter an error, it might be because the request hasn't reached the beacon chain yet. Please try again a few minutes later.

OPTIONS:
   --txHash value, --iTxID value  An Incognito transaction hash
   
```

### portalunshield
This function helps withdraw portal tokens (BTC) out of the Incognito network.
```shell
$ incognito-cli help portalunshield
NAME:
   incognito-cli portalunshield - Withdraw portal tokens (BTC) from the Incognito network.

USAGE:
   portalunshield --privateKey PRIVATE_KEY --externalAddress EXTERNAL_ADDRESS --amount AMOUNT [--tokenID TOKEN_ID]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   PORTAL

DESCRIPTION:
   This function helps withdraw portal tokens (BTC) out of the Incognito network.

OPTIONS:
   --privateKey value, -p value            A base58-encoded Incognito private key
   --externalAddress value, --eAddr value  A valid remote address for the currently-processed tokenID. User MUST make sure this address is valid to avoid the loss of money.
   --amount value, --amt value             The Incognito amount of the action (default: 0)
   --tokenID value, --id value             The Incognito tokenID of the un-shielding asset (default: "b832e5d3b1f01a4f0623f7fe91d6673461e1f5d37d91fe78c5c2e6183ff39696")
   
```

### portalunshieldstatus
This function helps retrieve the status of a portal un-shielding request.
Status should be understood as: 0 - waiting; 1 - processed but not completed; 2 - completed; 3 - rejected.
If you encounter an error saying "unexpected end of JSON input", it might be because the request hasn't reached the beacon chain yet. Please try again a few minutes later.
```shell
$ incognito-cli help portalunshieldstatus
NAME:
   incognito-cli portalunshieldstatus - Get the status of a portal un-shielding request.

USAGE:
   portalunshieldstatus --txHash TX_HASH

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   PORTAL

DESCRIPTION:
   This function helps retrieve the status of a portal un-shielding request.
   Status should be understood as: 0 - waiting; 1 - processed but not completed; 2 - completed; 3 - rejected.
   If you encounter an error saying "unexpected end of JSON input", it might be because the request hasn't reached the beacon chain yet. Please try again a few minutes later.

OPTIONS:
   --txHash value, --iTxID value  An Incognito transaction hash
   
```

## TRANSACTIONS
### checkreceiver
This function checks if an OTA key is a receiver of a transaction. If so, it will try to decrypt the received outputs and return the receiving info.
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
   This function checks if an OTA key is a receiver of a transaction. If so, it will try to decrypt the received outputs and return the receiving info.

OPTIONS:
   --txHash value, --iTxID value    An Incognito transaction hash
   --otaKey value, --ota value      A base58-encoded ota key
   --readonlyKey value, --ro value  A base58-encoded read-only key
   
```

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
   --privateKey value, -p value  A base58-encoded Incognito private key
   --tokenID value, --id value   The Incognito ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   --numThreads value            Number of threads used in this action (default: 4)
   --enableLog                   Enable log for this action (default: false)
   --logFile value               Location of the log file (default: "os.Stdout")
   
```

### convertall
This function helps convert UTXOs v1 of a user to UTXO v2 for all assets. It will automatically check for all UTXOs v1 of all tokens and convert them. Please note that this process is time-consuming and requires a considerable amount of CPU.
```shell
$ incognito-cli help convertall
NAME:
   incognito-cli convertall - Convert UTXOs of an account for all assets.

USAGE:
   convertall --privateKey PRIVATE_KEY [--numThreads NUM_THREADS] [--logFile LOG_FILE]

   OPTIONAL flags are denoted by a [] bracket.

CATEGORY:
   TRANSACTIONS

DESCRIPTION:
   This function helps convert UTXOs v1 of a user to UTXO v2 for all assets. It will automatically check for all UTXOs v1 of all tokens and convert them. Please note that this process is time-consuming and requires a considerable amount of CPU.

OPTIONS:
   --privateKey value, -p value  A base58-encoded Incognito private key
   --numThreads value            Number of threads used in this action (default: 4)
   --logFile value               Location of the log file (default: "os.Stdout")
   
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
   --privateKey value, -p value   A base58-encoded Incognito private key
   --address value, --addr value  A base58-encoded payment address
   --amount value, --amt value    The Incognito amount of the action (default: 0)
   --tokenID value, --id value    The Incognito ID of the token (default: "0000000000000000000000000000000000000000000000000000000000000004")
   --fee value                    The PRV amount for paying the transaction fee (default: 100)
   --version value, -v value      Version of the transaction (1 or 2) (default: 2)
   
```

<!-- commandsstop -->
