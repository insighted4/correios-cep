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

package observ

import (
	"fmt"
	"time"

	"contrib.go.opencensus.io/exporter/prometheus"
	"contrib.go.opencensus.io/exporter/stackdriver"
	datadog "github.com/DataDog/opencensus-go-exporter-datadog"
	"github.com/gin-gonic/gin"
	"github.com/insighted4/correios-cep/pkg/errors"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
)

// RegisterStatsExporter determines the type of StatsExporter service for exporting stats from Opencensus
// Currently it supports: prometheus.
func RegisterStatsExporter(r gin.IRouter, statsExporter, service string) (func(), error) {
	const op errors.Op = "observ.RegisterStatsExporter"
	stop := func() {}
	var err error
	switch statsExporter {
	case "prometheus":
		if err := registerPrometheusExporter(r, service); err != nil {
			return nil, errors.E(op, err)
		}
	case "stackdriver":
		if stop, err = registerStatsStackDriverExporter(service); err != nil {
			return nil, errors.E(op, err)
		}
	case "datadog":
		if stop, err = registerStatsDataDogExporter(service); err != nil {
			return nil, errors.E(op, err)
		}
	case "":
		return nil, errors.E(op, "StatsExporter not specified. Stats won't be collected")
	default:
		return nil, errors.E(op, fmt.Sprintf("StatsExporter %s not supported. Please open PR or an issue at github.com/gomods/athens", statsExporter))
	}
	if err = registerViews(); err != nil {
		return nil, errors.E(op, err)
	}

	return stop, nil
}

// registerPrometheusExporter creates exporter that collects stats for Prometheus.
func registerPrometheusExporter(r gin.IRouter, service string) error {
	const op errors.Op = "observ.registerPrometheusExporter"
	prom, err := prometheus.NewExporter(prometheus.Options{
		Namespace: service,
	})
	if err != nil {
		return errors.E(op, err)
	}

	r.GET("/metrics", gin.WrapH(prom))

	view.RegisterExporter(prom)

	return nil
}

func registerStatsDataDogExporter(service string) (func(), error) {
	const op errors.Op = "observ.registerStatsDataDogExporter"

	dd, err := datadog.NewExporter(datadog.Options{Service: service})
	if err != nil {
		return nil, errors.E(op, "Failed to initialize Datadog exporter", err)
	}

	if dd == nil {
		return nil, errors.E(op, "Failed to initialize Datadog exporter")
	}

	view.RegisterExporter(dd)
	return dd.Stop, nil
}

func registerStatsStackDriverExporter(projectID string) (func(), error) {
	const op errors.Op = "observ.registerStatsStackDriverExporter"

	sd, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: projectID,
	})
	if err != nil {
		return nil, errors.E(op, err)
	}

	view.RegisterExporter(sd)
	view.SetReportingPeriod(60 * time.Second)

	return sd.Flush, nil
}

// registerViews register stats which should be collected in Athens.
func registerViews() error {
	const op errors.Op = "observ.registerViews"
	if err := view.Register(
		ochttp.ServerRequestCountView,
		ochttp.ServerResponseBytesView,
		ochttp.ServerLatencyView,
		ochttp.ServerResponseCountByStatusCode,
		ochttp.ServerRequestBytesView,
		ochttp.ServerRequestCountByMethod,
		ochttp.ClientReceivedBytesDistribution,
		ochttp.ClientRoundtripLatencyDistribution,
		ochttp.ClientCompletedCount,
	); err != nil {
		return errors.E(op, err)
	}

	return nil
}
