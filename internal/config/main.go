package config

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/copus"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type Config interface {
	comfig.Logger
	pgdb.Databaser
	types.Copuser
	comfig.Listenerer
	DatabaseURL() string
	GetApiBlockFetchStart() int64
	GetWsClient() *ethclient.Client
	GetHttpsClient() *ethclient.Client
	GetContractAddress() common.Address
}

type config struct {
	comfig.Logger
	pgdb.Databaser
	types.Copuser
	comfig.Listenerer
	getter           kv.Getter
	client_ws        *ethclient.Client
	client_https     *ethclient.Client
	contract_address common.Address
}

func (c *config) DatabaseURL() string {
	dbMap, err := kv.MustFromEnv().GetStringMap("db")
	if err != nil {
		panic(err)
	}
	return dbMap["url"].(string)
}

func (c *config) GetApiBlockFetchStart() int64 {
	dbMap, err := kv.MustFromEnv().GetStringMap("api")
	if err != nil {
		panic(err)
	}

	return int64(dbMap["block_fetch_start"].(int))
}

func getInfuraWsEndpoint() string {
	dbMap, err := kv.MustFromEnv().GetStringMap("api")
	if err != nil {
		panic(err)
	}
	return dbMap["websocket_endpoint"].(string)
}

func getInfuraHttpsEndpoint() string {
	dbMap, err := kv.MustFromEnv().GetStringMap("api")
	if err != nil {
		panic(err)
	}
	return dbMap["https_endpoint"].(string)
}

func getInfuraAPIKey() string {
	dbMap, err := kv.MustFromEnv().GetStringMap("api")
	if err != nil {
		panic(err)
	}
	return dbMap["key"].(string)
}
func getApiTokenAddress() string {
	dbMap, err := kv.MustFromEnv().GetStringMap("api")
	if err != nil {
		panic(err)
	}
	return dbMap["token_address"].(string)
}
func (c *config) GetWsClient() *ethclient.Client {
	return c.client_ws
}

func (c *config) GetHttpsClient() *ethclient.Client {
	return c.client_https
}

func (c *config) GetContractAddress() common.Address {
	return c.contract_address
}

func New(getter kv.Getter) Config {
	apiKey := getInfuraAPIKey()

	https_endpoint := getInfuraHttpsEndpoint()
	httpsURL := https_endpoint + apiKey

	tokenAddress := getApiTokenAddress()
	contractAddress := common.HexToAddress(tokenAddress)

	client_https, err := ethclient.Dial(httpsURL)
	if err != nil {
		log.Fatalf("Failed to connect to the https Ethereum client: %v", err)
	}

	ws_endpoint := getInfuraWsEndpoint()
	websocketURL := ws_endpoint + apiKey

	client_ws, err := ethclient.Dial(websocketURL)
	if err != nil {
		log.Fatalf("Failed to connect to the WebSocket Ethereum client: %v", err)
	}

	return &config{
		getter:           getter,
		Databaser:        pgdb.NewDatabaser(getter),
		Copuser:          copus.NewCopuser(getter),
		Listenerer:       comfig.NewListenerer(getter),
		Logger:           comfig.NewLogger(getter, comfig.LoggerOpts{}),
		client_ws:        client_ws,
		client_https:     client_https,
		contract_address: contractAddress,
	}
}
