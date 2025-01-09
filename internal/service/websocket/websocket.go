package websocket

import (
	"context"
	"log"

	"github.com/OctaneAL/ETH-Tracker/internal/config"
	"github.com/OctaneAL/ETH-Tracker/internal/data/pg"
	"github.com/OctaneAL/ETH-Tracker/internal/erc20"
	"github.com/OctaneAL/ETH-Tracker/internal/models"
	"github.com/OctaneAL/ETH-Tracker/internal/service/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type Request struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type SubscriptionResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result"`
	Error   interface{} `json:"error"`
}

type EventData struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Subscription string `json:"subscription"`
		Result       string `json:"result"`
	} `json:"params"`
}

func SubscribeToLogs(cfg config.Config) {
	client_ws := cfg.GetWsClient()

	contractAddress := cfg.GetContractAddress()

	filterer, err := erc20.NewStorageFilterer(contractAddress, client_ws)
	if err != nil {
		log.Fatalf("Failed to create filterer: %v", err)
	}

	database := pg.NewMasterQ(cfg.DB())

	lastTransaction, err := database.Trans().GetLastRecord()
	if err != nil {
		log.Fatalf("Failed to get last transaction: %v", err)
	}
	blockFetchStart := uint64(lastTransaction.BlockNumber)

	transferChan := make(chan *erc20.StorageTransfer)
	subscription, err := filterer.WatchTransfer(&bind.WatchOpts{
		Start:   &blockFetchStart, // Start Watching from the last added block
		Context: context.Background(),
	}, transferChan, nil, nil)
	if err != nil {
		log.Fatalf("Failed to watch Transfer events: %v", err)
	}

	blockHash := models.BlockHash{
		BlockNumber: nil,
		Timestamp:   nil,
	}

	go func() {
		for {
			select {
			case err := <-subscription.Err():
				log.Printf("Subscription error: %v", err)
				return
			case event := <-transferChan:
				utils.ProcessTransferEvent(event, filterer, database, &blockHash, client_ws)
			}
		}
	}()

	select {}
}
