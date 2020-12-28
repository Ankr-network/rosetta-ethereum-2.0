package services

import (
	"context"
	"testing"

	"rosetta-ethereum-2.0/configuration"
	"rosetta-ethereum-2.0/ethereum"
	mocks "rosetta-ethereum-2.0/mocks/services"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/stretchr/testify/assert"
)

var (
	middlewareVersion     = "0.0.1"
	defaultNetworkOptions = &types.NetworkOptionsResponse{
		Version: &types.Version{
			RosettaVersion:    types.RosettaAPIVersion,
			NodeVersion:       "1.0.5",
			MiddlewareVersion: &middlewareVersion,
		},
		Allow: &types.Allow{
			OperationStatuses:       ethereum.OperationStatuses,
			OperationTypes:          ethereum.OperationTypes,
			Errors:                  Errors,
			HistoricalBalanceLookup: ethereum.HistoricalBalanceSupported,
			CallMethods:             ethereum.CallMethods,
		},
	}

	networkIdentifier = &types.NetworkIdentifier{
		Network:    ethereum.MainnetNetwork,
		Blockchain: ethereum.Blockchain,
	}
)

func TestNetworkEndpoints_Offline(t *testing.T) {
	cfg := &configuration.Configuration{
		Mode:    configuration.Offline,
		Network: networkIdentifier,
	}
	mockClient := &mocks.Client{}
	servicer := NewNetworkAPIService(cfg, mockClient)
	ctx := context.Background()

	networkList, err := servicer.NetworkList(ctx, nil)
	assert.Nil(t, err)
	assert.Equal(t, []*types.NetworkIdentifier{
		networkIdentifier,
	}, networkList.NetworkIdentifiers)

	networkStatus, err := servicer.NetworkStatus(ctx, nil)
	assert.Nil(t, networkStatus)
	assert.Equal(t, ErrUnavailableOffline.Code, err.Code)
	assert.Equal(t, ErrUnavailableOffline.Message, err.Message)

	networkOptions, err := servicer.NetworkOptions(ctx, nil)
	assert.Nil(t, err)
	assert.Equal(t, defaultNetworkOptions, networkOptions)

	mockClient.AssertExpectations(t)
}

func TestNetworkEndpoints_Online(t *testing.T) {
	cfg := &configuration.Configuration{
		Mode:                   configuration.Online,
		Network:                networkIdentifier,
		GenesisBlockIdentifier: ethereum.MainnetGenesisBlockIdentifier,
	}
	mockClient := &mocks.Client{}
	servicer := NewNetworkAPIService(cfg, mockClient)
	ctx := context.Background()

	networkList, err := servicer.NetworkList(ctx, nil)
	assert.Nil(t, err)
	assert.Equal(t, []*types.NetworkIdentifier{
		networkIdentifier,
	}, networkList.NetworkIdentifiers)

	currentBlock := &types.BlockIdentifier{
		Index: 10,
		Hash:  "block 10",
	}

	currentTime := int64(1000000000000)

	syncStatus := &types.SyncStatus{
		CurrentIndex: types.Int64(100),
	}

	peers := []*types.Peer{
		{
			PeerID: "77.93.223.9:8333",
		},
	}

	mockClient.On(
		"Status",
		ctx,
	).Return(
		currentBlock,
		currentTime,
		syncStatus,
		peers,
		nil,
	)
	networkStatus, err := servicer.NetworkStatus(ctx, nil)
	assert.Nil(t, err)
	assert.Equal(t, &types.NetworkStatusResponse{
		GenesisBlockIdentifier: nil,
		CurrentBlockIdentifier: currentBlock,
		CurrentBlockTimestamp:  currentTime,
		Peers:                  peers,
		SyncStatus:             syncStatus,
	}, networkStatus)

	networkOptions, err := servicer.NetworkOptions(ctx, nil)
	assert.Nil(t, err)
	assert.Equal(t, defaultNetworkOptions, networkOptions)

	mockClient.AssertExpectations(t)
}
