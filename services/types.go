package services

import (
	"context"

	"github.com/coinbase/rosetta-sdk-go/types"
)

// Client is used by the servicers to get block
// data and to submit transactions.
type Client interface {
	Status(context.Context) (
		*types.BlockIdentifier,
		int64,
		*types.SyncStatus,
		[]*types.Peer,
		error,
	)
}
