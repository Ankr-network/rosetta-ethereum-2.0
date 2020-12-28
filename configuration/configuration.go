package configuration

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"rosetta-ethereum-2.0/ethereum"

	"github.com/coinbase/rosetta-sdk-go/types"
)

// Mode is the setting that determines if
// the implementation is "online" or "offline".
type Mode string

const (
	// Online is when the implementation is permitted
	// to make outbound connections.
	Online Mode = "ONLINE"

	// Offline is when the implementation is not permitted
	// to make outbound connections.
	Offline Mode = "OFFLINE"

	// Mainnet is the Ethereum 2.0 Mainnet.
	Mainnet string = "MAINNET"

	// Testnet is Ethereum 2.0 Testnet.
	Testnet string = "TESTNET"

	// DataDirectory is the default location for all
	// persistent data.
	DataDirectory = "/data"

	// ModeEnv is the environment variable read
	// to determine mode.
	ModeEnv = "MODE"

	// NetworkEnv is the environment variable
	// read to determine network.
	NetworkEnv = "NETWORK"

	// PortEnv is the environment variable
	// read to determine the port for the Rosetta
	// implementation.
	PortEnv = "PORT"

	// BeaconRPCEnv is an optional environment variable
	// used to connect rosetta-ethereum to an already
	// running beacon node.
	BeaconRPCEnv = "BEACON_RPC"

	// HTTPWeb3ProviderEnv is the environment variable
	// used to connect beacon-node to an already synced
	// ethereum node
	HTTPWeb3ProviderEnv = "WEB3PROVIDER"

	// DefaultHTTPWeb3Provider is the default URL
	// of already synced ethereum node used for connect
	// beacon-node to connect to
	DefaultHTTPWeb3Provider = "http://localhost:8545"

	// DefaultRPCURL is the default URL for
	// a running beacon node. This is used
	// when BeaconRPCEnv is not populated.
	DefaultRPCURL = "http://localhost:4000"

	// MiddlewareVersion is the version of rosetta-ethereum.
	MiddlewareVersion = "0.0.1"
)

// Configuration determines how
type Configuration struct {
	Mode                   Mode
	Network                *types.NetworkIdentifier
	GenesisBlockIdentifier *types.BlockIdentifier
	BeaconURL              string
	RemoteBeacon           bool
	Port                   int
	PrysmArguments         string
}

// LoadConfiguration attempts to create a new Configuration
// using the ENVs in the environment.
func LoadConfiguration() (*Configuration, error) {
	config := &Configuration{}

	modeValue := Mode(os.Getenv(ModeEnv))
	switch modeValue {
	case Online:
		config.Mode = Online
	case Offline:
		config.Mode = Offline
	case "":
		return nil, errors.New("MODE must be populated")
	default:
		return nil, fmt.Errorf("%s is not a valid mode", modeValue)
	}

	httpWeb3Provider := DefaultHTTPWeb3Provider
	envWeb3Provider := os.Getenv(HTTPWeb3ProviderEnv)
	if len(envWeb3Provider) > 0 {
		httpWeb3Provider = envWeb3Provider
	}

	networkValue := os.Getenv(NetworkEnv)
	switch networkValue {
	case Mainnet:
		config.Network = &types.NetworkIdentifier{
			Blockchain: ethereum.Blockchain,
			Network:    ethereum.MainnetNetwork,
		}
		config.PrysmArguments = ethereum.MainnetPrysmArguments + "--http-web3provider=" + httpWeb3Provider
	case Testnet:
		config.Network = &types.NetworkIdentifier{
			Blockchain: ethereum.Blockchain,
			Network:    ethereum.TestnetNetwork,
		}
		config.PrysmArguments = ethereum.TestnetPrysmArguments + "--http-web3provider=" + httpWeb3Provider
	case "":
		return nil, errors.New("NETWORK must be populated")
	default:
		return nil, fmt.Errorf("%s is not a valid network", networkValue)
	}

	config.BeaconURL = DefaultRPCURL
	envBeaconRPC := os.Getenv(BeaconRPCEnv)
	if len(envBeaconRPC) > 0 {
		config.RemoteBeacon = true
		config.BeaconURL = envBeaconRPC
	}

	portValue := os.Getenv(PortEnv)
	if len(portValue) == 0 {
		return nil, errors.New("PORT must be populated")
	}

	port, err := strconv.Atoi(portValue)
	if err != nil || len(portValue) == 0 || port <= 0 {
		return nil, fmt.Errorf("%w: unable to parse port %s", err, portValue)
	}
	config.Port = port

	return config, nil
}
