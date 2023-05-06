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

package handler

import (
	"fmt"
	"net/http"
	"regexp"
	"runtime"
	"strings"
	"time"

	gosundheit "github.com/AppsFlyer/go-sundheit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/insighted4/correios-cep/correios"
	"github.com/insighted4/correios-cep/pkg/app"
	"github.com/insighted4/correios-cep/pkg/errors"
	"github.com/insighted4/correios-cep/pkg/log"
	"github.com/insighted4/correios-cep/pkg/version"
	"github.com/insighted4/correios-cep/storage"
	"github.com/sirupsen/logrus"
)

const Prefix = "/api/v1"

func New(correios correios.Correios, storage storage.Storage, health gosundheit.Health, release bool, now func() time.Time) http.Handler {
	if release {
		gin.SetMode(gin.ReleaseMode)
	}

	logger := log.WithField("component", "handler")

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(cors.Default())
	router.Use(LoggerMiddleware(logger, now, time.RFC3339, true))
	router.NoRoute(notFoundHandler())

	router.GET("/", rootHandler())
	router.GET("/health", healthHandler(health, logger))
	router.GET("/ping", pingHandler())

	api := router.Group(Prefix)
	api.GET("/addresses", listAddressHandler(storage, logger))
	api.GET("/addresses/:cep", getAddressHandler(correios, storage, logger))

	return router
}

func rootHandler() gin.HandlerFunc {
	var root = gin.H{
		"server":          app.Description,
		"arch":            runtime.GOARCH,
		"build_time":      version.BuildTime,
		"commit":          version.CommitHash,
		"os":              runtime.GOOS,
		"runtime_version": runtime.Version(),
		"version":         version.Version,
	}

	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, root)
	}
}

func notFoundHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		abortWithStatus(ctx, http.StatusNotFound, http.StatusText(http.StatusNotFound), nil)
	}
}

func healthHandler(health gosundheit.Health, logger logrus.FieldLogger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		result, healthy := health.Results()
		if !healthy {
			abortWithStatus(ctx, http.StatusInternalServerError, "Health check failed.", result)
			return
		}

		logger.Infof("Health check passed")
		ctx.JSON(http.StatusOK, result)
	}
}

func pingHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	}
}

type HTTPErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func (er *HTTPErrorResponse) Error() string {
	return fmt.Sprintf("%d - %s", er.Code, er.Message)
}

func abortWithStatus(ctx *gin.Context, code int, message string, details interface{}) {
	ctx.AbortWithStatusJSON(code, &HTTPErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	})
}

var newLine = regexp.MustCompile(`\r?\n?\t`)

func abortWithError(ctx *gin.Context, err error, details interface{}) {
	var code int
	var msg string
	var typedError errors.Error
	switch {
	case errors.AsErr(err, &typedError):
		code = typedError.Kind
		msg = newLine.ReplaceAllString(typedError.Error(), " ")
		if index := strings.Index(msg, ":"); len(msg) > index+1 {
			msg = strings.TrimSpace(msg[index+1:])
		}
	default:
		code = errors.KindUnexpected
		msg = http.StatusText(http.StatusInternalServerError)
	}

	ctx.AbortWithStatusJSON(code, &HTTPErrorResponse{
		Code:    code,
		Message: msg,
		Details: details,
	})
}
