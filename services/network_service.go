package services

import (
	"context"

	"rosetta-ethereum-2.0/configuration"
	"rosetta-ethereum-2.0/ethereum"

	"github.com/coinbase/rosetta-sdk-go/asserter"
	"github.com/coinbase/rosetta-sdk-go/types"
)

// NetworkAPIService implements the server.NetworkAPIServicer interface.
type NetworkAPIService struct {
	config *configuration.Configuration
	client Client
}

// NewNetworkAPIService creates a new instance of a NetworkAPIService.
func NewNetworkAPIService(
	cfg *configuration.Configuration,
	client Client,
) *NetworkAPIService {
	return &NetworkAPIService{
		config: cfg,
		client: client,
	}
}

// NetworkList implements the /network/list endpoint
func (s *NetworkAPIService) NetworkList(
	ctx context.Context,
	request *types.MetadataRequest,
) (*types.NetworkListResponse, *types.Error) {
	return &types.NetworkListResponse{
		NetworkIdentifiers: []*types.NetworkIdentifier{s.config.Network},
	}, nil
}

// NetworkOptions implements the /network/options endpoint.
func (s *NetworkAPIService) NetworkOptions(
	ctx context.Context,
	request *types.NetworkRequest,
) (*types.NetworkOptionsResponse, *types.Error) {
	return &types.NetworkOptionsResponse{
		Version: &types.Version{
			NodeVersion:       ethereum.NodeVersion,
			RosettaVersion:    types.RosettaAPIVersion,
			MiddlewareVersion: types.String(configuration.MiddlewareVersion),
		},
		Allow: &types.Allow{
			Errors:                  Errors,
			OperationTypes:          ethereum.OperationTypes,
			OperationStatuses:       ethereum.OperationStatuses,
			HistoricalBalanceLookup: ethereum.HistoricalBalanceSupported,
			CallMethods:             nil,
		},
	}, nil
}

// NetworkStatus implements the /network/status endpoint.
func (s *NetworkAPIService) NetworkStatus(
	ctx context.Context,
	request *types.NetworkRequest,
) (*types.NetworkStatusResponse, *types.Error) {
	if s.config.Mode != configuration.Online {
		return nil, ErrUnavailableOffline
	}

	currentBlock, currentTime, syncStatus, peers, err := s.client.Status(ctx)
	if err != nil {
		return nil, wrapErr(ErrBeacon, err)
	}

	if currentTime < asserter.MinUnixEpoch {
		return nil, ErrBeaconNotReady
	}

	return &types.NetworkStatusResponse{
		CurrentBlockIdentifier: currentBlock,
		CurrentBlockTimestamp:  currentTime,
		GenesisBlockIdentifier: s.config.GenesisBlockIdentifier,
		SyncStatus:             syncStatus,
		Peers:                  peers,
	}, nil
}
