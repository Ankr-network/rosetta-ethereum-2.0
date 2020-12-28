package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMempoolEndpoints(t *testing.T) {
	servicer := NewMempoolAPIService()
	ctx := context.Background()

	mem, err := servicer.Mempool(ctx, nil)
	assert.Nil(t, mem)
	assert.Equal(t, ErrUnimplemented.Code, err.Code)
	assert.Equal(t, ErrUnimplemented.Message, err.Message)

	memTransaction, err := servicer.MempoolTransaction(ctx, nil)
	assert.Nil(t, memTransaction)
	assert.Equal(t, ErrUnimplemented.Code, err.Code)
	assert.Equal(t, ErrUnimplemented.Message, err.Message)
}
