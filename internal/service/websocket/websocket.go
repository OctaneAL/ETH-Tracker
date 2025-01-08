package websocket

import (
	"context"
	"log"

	"github.com/OctaneAL/ETH-Tracker/internal/config"
	"github.com/OctaneAL/ETH-Tracker/internal/data"
	"github.com/OctaneAL/ETH-Tracker/internal/data/pg"
	"github.com/OctaneAL/ETH-Tracker/internal/erc20"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
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
	ws_endpoint := cfg.GetInfuraWsEndpoint()
	// https_endpoint := cfg.GetInfuraHttpsEndpoint()

	apiKey := cfg.GetInfuraAPIKey()

	tokenAddress := cfg.GetApiTokenAddress()

	websocketURL := ws_endpoint + apiKey
	// httpsURL := https_endpoint + apiKey

	client, err := ethclient.Dial(websocketURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	contractAddress := common.HexToAddress(tokenAddress)

	filterer, err := erc20.NewStorageFilterer(contractAddress, client)
	if err != nil {
		log.Fatalf("Failed to create filterer: %v", err)
	}

	database := pg.NewMasterQ(cfg.DB())

	transferChan := make(chan *erc20.StorageTransfer)
	subscription, err := filterer.WatchTransfer(&bind.WatchOpts{
		Context: context.Background(),
	}, transferChan, nil, nil)
	if err != nil {
		log.Fatalf("Failed to watch Transfer events: %v", err)
	}

	for {
		select {
		case err := <-subscription.Err():
			log.Printf("Subscription error: %v", err)
			return
		case event := <-transferChan:
			log.Printf("Transfer event received: From %s, To %s, Value %d",
				event.From.Hex(), event.To.Hex(), event.Value)

			transactionDetails, err := filterer.ParseTransfer(event.Raw)
			if err != nil {
				log.Fatalf("Failed to Parse Tranfer: %v", err)
			}

			transaction := data.InsertTransaction{
				BalanceNumeric:   transactionDetails.Value.String(),
				Sender:           transactionDetails.From.String(),
				Recipient:        transactionDetails.To.String(),
				TransactionHash:  event.Raw.TxHash.String(),
				TransactionIndex: string(event.Raw.TxIndex),
			}

			if transaction.BalanceNumeric != "0" {
				database.Trans().Insert(transaction)
			}

		}
	}
}
