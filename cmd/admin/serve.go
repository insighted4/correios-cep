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
	"context"
	"fmt"
	"os"
	"time"

	"github.com/insighted4/correios-cep/pkg/log"
	"github.com/insighted4/correios-cep/pkg/net"
	"github.com/insighted4/correios-cep/server"
	"github.com/insighted4/correios-cep/storage/postgres"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func commandServe() *cobra.Command {
	var (
		databaseURL string
		logFormat   string
		logLevel    string
		addr        string
	)

	cmd := cobra.Command{
		Use:     "serve",
		Short:   "Start HTTP server",
		Example: fmt.Sprintf("%s serve", shortDescription),
		Run: func(cmd *cobra.Command, args []string) {
			if err := serve(); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVar(&databaseURL, "database-url", "", "database connection string (ex.: postgresql://username:password@localhost/cep)")
	_ = viper.BindPFlag("database_url", cmd.Flags().Lookup("database-url"))

	cmd.Flags().StringVar(&logFormat, "log-format", log.DefaultFormat, "logger format")
	_ = viper.BindPFlag("log_format", cmd.Flags().Lookup("log-format"))

	cmd.Flags().StringVar(&logLevel, "log-level", log.DefaultLevel, "logger level")
	_ = viper.BindPFlag("log_level", cmd.Flags().Lookup("log-level"))

	cmd.Flags().StringVar(&addr, "addr", net.DefaultAddr, "HTTP bind address")
	_ = viper.BindPFlag("addr", cmd.Flags().Lookup("addr"))

	return &cmd
}

func serve() error {
	log.SetLogger(newLogger())

	pgCfg, err := newPostgresOptions()
	if err != nil {
		return err
	}

	now := time.Now().UTC

	ctx := context.Background()

	pg, err := postgres.Connect(ctx, pgCfg, now)
	if err != nil {
		return err
	}

	cfg := newServerConfig(pg)
	cfg.Now = now

	s := server.New(cfg)
	if err := s.Run(); err != nil {
		return err
	}

	log.Infof("Finished")
	return nil
}
