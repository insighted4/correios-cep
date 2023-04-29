package main

import (
	"fmt"

	"github.com/insighted4/correios-cep/pkg/log"
	"github.com/insighted4/correios-cep/pkg/net"
	"github.com/insighted4/correios-cep/server"
	"github.com/insighted4/correios-cep/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func newLogger() logrus.FieldLogger {
	return log.New(viper.GetString("log_level"), viper.GetString("log_format"))
}

func newServerConfig(storage storage.Storage) server.Config {
	return server.Config{
		HTTPServerConfig: net.HTTPServerConfig{
			Addr: viper.GetString("addr"),
		},
		ReleaseMode: viper.GetString("log_level") != "debug",
		Storage:     storage,
	}
}

func newPostgresOptions() (*pgxpool.Config, error) {
	databaseURL := viper.GetString("database_url")
	if databaseURL == "" {
		return nil, fmt.Errorf("invalid database URL")
	}

	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
