package cli

import (
	"context"
	"log"
	"sync"

	"github.com/OctaneAL/ETH-Tracker/internal/config"
	"github.com/OctaneAL/ETH-Tracker/internal/db"
	"github.com/OctaneAL/ETH-Tracker/internal/service"
	"github.com/OctaneAL/ETH-Tracker/internal/service/handlers"
	"github.com/OctaneAL/ETH-Tracker/internal/service/websocket"
	"github.com/alecthomas/kingpin"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3"
)

func RunServiceCommand(cfg config.Config) {
	var wg sync.WaitGroup

	// ctx, _ := context.WithCancel(context.Background())

	ctx := handlers.CtxDB(context.Background(), db.NewDB(cfg.DatabaseURL()))(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Starting service...")
		service.Run(cfg)
		log.Println("Service stopped.")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Starting WebSocket subscription...")
		websocket.SubscribeToLogs(ctx, cfg)
		log.Println("WebSocket subscription stopped.")
	}()

	wg.Wait()
}

func Run(args []string) bool {
	log := logan.New()

	defer func() {
		if rvr := recover(); rvr != nil {
			log.WithRecover(rvr).Error("app panicked")
		}
	}()

	cfg := config.New(kv.MustFromEnv())
	log = cfg.Log()

	app := kingpin.New("ETH-Tracker", "")

	runCmd := app.Command("run", "run command")
	serviceCmd := runCmd.Command("service", "run service") // you can insert custom help

	migrateCmd := app.Command("migrate", "migrate command")
	migrateUpCmd := migrateCmd.Command("up", "migrate db up")
	migrateDownCmd := migrateCmd.Command("down", "migrate db down")

	// custom commands go here...

	cmd, err := app.Parse(args[1:])
	if err != nil {
		log.WithError(err).Error("failed to parse arguments")
		return false
	}

	switch cmd {
	case serviceCmd.FullCommand():
		RunServiceCommand(cfg)
		// service.Run(cfg)
	case migrateUpCmd.FullCommand():
		err = MigrateUp(cfg)
	case migrateDownCmd.FullCommand():
		err = MigrateDown(cfg)
	// handle any custom commands here in the same way
	default:
		log.Errorf("unknown command %s", cmd)
		return false
	}
	if err != nil {
		log.WithError(err).Error("failed to exec cmd")
		return false
	}
	return true
}
