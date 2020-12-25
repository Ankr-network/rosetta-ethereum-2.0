package ethereum

import (
	"context"
	"encoding/hex"
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

var (
	beaconRPC = "167.71.156.69:4000"
)

// Client allows for querying a set of specific Ethereum endpoints in an
// idempotent manner. Client relies on the eth_*, debug_*, and admin_*
// methods and on the graphql endpoint.
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
	int64,
	*RosettaTypes.SyncStatus,
	[]*RosettaTypes.Peer,
	error,
) {
	chainHead, err := ec.chainHead(ctx)
	if err != nil {
		return nil, -1, nil, nil, err
	}

	genesis, err := ec.genesis(ctx)
	if err != nil {
		return nil, -1, nil, nil, err
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
		return nil, -1, nil, nil, err
	}

	return &RosettaTypes.BlockIdentifier{
			Hash:  hex.EncodeToString(chainHead.GetHeadBlockRoot()),
			Index: int64(chainHead.GetHeadSlot()),
		},
		timeutils.Now().Unix(),
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

func getHighestBlock(genesisTimeSec int64) uint64 {
	now := timeutils.Now().Unix()
	genesis := int64(genesisTimeSec)
	if now < genesis {
		return 0
	}
	return uint64(now-genesis) / secondsPerSlotCreation
}

func convertTime(time uint64) int64 {
	return int64(time) * 1000
}
