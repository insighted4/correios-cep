// Copyright 2023 The Correios CEP Admin Authors
//
// Licensed under the AGPL, Version 3.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.gnu.org/licenses/agpl-3.0.en.html
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
