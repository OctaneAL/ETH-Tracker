package websocket

import (
	"context"
	"log"
	"time"

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
	// https_endpoint := cfg.GetInfuraHttpSsEndpoint()

	apiKey := cfg.GetInfuraAPIKey()

	tokenAddress := cfg.GetApiTokenAddress()

	websocketURL := ws_endpoint + apiKey
	// httpsURL := https_endpoint + apiKey

	client_ws, err := ethclient.Dial(websocketURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	contractAddress := common.HexToAddress(tokenAddress)

	filterer, err := erc20.NewStorageFilterer(contractAddress, client_ws)
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

	go func() {
		for {
			select {
			case err := <-subscription.Err():
				log.Printf("Subscription error: %v", err)
				return
			case event := <-transferChan:
				log.Printf("Transfer event received: From %s, To %s, Value %d, Block %d",
					event.From.Hex(), event.To.Hex(), event.Value, event.Raw.BlockNumber)

				transactionDetails, err := filterer.ParseTransfer(event.Raw)
				if err != nil {
					log.Fatalf("Failed to Parse Tranfer: %v", err)
				}

				// Takes too much time

				// blockNumber := big.NewInt(int64(event.Raw.BlockNumber))
				// block, err := client_ws.BlockByNumber(context.Background(), blockNumber)
				// if err != nil {
				// 	log.Printf("Failed to fetch block: %v", err)
				// 	continue
				// }

				// timestamp := time.Unix(int64(block.Time()), 0)
				// log.Printf("Timestamp for transfer: %s\n", timestamp)

				timestamp := time.Now()

				transaction := data.Transaction{
					BalanceNumeric:   transactionDetails.Value.String(),
					Sender:           transactionDetails.From.String(),
					Recipient:        transactionDetails.To.String(),
					TransactionHash:  event.Raw.TxHash.String(),
					TransactionIndex: int64(event.Raw.TxIndex),
					BlockNumber:      int64(event.Raw.BlockNumber),
					Timestamp:        timestamp,
				}

				if transaction.BalanceNumeric != "0" {
					database.Trans().Insert(transaction)
				}

			}
		}
	}()

	select {}
}
