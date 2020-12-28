package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountBalance(t *testing.T) {
	servicer := NewAccountAPIService()

	ctx := context.Background()

	bal, err := servicer.AccountBalance(ctx, nil)
	assert.Nil(t, bal)
	assert.Equal(t, ErrUnimplemented.Code, err.Code)
	assert.Equal(t, ErrUnimplemented.Message, err.Message)

	coins, err := servicer.AccountCoins(ctx, nil)
	assert.Nil(t, coins)
	assert.Equal(t, ErrUnimplemented.Code, err.Code)
	assert.Equal(t, ErrUnimplemented.Message, err.Message)

	mockClient.AssertExpectations(t)
}
