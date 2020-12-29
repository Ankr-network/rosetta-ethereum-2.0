package ethereum

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"rosetta-ethereum-2.0/timeutils"

	RosettaTypes "github.com/coinbase/rosetta-sdk-go/types"
	types "github.com/gogo/protobuf/types"
	pb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	grpc "google.golang.org/grpc"
)

const (
	secondsPerSlotCreation = uint64(12)
	grpcTimeout            = 120 * time.Second
)

// Client allows for querying a set of specific Ethereum 2.0 endpoints in an
// idempotent manner.
//
type Client struct {
	url               string
	nodeClient        pb.NodeClient
	beaconChainClient pb.BeaconChainClient
	conn              *grpc.ClientConn
}

func NewClient(ctx context.Context, url string) (*Client, error) {
	conn, err := grpc.DialContext(ctx, url, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect to remote wallet: %s", err)
	}

	bcc := pb.NewBeaconChainClient(conn)
	nc := pb.NewNodeClient(conn)

	return &Client{
		url:               url,
		nodeClient:        nc,
		beaconChainClient: bcc,
		conn:              conn,
	}, nil
}

// Close shuts down the RPC client connection.
func (ec *Client) Close() {
	ec.conn.Close()
}

func (ec *Client) Status(ctx context.Context) (
	*RosettaTypes.BlockIdentifier,
	*RosettaTypes.BlockIdentifier,
	int64,
	*RosettaTypes.SyncStatus,
	[]*RosettaTypes.Peer,
	error,
) {
	chainHead, err := ec.chainHead(ctx)
	if err != nil {
		return nil, nil, -1, nil, nil, err
	}

	genesis, err := ec.genesis(ctx)
	if err != nil {
		return nil, nil, -1, nil, nil, err
	}

	genesisTime := genesis.GetGenesisTime()
	highestBlock := getHighestBlock(genesisTime.GetSeconds())

	var syncStatus *RosettaTypes.SyncStatus
	currentIndex := int64(chainHead.GetHeadSlot())
	targetIndex := int64(highestBlock)

	stage := "synced"
	synced := true
	if currentIndex == 0 {
		stage = "deposit processing"
		synced = false
	} else if currentIndex < targetIndex && currentIndex > 0 {
		stage = "syncing"
		synced = false
	}

	syncStatus = &RosettaTypes.SyncStatus{
		CurrentIndex: &currentIndex,
		TargetIndex:  &targetIndex,
		Stage:        &stage,
		Synced:       &synced,
	}

	peers, err := ec.peers(ctx)
	if err != nil {
		return nil, nil, -1, nil, nil, err
	}

	return &RosettaTypes.BlockIdentifier{
			Hash:  hex.EncodeToString(chainHead.GetHeadBlockRoot()),
			Index: int64(chainHead.GetHeadSlot()),
		},
		&RosettaTypes.BlockIdentifier{
			Hash:  hex.EncodeToString(genesis.GetGenesisValidatorsRoot()),
			Index: 1,
		},
		timeutils.Now().Unix() * 1000,
		syncStatus,
		peers,
		nil
}

func (ec *Client) chainHead(ctx context.Context) (*pb.ChainHead, error) {
	res, err := ec.beaconChainClient.GetChainHead(ctx, &types.Empty{})
	if err != nil {
		log.Fatalf("could not get chain head: %s", err)
	}
	return res, nil
}

// Peers retrieves all peers of the node.
func (ec *Client) peers(ctx context.Context) ([]*RosettaTypes.Peer, error) {
	res, err := ec.nodeClient.ListPeers(ctx, &types.Empty{})
	if err != nil {
		log.Fatalf("could not list peers: %s", err)
	}
	info := res.GetPeers()

	peers := make([]*RosettaTypes.Peer, len(info))
	for i, peerInfo := range info {
		peers[i] = &RosettaTypes.Peer{
			PeerID: peerInfo.PeerId,
			Metadata: map[string]interface{}{
				"address":          peerInfo.Address,
				"direction":        peerInfo.Direction,
				"connection_state": peerInfo.ConnectionState,
				"enr":              peerInfo.Enr,
			},
		}
	}

	return peers, nil
}

// Genesis retrieves details of chain's genesis.
func (ec *Client) genesis(ctx context.Context) (*pb.Genesis, error) {
	genesis, err := ec.nodeClient.GetGenesis(ctx, &types.Empty{})
	if err != nil {
		log.Fatalf("could not retrieve genesis: %s", err)
	}

	return genesis, nil
}

