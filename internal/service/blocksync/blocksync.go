package blocksync

import (
	"context"

	"github.com/OctaneAL/ETH-Tracker/internal/config"
	"github.com/OctaneAL/ETH-Tracker/internal/data/pg"
	"github.com/OctaneAL/ETH-Tracker/internal/erc20"
	"github.com/OctaneAL/ETH-Tracker/internal/models"
	"github.com/OctaneAL/ETH-Tracker/internal/service/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func FetchMissedBlocks(cfg config.Config) {
	logger := cfg.Log()

	client_https := cfg.GetHttpsClient()

	contractAddress := cfg.GetContractAddress()

	database := pg.NewMasterQ(cfg.DB())

	blockFetchStart := cfg.GetApiBlockFetchStart()

	lastTransaction, err := database.Trans().GetLastRecord()
	if err == nil {
		blockFetchStart = lastTransaction.BlockNumber
	}
	header, err := client_https.HeaderByNumber(context.Background(), nil)
	if err != nil {
		logger.Fatalf("Failed to get latest block header: %v", err)
	}

	logger.Infof("The latest block number is: %d\n", header.Number.Uint64())

	var blockFetchEnd uint64 = header.Number.Uint64()

	filterer, err := erc20.NewStorageFilterer(contractAddress, client_https)
	if err != nil {
		logger.Fatalf("Failed to create filterer: %v", err)
	}

	// Failed to filter Transfer events: query returned more than 10000 results. Try with this block range [0x149524A, 0x1495432].
	if (blockFetchEnd - uint64(blockFetchStart)) >= 100 {
		blockFetchStart = int64(blockFetchEnd) - 100 + 1
	}

	filterOpts := &bind.FilterOpts{
		Start: uint64(blockFetchStart),
		End:   &blockFetchEnd,
	}

	logger.Infof("Starting at block %d, ending at block %d\n", blockFetchStart, blockFetchEnd)

	iter, err := filterer.FilterTransfer(filterOpts, nil, nil)
	if err != nil {
		logger.Fatalf("Failed to filter Transfer events: %v", err)
	}

	blockHash := models.BlockHash{
		BlockNumber: nil,
		Timestamp:   nil,
	}

	logger.Info("Transfer Events:")
	for iter.Next() {
		event := iter.Event

		utils.ProcessTransferEvent(event, filterer, database, &blockHash, client_https, logger)
	}
}
