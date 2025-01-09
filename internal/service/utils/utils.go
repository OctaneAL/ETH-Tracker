package utils

import (
	"log"
	"time"

	"github.com/OctaneAL/ETH-Tracker/internal/data"
	"github.com/OctaneAL/ETH-Tracker/internal/erc20"
)

func ProcessTransferEvent(event *erc20.StorageTransfer, filterer *erc20.StorageFilterer, database data.MasterQ) {
	log.Printf("Transfer event received: From %s, To %s, Value %d, Block %d",
		event.From.Hex(), event.To.Hex(), event.Value, event.Raw.BlockNumber)

	transactionDetails, err := filterer.ParseTransfer(event.Raw)
	if err != nil {
		log.Fatalf("Failed to Parse Tranfer: %v", err)
	}

	// Takes too much time

	// blockNumber := big.NewInt(int64(event.Raw.BlockNumber))
	// block, err := client_https.BlockByNumber(context.Background(), blockNumber)
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
