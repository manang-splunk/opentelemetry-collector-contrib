// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package healthcheckextension // import "github.com/open-telemetry/opentelemetry-collector-contrib/extension/healthcheckextension"

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/jaegertracing/jaeger/pkg/healthcheck"
	"go.opencensus.io/stats/view"
	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
)

type healthCheckExtension struct {
	config   Config
	logger   *zap.Logger
	state    *healthcheck.HealthCheck
	server   http.Server
	stopCh   chan struct{}
	exporter *healthCheckExporter
}

var _ component.PipelineWatcher = (*healthCheckExtension)(nil)

func (hc *healthCheckExtension) Start(_ context.Context, host component.Host) error {

	hc.logger.Info("Starting health_check extension", zap.Any("config", hc.config))

	// Initialize listener
	var (
		ln  net.Listener
		err error
	)
	if hc.config.Port != 0 && hc.config.TCPAddr.Endpoint == defaultEndpoint {
		hc.logger.Warn("`Port` is deprecated, use `Endpoint` instead")
		portStr := ":" + strconv.Itoa(int(hc.config.Port))
		ln, err = net.Listen("tcp", portStr)
	} else {
		ln, err = hc.config.TCPAddr.Listen()
	}
	if err != nil {
		return err
	}

	if !hc.config.CheckCollectorPipeline.Enabled {
		// Mount HC handler
		mux := http.NewServeMux()
		mux.Handle(hc.config.Path, hc.state.Handler())
		hc.server.Handler = mux
		hc.stopCh = make(chan struct{})
		go func() {
			defer close(hc.stopCh)

			// The listener ownership goes to the server.
			if err = hc.server.Serve(ln); err != http.ErrServerClosed && err != nil {
				host.ReportFatalError(err)
			}
		}()
	} else {
		// collector pipeline health check
		hc.exporter = newHealthCheckExporter()
		view.RegisterExporter(hc.exporter)

		interval, err := time.ParseDuration(hc.config.CheckCollectorPipeline.Interval)
		if err != nil {
			return err
		}

		// ticker used by collector pipeline health check for rotation
		ticker := time.NewTicker(time.Second)

		mux := http.NewServeMux()
		mux.Handle(hc.config.Path, hc.handler())
		hc.server.Handler = mux
		hc.stopCh = make(chan struct{})
		go func() {
			defer close(hc.stopCh)
			defer view.UnregisterExporter(hc.exporter)

			go func() {
				for {
					select {
					case <-ticker.C:
						hc.exporter.rotate(interval)
					case <-hc.stopCh:
						return
					}
				}
			}()

			if err = hc.server.Serve(ln); err != http.ErrServerClosed && err != nil {
				host.ReportFatalError(err)
			}

		}()
	}

	return nil
}

// new handler function used for check collector pipeline
func (hc *healthCheckExtension) handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if hc.check() && hc.state.Get() == healthcheck.Ready {
			w.WriteHeader(200)
		} else {
			w.WriteHeader(500)
		}
	})
}

func (hc *healthCheckExtension) check() bool {
	return hc.exporter.checkHealthStatus(hc.config.CheckCollectorPipeline.ExporterFailureThreshold)
}

func (hc *healthCheckExtension) Shutdown(context.Context) error {
	err := hc.server.Close()
	if hc.stopCh != nil {
		<-hc.stopCh
	}
	return err
}

func (hc *healthCheckExtension) Ready() error {
	hc.state.Set(healthcheck.Ready)
	return nil
}

func (hc *healthCheckExtension) NotReady() error {
	hc.state.Set(healthcheck.Unavailable)
	return nil
}

func newServer(config Config, logger *zap.Logger) *healthCheckExtension {
	hc := &healthCheckExtension{
		config: config,
		logger: logger,
		state:  healthcheck.New(),
		server: http.Server{},
	}

	hc.state.SetLogger(logger)

	return hc
}
