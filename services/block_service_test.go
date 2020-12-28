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

func TestBlockService_Offline(t *testing.T) {
	cfg := &configuration.Configuration{
		Mode: configuration.Offline,
	}
	mockClient := &mocks.Client{}
	servicer := NewBlockAPIService(cfg, mockClient)
	ctx := context.Background()

	block, err := servicer.Block(ctx, &types.BlockRequest{})
	assert.Nil(t, block)
	assert.Equal(t, ErrUnavailableOffline.Code, err.Code)
	assert.Equal(t, ErrUnavailableOffline.Message, err.Message)

	blockTransaction, err := servicer.BlockTransaction(ctx, &types.BlockTransactionRequest{})
	assert.Nil(t, blockTransaction)
	assert.Equal(t, ErrUnimplemented.Code, err.Code)
	assert.Equal(t, ErrUnimplemented.Message, err.Message)

	mockClient.AssertExpectations(t)
}

func TestBlockService_Online(t *testing.T) {
	cfg := &configuration.Configuration{
		Mode: configuration.Online,
	}
	mockClient := &mocks.Client{}
	servicer := NewBlockAPIService(cfg, mockClient)
	ctx := context.Background()

	block := &types.Block{
		BlockIdentifier: &types.BlockIdentifier{
			Index: 100,
			Hash:  "block 100",
		},
	}

	blockResponse := &types.BlockResponse{
		Block: block,
	}

	t.Run("nil identifier", func(t *testing.T) {
		mockClient.On(
			"Block",
			ctx,
			(*types.PartialBlockIdentifier)(nil),
		).Return(
			block,
			nil,
		).Once()
		b, err := servicer.Block(ctx, &types.BlockRequest{})
		assert.Nil(t, err)
		assert.Equal(t, blockResponse, b)
	})

	t.Run("populated identifier", func(t *testing.T) {
		pbIdentifier := types.ConstructPartialBlockIdentifier(block.BlockIdentifier)
		mockClient.On("Block", ctx, pbIdentifier).Return(block, nil).Once()
		b, err := servicer.Block(ctx, &types.BlockRequest{
			BlockIdentifier: pbIdentifier,
		})
		assert.Nil(t, err)
		assert.Equal(t, blockResponse, b)
	})

	t.Run("orphaned block", func(t *testing.T) {
		pbIdentifier := types.ConstructPartialBlockIdentifier(block.BlockIdentifier)
		mockClient.On("Block", ctx, pbIdentifier).Return(nil, ethereum.ErrBlockOrphaned).Once()
		b, err := servicer.Block(ctx, &types.BlockRequest{
			BlockIdentifier: pbIdentifier,
		})

		assert.Nil(t, b)
		assert.Equal(t, ErrBlockOrphaned.Code, err.Code)
		assert.Equal(t, ErrBlockOrphaned.Message, err.Message)
		assert.Equal(t, ErrBlockOrphaned.Retriable, err.Retriable)
	})

	mockClient.AssertExpectations(t)
}
