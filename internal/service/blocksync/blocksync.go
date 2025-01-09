package blocksync

import (
	"context"
	"log"

	"github.com/OctaneAL/ETH-Tracker/internal/config"
	"github.com/OctaneAL/ETH-Tracker/internal/data/pg"
	"github.com/OctaneAL/ETH-Tracker/internal/erc20"
	"github.com/OctaneAL/ETH-Tracker/internal/service/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func FetchMissedBlocks(cfg config.Config) {
	// https_endpoint := cfg.GetInfuraHttpsEndpoint()
	// apiKey := cfg.GetInfuraAPIKey()
	// httpsURL := https_endpoint + apiKey

	// tokenAddress := cfg.GetApiTokenAddress()
	// contractAddress := common.HexToAddress(tokenAddress)

	// client_https, err := ethclient.Dial(httpsURL)
	// if err != nil {
	// 	log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	// }
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
		log.Fatalf("Failed to get latest block header: %v", err)
	}

	log.Printf("The latest block number is: %d\n", header.Number.Uint64())

	var blockFetchEnd uint64 = header.Number.Uint64()

	filterer, err := erc20.NewStorageFilterer(contractAddress, client_https)
	if err != nil {
		log.Fatalf("Failed to create filterer: %v", err)
	}

	// Failed to filter Transfer events: query returned more than 10000 results. Try with this block range [0x149524A, 0x1495432].
	if (blockFetchEnd - uint64(blockFetchStart)) >= 100 {
		blockFetchStart = int64(blockFetchEnd) - 100 + 1
	}

	filterOpts := &bind.FilterOpts{
		Start: uint64(blockFetchStart),
		End:   &blockFetchEnd,
	}

	log.Printf("Starting at block %d, ending at block %d\n", blockFetchStart, blockFetchEnd)

	iter, err := filterer.FilterTransfer(filterOpts, nil, nil)
	if err != nil {
		log.Fatalf("Failed to filter Transfer events: %v", err)
	}

	log.Println("Transfer Events:")
	for iter.Next() {
		event := iter.Event

		utils.ProcessTransferEvent(event, filterer, database)
	}
}
