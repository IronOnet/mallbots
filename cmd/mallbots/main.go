package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/irononet/mallbots/baskets"
	"github.com/irononet/mallbots/customers"
	"github.com/irononet/mallbots/depots"
	"github.com/irononet/mallbots/internal/config"
	"github.com/irononet/mallbots/internal/logger"
	"github.com/irononet/mallbots/internal/monolith"
	"github.com/irononet/mallbots/internal/rpc"
	"github.com/irononet/mallbots/internal/waiter"
	"github.com/irononet/mallbots/internal/web"
	"github.com/irononet/mallbots/notifications"
	"github.com/irononet/mallbots/ordering"
	"github.com/irononet/mallbots/payments"
	"github.com/irononet/mallbots/stores"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run() (err error) {
	var cfg config.AppConfig
	// parse config/env/...
	cfg, err = config.InitConfig()
	if err != nil {
		return err
	}

	m := app{cfg: cfg}

	// init infra
	m.db, err = sql.Open("pgx", cfg.PG.Conn)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(m.db)
	// init nats and jetstream
	m.nc, err = nats.Connect(cfg.Nats.URL)
	if err != nil{
		return err
	}
	defer m.nc.Close()

	m.js, err = initJetstream(cfg.Nats, m.nc)
	if err != nil{
		return err

	}
	m.logger = logger.New(logger.LogConfig{
		Environment: cfg.Environment,
		LogLevel:    logger.Level(cfg.LogLevel),
	})
	m.rpc = initRpc(cfg.Rpc)
	m.mux = initMux(cfg.Web)
	m.waiter = waiter.New(waiter.CatchSignals())

	// init modules
	m.modules = []monolith.Module{
		&baskets.Module{},
		&customers.Module{},
		&depots.Module{},
		&notifications.Module{},
		&ordering.Module{},
		&payments.Module{},
		&stores.Module{},
	}

	if err = m.StartupModules(); err != nil {
		return err
	}

	// Mount general web resources
	m.mux.Mount("/", http.FileServer(http.FS(web.WebUI)))

	fmt.Println("Started mallbots application")
	defer fmt.Println("stopped mallbots application")

	m.waiter.Add(
		m.waitForWeb,
		m.waitForRpc,
		m.waitForStream,
	)

	return m.waiter.Wait()
}

func initRpc(_ rpc.RpcConfig) *grpc.Server {
	server := grpc.NewServer()
	reflection.Register(server)

	return server
}

func initMux(_ web.WebConfig) *chi.Mux {
	return chi.NewMux()
}

func initJetstream(cfg config.NatsConfig, nc *nats.Conn) (nats.JetStreamContext, error){
	js, err := nc.JetStream()
	if err != nil{
		return nil, err
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name: cfg.Stream,
		Subjects: []string{fmt.Sprintf("%s.>", cfg.Stream)},
	})

	return js, err
}