func (ec *Client) Block(
	ctx context.Context,
	blockIdentifier *RosettaTypes.PartialBlockIdentifier,
) (*RosettaTypes.Block, error) {
	fmt.Printf("blockIdentifier: %s", blockIdentifier)
	if blockIdentifier != nil {
		if blockIdentifier.Hash != nil {
			return ec.blockByHash(ctx, *blockIdentifier.Hash)
		}

		if blockIdentifier.Index != nil {
			return ec.blockByIndex(ctx, *blockIdentifier.Index)
		}
	}

	return nil, errors.New("Query must be hash or index")
}

func (ec *Client) blockByIndex(ctx context.Context, block int64) (*RosettaTypes.Block, error) {
	var slot *pb.ListBlocksRequest_Slot
	fmt.Printf("Block number: %s", slot.Slot)
	slot.Slot = uint64(block)
	var in *pb.ListBlocksRequest
	in.QueryFilter = slot

	res, err := ec.beaconChainClient.ListBlocks(ctx, in)
	if err != nil {
		log.Fatalf("could not get block by slot index: %s", err)
	}
	fmt.Printf("ListBlocks: %s", res)
	rosettaBlock, err := ec.parseBeaconBlock(ctx, res)
	if err != nil {
		return nil, err
	}
	fmt.Printf("rosettaBlock: %s", rosettaBlock)
	return rosettaBlock, nil
}

func (ec *Client) blockByHash(ctx context.Context, hash string) (*RosettaTypes.Block, error) {
	var root *pb.ListBlocksRequest_Root
	root.Root = []byte(hash)
	var in *pb.ListBlocksRequest
	in.QueryFilter = root

	res, err := ec.beaconChainClient.ListBlocks(ctx, in)
	if err != nil {
		log.Fatalf("could not get block by root hash: %s", err)
	}
	rosettaBlock, err := ec.parseBeaconBlock(ctx, res)
	if err != nil {
		return nil, err
	}
	return rosettaBlock, nil
}

func (ec *Client) parseBeaconBlock(ctx context.Context, block *pb.ListBlocksResponse) (*RosettaTypes.Block, error) {
	if len(block.BlockContainers) < 1 {
		return nil, fmt.Errorf("No blocks were found %g", len(block.BlockContainers)) // !TODO add error output
	}

	var parentBlockIdentifier *RosettaTypes.BlockIdentifier
	b := block.BlockContainers[0]
	if b.Block.Block.Slot != 0 {
		parentBlockIdentifier = &RosettaTypes.BlockIdentifier{
			Index: int64(b.Block.Block.Slot - 1),
			Hash:  string(b.Block.Block.ParentRoot),
		}
	}
	fmt.Printf("parentBlockIdentifier:  %s", parentBlockIdentifier)

	timestamp, err := ec.getBlockTimestamp(ctx, int64(b.Block.Block.Slot))
	if err != nil {
		return nil, err
	}
	fmt.Printf("timestamp:  %s", timestamp)
	return &RosettaTypes.Block{
		BlockIdentifier: &RosettaTypes.BlockIdentifier{
			Index: int64(b.Block.Block.Slot),
			Hash:  string(b.BlockRoot),
		},
		ParentBlockIdentifier: parentBlockIdentifier,
		//The timestamp in milliseconds because some blockchains produce block more often than once a second.
		Timestamp:    timestamp * 1000,
		Transactions: nil,
		Metadata: map[string]interface{}{
			"epoch": int64(b.Block.Block.Slot) / 32,
			// "attestations": b.Block.Block.Body,
		},
	}, nil
}

func getHighestBlock(genesisTimeSec int64) uint64 {
	now := timeutils.Now().Unix()
	genesis := int64(genesisTimeSec)
	if now < genesis {
		return 0
	}
	return uint64(now-genesis) / secondsPerSlotCreation
}

func (ec *Client) getBlockTimestamp(ctx context.Context, blockNumber int64) (int64, error) {
	genesis, err := ec.genesis(ctx)
	if err != nil {
		return 0, err
	}
	genesisTime := genesis.GetGenesisTime()

	return (12 * blockNumber) + int64(genesisTime.GetSeconds()), nil

}

func convertTime(time uint64) int64 {
	return int64(time) * 1000
}
