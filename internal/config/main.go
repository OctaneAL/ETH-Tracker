package config

import (
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
	GetInfuraWsEndpoint() string
	GetInfuraHttpsEndpoint() string
	GetInfuraAPIKey() string
	GetApiTokenAddress() string
}

type config struct {
	comfig.Logger
	pgdb.Databaser
	types.Copuser
	comfig.Listenerer
	getter kv.Getter
}

func (c *config) DatabaseURL() string {
	dbMap, err := kv.MustFromEnv().GetStringMap("db")
	if err != nil {
		panic(err)
	}
	return dbMap["url"].(string)
}

func (c *config) GetInfuraWsEndpoint() string {
	dbMap, err := kv.MustFromEnv().GetStringMap("api")
	if err != nil {
		panic(err)
	}
	return dbMap["websocket_endpoint"].(string)
}

func (c *config) GetInfuraHttpsEndpoint() string {
	dbMap, err := kv.MustFromEnv().GetStringMap("api")
	if err != nil {
		panic(err)
	}
	return dbMap["https_endpoint"].(string)
}

func (c *config) GetInfuraAPIKey() string {
	dbMap, err := kv.MustFromEnv().GetStringMap("api")
	if err != nil {
		panic(err)
	}
	return dbMap["key"].(string)
}

func (c *config) GetApiTokenAddress() string {
	dbMap, err := kv.MustFromEnv().GetStringMap("api")
	if err != nil {
		panic(err)
	}
	return dbMap["token_address"].(string)
}

func New(getter kv.Getter) Config {
	return &config{
		getter:     getter,
		Databaser:  pgdb.NewDatabaser(getter),
		Copuser:    copus.NewCopuser(getter),
		Listenerer: comfig.NewListenerer(getter),
		Logger:     comfig.NewLogger(getter, comfig.LoggerOpts{}),
	}
}
