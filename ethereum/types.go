package ethereum

import (
	"fmt"

	"github.com/coinbase/rosetta-sdk-go/types"
)

const (
	// NodeVersion is the version of prysm we are using.
	NodeVersion = "1.0.5"

	// Blockchain is Ethereum.
	Blockchain string = "Ethereum 2.0"

	// MainnetNetwork is the value of the network
	// in MainnetNetworkIdentifier.
	MainnetNetwork string = "Mainnet"

	// TestnetNetwork is the value of the network
	// in TestnetNetworkIdentifier.
	TestnetNetwork string = "Pyrmont"

	// HistoricalBalanceSupported is whether
	// historical balance is supported.
	HistoricalBalanceSupported = false

	// Symbol is the symbol value
	// used in Currency.
	Symbol = "ETH"

	// Decimals is the decimals value
	// used in Currency.
	Decimals = 18

	// SuccessStatus is the status of any
	// Ethereum operation considered successful.
	SuccessStatus = "SUCCESS"

	// FailureStatus is the status of any
	// Ethereum operation considered unsuccessful.
	FailureStatus = "FAILURE"

	// InputOpType is used to describe
	// INPUT.
	InputOpType = "INPUT"

	// OutputOpType is used to describe
	// OUTPUT.
	OutputOpType = "OUTPUT"

	// CoinbaseOpType is used to describe
	// Coinbase.
	CoinbaseOpType = "COINBASE"

	// MainnetPrysmArguments are the arguments to start a mainnet Prysm instance.
	MainnetPrysmArguments = `--config-file=/app/ethereum/prysm-config.yaml --datadir=/data`
)

var (
	// TestnetPrysmArguments are the arguments to start a ropsten Prysm instance.
	TestnetPrysmArguments = fmt.Sprintf("%s --pyrmont", MainnetPrysmArguments)

	// MainnetGenesisBlockIdentifier is the *types.BlockIdentifier
	// of the mainnet genesis block.
	MainnetGenesisBlockIdentifier = &types.BlockIdentifier{}

	// TestnetGenesisBlockIdentifier is the *types.BlockIdentifier
	// of the testnet genesis block.
	TestnetGenesisBlockIdentifier = &types.BlockIdentifier{}

	// Currency is the *types.Currency for all
	// Ethereum networks.
	Currency = &types.Currency{
		Symbol:   Symbol,
		Decimals: Decimals,
	}

	// OperationTypes are all supported operation.Types.
	OperationTypes = []string{
		InputOpType,
		OutputOpType,
		CoinbaseOpType,
	}

	// OperationStatuses are all supported operation statuses.
	OperationStatuses = []*types.OperationStatus{
		{
			Status:     SuccessStatus,
			Successful: true,
		},
		{
			Status:     FailureStatus,
			Successful: false,
		},
	}

	// CallMethods are all supported call methods.
	CallMethods = []string{}
)
