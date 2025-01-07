package websocket

import (
	"context"
	"fmt"
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
	fmt.Println(websocketURL)

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

	fmt.Println("Good shit - good going")

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

			// transactionDetails := getTransactionByAddress()
			// transactionDetails := getTransactionByHash(eventData.Params.Result, httpsURL)

			// log.Printf("Transaction details: %+v\n", transactionDetails)

			// if transactionDetails != nil {
			// 	database.Trans().Insert(*transactionDetails)
			// }

		}
	}

	// conn, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	// if err != nil {
	// 	log.Fatalf("Failed to connect to WebSocket: %v", err)
	// }
	// defer conn.Close()

	// subscriptionData := Request{
	// 	Jsonrpc: "2.0",
	// 	ID:      1,
	// 	Method:  "eth_subscribe",
	// 	Params: []interface{}{
	// 		"newPendingTransactions",
	// 		map[string]interface{}{
	// 			"address": tokenAddress, // Filter by USDC address
	// 		},
	// 	},
	// }

	// subscriptionJSON, err := json.Marshal(subscriptionData)
	// if err != nil {
	// 	log.Fatalf("Failed to marshal subscription data: %v", err)
	// }

	// err = conn.WriteMessage(websocket.TextMessage, subscriptionJSON)
	// if err != nil {
	// 	log.Fatalf("Failed to send subscription request: %v", err)
	// }
	// log.Println("Subscription request sent.")

	// _, confirmationResponse, err := conn.ReadMessage()
	// if err != nil {
	// 	log.Fatalf("Failed to read confirmation response: %v", err)
	// }

	// var confirmationData SubscriptionResponse
	// err = json.Unmarshal(confirmationResponse, &confirmationData)
	// if err != nil {
	// 	log.Fatalf("Failed to unmarshal confirmation response: %v", err)
	// }
	// log.Printf("Subscription confirmed: %+v\n", confirmationData)

	// if confirmationData.Result == nil {
	// 	log.Println("Failed to create subscription.")
	// 	return
	// }

	// log.Println("Listening for events...")

	// for {
	// 	_, eventResponse, err := conn.ReadMessage()
	// 	if err != nil {
	// 		// log.Printf("Failed to read event data: %v\n", err)
	// 		return
	// 	}

	// 	var eventData EventData
	// 	err = json.Unmarshal(eventResponse, &eventData)
	// 	if err != nil {
	// 		// log.Printf("Failed to unmarshal event data: %v\n", err)
	// 		continue
	// 	}

	// 	// log.Printf("Received event data: %+v\n", eventData)

	// 	transactionDetails := getTransactionByHash(eventData.Params.Result, httpsURL)

	// 	// log.Printf("Transaction details: %+v\n", transactionDetails)

	// 	if transactionDetails != nil {
	// 		database.Trans().Insert(*transactionDetails)
	// 	}
	// }
}

// func getTransactionByHash(txHash, httpsURL string) *data.InsertTransaction {
// 	requestData := Request{
// 		Jsonrpc: "2.0",
// 		ID:      1,
// 		Method:  "eth_getTransactionByHash",
// 		Params: []interface{}{
// 			txHash,
// 		},
// 	}

// 	requestJSON, err := json.Marshal(requestData)
// 	if err != nil {
// 		// log.Fatalf("Failed to marshal subscription data: %v", err)
// 	}

// 	resp, err := http.Post(httpsURL, "application/json", bytes.NewBuffer(requestJSON))
// 	if err != nil {
// 		// log.Fatalf("Failed to send HTTP request: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	var result map[string]interface{}
// 	err = json.NewDecoder(resp.Body).Decode(&result)
// 	if err != nil {
// 		// log.Fatalf("Failed to decode response: %v", err)
// 	}

// 	if result["result"] == nil {
// 		// log.Println("Failed to get transaction details.")
// 		return nil
// 	}

// 	result = result["result"].(map[string]interface{})

// 	balanceRaw := result["value"].(string)

// 	balanceNumeric := new(big.Int)
// 	_, success := balanceNumeric.SetString(balanceRaw[2:], 16)
// 	if !success {
// 		// log.Fatalf("Failed to parse balance: %s", balanceRaw)
// 	}

// 	if (balanceNumeric).Cmp(big.NewInt(0)) == 0 {
// 		// log.Println("Transaction value is 0. Skipping...")
// 		return nil
// 	}

// 	// log.Println(result)

// 	sender := result["from"].(string)
// 	recipient := result["to"].(string)
// 	transactionHash := result["hash"].(string)
// 	transactionIndex := "0x0"
// 	if result["transactionIndex"] != nil {
// 		transactionIndex = result["transactionIndex"].(string)
// 	}

// 	transactionData := data.InsertTransaction{
// 		BalanceNumeric:   balanceNumeric.String(),
// 		Sender:           sender,
// 		Recipient:        recipient,
// 		TransactionHash:  transactionHash,
// 		TransactionIndex: transactionIndex,
// 	}

// 	return &transactionData
// }
