package services

import (
	"context"

	"github.com/coinbase/rosetta-sdk-go/types"
)

// AccountAPIService implements the server.AccountAPIServicer interface.
type AccountAPIService struct{}

// NewAccountAPIService returns a new *AccountAPIService.
func NewAccountAPIService() *AccountAPIService {
	return &AccountAPIService{}
}

// AccountBalance implements /account/balance.
func (s *AccountAPIService) AccountBalance(
	ctx context.Context,
	request *types.AccountBalanceRequest,
) (*types.AccountBalanceResponse, *types.Error) {
	return nil, wrapErr(ErrUnimplemented, nil)
}

// AccountCoins implements /account/coins.
func (s *AccountAPIService) AccountCoins(
	ctx context.Context,
	request *types.AccountCoinsRequest,
) (*types.AccountCoinsResponse, *types.Error) {
	return nil, wrapErr(ErrUnimplemented, nil)
}
