package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstructionService(t *testing.T) {
	servicer := NewConstructionAPIService()
	ctx := context.Background()

	deriveResponse, err := servicer.ConstructionDerive(ctx, nil)
	assert.Nil(t, deriveResponse)
	assert.Equal(t, ErrUnimplemented.Code, err.Code)
	assert.Equal(t, ErrUnimplemented.Message, err.Message)

	preprocessResponse, err := servicer.ConstructionPreprocess(ctx, nil)
	assert.Nil(t, preprocessResponse)
	assert.Equal(t, ErrUnimplemented.Code, err.Code)
	assert.Equal(t, ErrUnimplemented.Message, err.Message)

	metadataResponse, err := servicer.ConstructionMetadata(ctx, nil)
	assert.Nil(t, metadataResponse)
	assert.Equal(t, ErrUnimplemented.Code, err.Code)
	assert.Equal(t, ErrUnimplemented.Message, err.Message)

	payloadsResponse, err := servicer.ConstructionPayloads(ctx, nil)
	assert.Nil(t, payloadsResponse)
	assert.Equal(t, ErrUnimplemented.Code, err.Code)
	assert.Equal(t, ErrUnimplemented.Message, err.Message)

	combineResponse, err := servicer.ConstructionCombine(ctx, nil)
	assert.Nil(t, combineResponse)
	assert.Equal(t, ErrUnimplemented.Code, err.Code)
	assert.Equal(t, ErrUnimplemented.Message, err.Message)

	hashResponse, err := servicer.ConstructionHash(ctx, nil)
	assert.Nil(t, hashResponse)
	assert.Equal(t, ErrUnimplemented.Code, err.Code)
	assert.Equal(t, ErrUnimplemented.Message, err.Message)

	parseResponse, err := servicer.ConstructionParse(ctx, nil)
	assert.Nil(t, parseResponse)
	assert.Equal(t, ErrUnimplemented.Code, err.Code)
	assert.Equal(t, ErrUnimplemented.Message, err.Message)

	submitResponse, err := servicer.ConstructionSubmit(ctx, nil)
	assert.Nil(t, submitResponse)
	assert.Equal(t, ErrUnimplemented.Code, err.Code)
	assert.Equal(t, ErrUnimplemented.Message, err.Message)
}
