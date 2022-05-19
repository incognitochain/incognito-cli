package main

import "fmt"

const (
	UnexpectedError = iota
	UTXOVersionError
	NumThreadsError

	InvalidPrivateKeyError
	InvalidPaymentAddressError
	InvalidReadonlyKeyError
	InvalidOTAKeyError
	InvalidMiningKeyError
	InvalidTokenIDError

	GetBalanceError
	GetAllBalancesError
	GetAccountInfoError
	ConsolidateAccountError
	GetUnspentOutputCoinsError
	GetOutputCoinsError
	GetHistoryError
	SaveHistoryError
	GenerateMasterKeyError
	InvalidNumberShardsError
	InvalidShardError
	DeriveChildError
	ImportMnemonicError
	SubmitKeyError

	CreateStakingTransactionError
	CreateUnStakingTransactionError
	CreateWithdrawRewardTransactionError
	GetRewardAmountError
)

var errCodeMessages = map[int]struct {
	Code    int
	Message string
}{
	UnexpectedError:  {-1000, "Unexpected error"},
	UTXOVersionError: {-1001, "Expect version to be either 1 or 2"},
	NumThreadsError:  {-1002, "Expect numThreads to be greater than 0"},

	InvalidPrivateKeyError:     {-2000, "Invalid Incognito private key"},
	InvalidPaymentAddressError: {-2001, "Invalid Incognito payment address"},
	InvalidReadonlyKeyError:    {-2002, "Invalid Incognito readonly key"},
	InvalidOTAKeyError:         {-2003, "Invalid Incognito ota key"},
	InvalidMiningKeyError:      {-2004, "Invalid Incognito mining key"},
	InvalidTokenIDError:        {-2005, "Invalid Incognito tokenID"},

	GetBalanceError:            {-3000, "Error when retrieving balance"},
	GetAllBalancesError:        {-3001, "Error when retrieving all balances"},
	GetAccountInfoError:        {-3002, "Error when getting account info"},
	ConsolidateAccountError:    {-3003, "Consolidating error"},
	GetUnspentOutputCoinsError: {-3004, "Get UTXO error"},
	GetOutputCoinsError:        {-3005, "Get output coin error"},
	GetHistoryError:            {-3006, "Get account history error"},
	SaveHistoryError:           {-3007, "Save account history error"},
	GenerateMasterKeyError:     {-3008, "Generate master key error"},
	InvalidNumberShardsError:   {-3009, "Invalid number of shards"},
	InvalidShardError:          {-3010, "Invalid shard"},
	DeriveChildError:           {-3011, "Derive child error"},
	ImportMnemonicError:        {-3012, "Cannot import mnemonic"},
	SubmitKeyError:             {-3013, "Submit key error"},

	CreateStakingTransactionError:        {-4000, "Cannot create staking transaction"},
	CreateUnStakingTransactionError:      {-4001, "Cannot create un-staking transaction"},
	CreateWithdrawRewardTransactionError: {-4002, "Cannot create reward withdrawal transaction"},
	GetRewardAmountError:                 {-4003, "Cannot get reward amount"},
}

type appError struct {
	Code    int
	Message string
	Err     error
}

// Error satisfies the error interface and prints human-readable errors.
func (e appError) Error() error {
	if e.Err != nil {
		return fmt.Errorf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Errorf("[%d] %s", e.Code, e.Message)
}

func newAppError(key int, err ...error) error {
	res := appError{
		Code:    errCodeMessages[key].Code,
		Message: errCodeMessages[key].Message,
	}

	if len(err) > 0 {
		res.Err = err[0]
	}

	return res.Error()
}
