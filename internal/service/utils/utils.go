package utils

import (
	"context"
	"math/big"
	"time"

	"github.com/OctaneAL/ETH-Tracker/internal/data"
	"github.com/OctaneAL/ETH-Tracker/internal/erc20"
	"github.com/OctaneAL/ETH-Tracker/internal/models"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/logan/v3"
)

func ProcessTransferEvent(event *erc20.StorageTransfer, filterer *erc20.StorageFilterer, database data.MasterQ, blockHash *models.BlockHash, client *ethclient.Client, logger *logan.Entry) {
	logger.Infof("Transfer event received: From %s, To %s, Value %d, Block %d",
		event.From.Hex(), event.To.Hex(), event.Value, event.Raw.BlockNumber)

	transactionDetails, err := filterer.ParseTransfer(event.Raw)
	if err != nil {
		logger.Fatalf("Failed to Parse Tranfer: %v", err)
	}

	// Fetch block timestamp

	blockNumber := big.NewInt(int64(event.Raw.BlockNumber))

	timestamp := time.Now()

	if blockHash.BlockNumber == nil || (blockHash.BlockNumber.Uint64() != blockNumber.Uint64()) {
		block, err := client.BlockByNumber(context.Background(), blockNumber)
		if err != nil {
			logger.Infof("Failed to fetch block: %v", err)
		} else {
			timestamp = time.Unix(int64(block.Time()), 0)
		}

		blockHash.BlockNumber = blockNumber
		blockHash.Timestamp = &timestamp
	}

	transaction := data.Transaction{
		BalanceNumeric:   transactionDetails.Value.String(),
		Sender:           transactionDetails.From.String(),
		Recipient:        transactionDetails.To.String(),
		TransactionHash:  event.Raw.TxHash.String(),
		TransactionIndex: int64(event.Raw.TxIndex),
		BlockNumber:      int64(event.Raw.BlockNumber),
		Timestamp:        *blockHash.Timestamp,
	}

	if transaction.BalanceNumeric != "0" {
		database.Trans().Insert(transaction)
	}
}
