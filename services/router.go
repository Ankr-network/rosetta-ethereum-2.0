package services

import (
	"net/http"

	"rosetta-ethereum-2.0/configuration"

	"github.com/coinbase/rosetta-sdk-go/asserter"
	"github.com/coinbase/rosetta-sdk-go/server"
)

// NewBlockchainRouter creates a Mux http.Handler from a collection
// of server controllers.
func NewBlockchainRouter(
	config *configuration.Configuration,
	client Client,
	asserter *asserter.Asserter,
) http.Handler {
	networkAPIService := NewNetworkAPIService(config, client)
	networkAPIController := server.NewNetworkAPIController(
		networkAPIService,
		asserter,
	)

	return server.NewRouter(
		networkAPIController,
	)
}
