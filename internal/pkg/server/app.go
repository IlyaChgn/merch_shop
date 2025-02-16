package server

import (
	"context"
	authrepo "github.com/IlyaChgn/merch_shop/internal/pkg/auth/repository"
	"github.com/IlyaChgn/merch_shop/internal/pkg/config"
	routers "github.com/IlyaChgn/merch_shop/internal/pkg/server/delivery/router"
	pool "github.com/IlyaChgn/merch_shop/internal/pkg/server/repository"
	shoprepo "github.com/IlyaChgn/merch_shop/internal/pkg/shop/repository"
	"log"
	"net/http"
	"os"
	"time"
)

type Server struct {
	server *http.Server
}

type serverConfig struct {
	Address string
	Timeout time.Duration
	Handler http.Handler
}

func createServerConfig(addr string, timeout int, handler http.Handler) serverConfig {
	return serverConfig{
		Address: addr,
		Timeout: time.Second * time.Duration(timeout),
		Handler: handler,
	}
}

func createServer(config serverConfig) *http.Server {
	return &http.Server{
		Addr:         config.Address,
		ReadTimeout:  config.Timeout,
		WriteTimeout: config.Timeout,
		Handler:      config.Handler,
	}
}

func (srv *Server) Run() error {
	cfgPath := os.Getenv("CONFIG_PATH")

	cfg := config.ReadConfig(cfgPath)
	if cfg == nil {
		log.Fatal("The config wasn`t opened")
	}

	postgresURL := pool.NewConnectionString(cfg.Postgres.Username, cfg.Postgres.Password,
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName)

	postgresPool, err := pool.NewPostgresPool(postgresURL)
	if err != nil {
		log.Fatal("Something went wrong while creating postgres pool ", err)
	}

	err = postgresPool.Ping(context.Background())
	if err != nil {
		log.Fatal("Cannot connect to postgres database ", err)
	}

	authStorage := authrepo.NewAuthStorage(postgresPool, cfg.SecretKey)
	shopStorage := shoprepo.NewShopStorage(postgresPool)

	router := routers.NewRouter(authStorage, shopStorage)

	serverCfg := createServerConfig(cfg.Server.Host+cfg.Server.Port, cfg.Server.Timeout, router)
	srv.server = createServer(serverCfg)

	log.Printf("Server is listening on %s\n", cfg.Server.Host+cfg.Server.Port)

	return srv.server.ListenAndServe()
}
