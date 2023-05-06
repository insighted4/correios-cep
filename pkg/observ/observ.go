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
	"context"
	"fmt"

	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/stackdriver"
	datadog "github.com/DataDog/opencensus-go-exporter-datadog"
	"github.com/insighted4/correios-cep/pkg/errors"
	"go.opencensus.io/trace"
)

// RegisterExporter determines the type of TraceExporter service for exporting traces from opencensus
// User can choose from multiple tracing services (datadog, jaegar)
// RegisterExporter returns the 'Flush' function for that particular tracing service.
func RegisterExporter(traceExporter, url, service, env string) (func(), error) {
	const op errors.Op = "observ.RegisterExporter"
	switch traceExporter {
	case "jaeger":
		return registerJaegerExporter(url, service, env)
	case "datadog":
		return registerDatadogExporter(url, service, env)
	case "stackdriver":
		return registerStackdriverExporter(url, env)
	case "":
		return nil, errors.E(op, "Exporter not specified. Traces won't be exported")
	default:
		return nil, errors.E(op, fmt.Sprintf("Exporter %s not supported. Please open PR or an issue at github.com/gomods/athens", traceExporter))
	}
}

// registerJaegerExporter creates a jaeger exporter for exporting traces to opencensus.
// Currently uses the 'TraceExporter' variable in the config file.
// It should in the future have a nice sampling rate defined.
func registerJaegerExporter(url, service, env string) (func(), error) {
	const op errors.Op = "observ.registerJaegarExporter"
	if url == "" {
		return nil, errors.E(op, "Exporter URL is empty. Traces won't be exported")
	}
	ex, err := jaeger.NewExporter(jaeger.Options{
		Endpoint: url,
		Process: jaeger.Process{
			ServiceName: service,
			Tags: []jaeger.Tag{
				// IP Tag ensures Jaeger's clock isn't skewed.
				// If/when we have traces across different servers,
				// we should make this IP dynamic.
				jaeger.StringTag("ip", "127.0.0.1"),
			},
		},
	})
	if err != nil {
		return nil, errors.E(op, err)
	}
	traceRegisterExporter(ex, env)
	return ex.Flush, nil
}

func traceRegisterExporter(exporter trace.Exporter, env string) {
	trace.RegisterExporter(exporter)
	if env == "development" {
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	}
}

// registerDatadogTracerExporter creates a datadog exporter.
// Currently uses the 'TraceExporter' variable in the config file.
func registerDatadogExporter(url, service, env string) (func(), error) {
	const op errors.Op = "observ.registerDatadogExporter"
	ex, err := datadog.NewExporter(
		datadog.Options{
			TraceAddr: url,
			Service:   service,
		})
	if err != nil {
		return nil, errors.E(op, err)
	}
	traceRegisterExporter(ex, env)
	return ex.Stop, nil
}

func registerStackdriverExporter(projectID, env string) (func(), error) {
	const op errors.Op = "observ.registerStackdriverExporter"
	ex, err := stackdriver.NewExporter(stackdriver.Options{ProjectID: projectID})
	if err != nil {
		return nil, errors.E(op, err)
	}
	traceRegisterExporter(ex, env)
	return ex.Flush, nil
}

// StartSpan takes in a Context Interface and opName and starts a span. It returns the new attached ObserverContext
// and span.
func StartSpan(ctx context.Context, op string) (context.Context, *trace.Span) {
	return trace.StartSpan(ctx, op)
}
