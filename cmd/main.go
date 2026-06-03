package main

import (
	"context"
	"log"

	"github.com/caarlos0/env/v11"
	infra "github.com/el10savio/artistInsights/src/infra/clickhouse"
	"github.com/el10savio/artistInsights/src/lib"
	"github.com/el10savio/artistInsights/src/pkg/logger"
	"github.com/el10savio/artistInsights/src/srv"
)

type config struct {
	CHHost string `env:"CH_HOST" envDefault:"localhost"`
	CHPort uint16 `env:"CH_PORT" envDefault:"9000"`
	CHDB   string `env:"CH_DB"   envDefault:"artist_insights"`
	CHUser string `env:"CH_USER" envDefault:"admin"`
	CHPass string `env:"CH_PASS" envDefault:"pass"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("parse config: %v", err)
	}

	zapLogger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("init logger: %v", err)
	}
	defer zapLogger.Sync()

	ctx := context.Background()
	conn, err := infra.Connect(ctx, infra.Config{
		Host: cfg.CHHost,
		Port: cfg.CHPort,
		DB:   cfg.CHDB,
		User: cfg.CHUser,
		Pass: cfg.CHPass,
	})
	if err != nil {
		log.Fatalf("connect clickhouse: %v", err)
	}

	artistStore := &infra.ArtistStore{Conn: conn}
	artistService := lib.NewArtistService(artistStore)
	server := srv.New(artistService, zapLogger)

	if err := server.Start(":8080"); err != nil {
		log.Fatalf("server: %v", err)
	}
}
