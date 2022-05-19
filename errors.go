package main

import "fmt"

const (
	UnexpectedError = iota
	VersionError
	NumThreadsError
	InvalidAmountError
	InvalidIncognitoTxHashError
	UserInputError

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
	InsufficientBalanceError

	CreateStakingTransactionError
	CreateUnStakingTransactionError
	CreateWithdrawRewardTransactionError
	GetRewardAmountError

	CreateTransferTransactionError
	CreateConversionTransactionError
	GetReceivingInfoError

	CentralizedShieldError

	GetEVMNetworkError
	InvalidEVMTokenAddressError
	EVMTokenIDToIncognitoTokenIDError
	IncognitoTokenIDToEVMTokenIDError
	GetEVMTokenInfoError
	WrongEVMNetworkError
	GetEVMBalanceError
	NewEVMAccountError
	GetEVMBurnProofError
	CreateEVMShieldingTransactionError
	CreateEVMUnShieldingTransactionError
	EVMDepositError
	EVMWithdrawError
	GetEVMShieldingStatusError
	CreatePRVShieldingTransactionError
	CreatePRVUnShieldingTransactionError
	EVMBurnPRVError
	EVMMintPRVError

	GenerateShieldingAddressError
	BTCClientNotFoundError
	GetBTCConfirmationError
	NotEnoughBTCConfirmationError
	BuildBTCProofError
	CreatePortalShieldingTransactionError
	CreatePortalUnShieldingTransactionError
	GetPortalShieldingStatusError
	GetPortalUnShieldingStatusError
	InvalidExternalAddressError
)

var errCodeMessages = map[int]struct {
	Code    int
	Message string
}{
	UnexpectedError:             {-1000, "Unexpected error"},
	VersionError:                {-1001, "Expect version to be either 1 or 2"},
	NumThreadsError:             {-1002, "Expect numThreads to be greater than 0"},
	InvalidAmountError:          {-1003, "Invalid Incognito amount"},
	InvalidIncognitoTxHashError: {-1004, "Invalid Incognito txHash"},
	UserInputError:              {-1005, "User input error"},

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
	InsufficientBalanceError:   {-3014, "Insufficient Incognito balance error"},

	CreateStakingTransactionError:        {-4000, "Cannot create staking transaction"},
	CreateUnStakingTransactionError:      {-4001, "Cannot create un-staking transaction"},
	CreateWithdrawRewardTransactionError: {-4002, "Cannot create reward withdrawal transaction"},
	GetRewardAmountError:                 {-4003, "Cannot get reward amount"},

	CreateTransferTransactionError:   {-5000, "Cannot create transfer transaction"},
	CreateConversionTransactionError: {-5001, "Cannot create conversion transaction"},
	GetReceivingInfoError:            {-5002, "Cannot get receiving info"},

	CentralizedShieldError: {-6000, "Cannot create centralized shielding transaction"},

	GetEVMNetworkError:                   {-6100, "Cannot get EVM network"},
	InvalidEVMTokenAddressError:          {-6101, "Invalid EVM token address"},
	EVMTokenIDToIncognitoTokenIDError:    {-6102, "Cannot get Incognito tokenID from EVM tokenID"},
	IncognitoTokenIDToEVMTokenIDError:    {-6103, "Cannot get EVM tokenID from Incognito tokenID"},
	GetEVMTokenInfoError:                 {-6104, "Cannot get EVM token info"},
	WrongEVMNetworkError:                 {-6105, "Wrong EVM network"},
	GetEVMBalanceError:                   {-6116, "Cannot get EVM balance"},
	NewEVMAccountError:                   {-6117, "Cannot create new EVM account"},
	GetEVMBurnProofError:                 {-6118, "Cannot get EVM burn proof"},
	CreateEVMShieldingTransactionError:   {-6119, "Cannot create EVM shielding transaction"},
	CreateEVMUnShieldingTransactionError: {-6120, "Cannot create EVM un-shielding transaction"},
	EVMDepositError:                      {-6121, "Cannot deposit EVM token to the smart contract"},
	EVMWithdrawError:                     {-6122, "Cannot withdraw EVM token from the smart contract"},
	GetEVMShieldingStatusError:           {-6123, "Cannot retrieve EVM shielding transaction"},
	CreatePRVShieldingTransactionError:   {-6124, "Cannot create PRV shielding transaction"},
	CreatePRVUnShieldingTransactionError: {-6125, "Cannot create PRV un-shielding transaction"},
	EVMBurnPRVError:                      {-6126, "Cannot burn PRV on EVM network"},
	EVMMintPRVError:                      {-6127, "Cannot mint PRV on EVM network"},

	GenerateShieldingAddressError:           {-6200, "Cannot generate shielding address"},
	BTCClientNotFoundError:                  {-6201, "BTC client not found"},
	GetBTCConfirmationError:                 {-6202, "Cannot get BTC confirmation"},
	NotEnoughBTCConfirmationError:           {-6203, "Need at least 6 confirmations"},
	BuildBTCProofError:                      {-6204, "Cannot build BTC proof"},
	CreatePortalShieldingTransactionError:   {-6205, "Cannot create portal shielding transaction"},
	CreatePortalUnShieldingTransactionError: {-6206, "Cannot create portal un-shielding transaction"},
	GetPortalShieldingStatusError:           {-6207, "Cannot retrieve portal shielding status"},
	GetPortalUnShieldingStatusError:         {-6208, "Cannot retrieve portal un-shielding status"},
	InvalidExternalAddressError:             {-6209, "Invalid external address"},
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
