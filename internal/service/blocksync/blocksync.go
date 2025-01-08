package blocksync

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

func FetchMissedBlocks(cfg config.Config) {
	https_endpoint := cfg.GetInfuraHttpsEndpoint()
	apiKey := cfg.GetInfuraAPIKey()
	httpsURL := https_endpoint + apiKey

	tokenAddress := cfg.GetApiTokenAddress()
	contractAddress := common.HexToAddress(tokenAddress)

	client_https, err := ethclient.Dial(httpsURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	database := pg.NewMasterQ(cfg.DB())

	blockFetchStart := cfg.GetApiBlockFetchStart()

	lastTransaction, err := database.Trans().GetLastRecord()
	if err == nil {
		blockFetchStart = lastTransaction.BlockNumber + 1
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
}
