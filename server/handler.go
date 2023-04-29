package server

import (
	"fmt"
	"net/http"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/insighted4/correios-cep/pkg/app"
	"github.com/insighted4/correios-cep/pkg/errors"
	"github.com/insighted4/correios-cep/pkg/version"
)

const Prefix = "/api/v1"

var root = gin.H{
	"service":         app.Description,
	"arch":            runtime.GOARCH,
	"build_time":      version.BuildTime,
	"commit":          version.CommitHash,
	"os":              runtime.GOOS,
	"runtime_version": runtime.Version(),
	"version":         version.Version,
}

func (s *service) newHandler() http.Handler {
	if s.cfg.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(cors.Default())
	router.Use(LoggerMiddleware(s.logger, s.now, time.RFC3339, true))
	router.NoRoute(s.handleNotFound)

	router.GET("/", s.handleRoot)
	router.GET("/health", s.handleHealth)
	router.GET("/ping", s.handlePing)

	return router
}

func (s *service) handleRoot(c *gin.Context) {
	c.JSON(http.StatusOK, root)
}

func (s *service) handleNotFound(c *gin.Context) {
	s.abortWithStatus(c, http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

func (s *service) handleHealth(c *gin.Context) {
	if s.health.IsHealthy() {
		s.abortWithStatus(c, http.StatusInternalServerError, "Health check failed.")
		return
	}

	s.logger.Infof("Health check passed")
}

func (s *service) handlePing(c *gin.Context) {
	c.Status(http.StatusOK)
}

type HTTPErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (er *HTTPErrorResponse) Error() string {
	return fmt.Sprintf("%d - %s", er.Code, er.Message)
}

func (s *service) abortWithStatus(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, &HTTPErrorResponse{
		Code:    code,
		Message: message,
	})
}

var newLine = regexp.MustCompile(`\r?\n?\t`)

func (s *service) abortWithError(ctx *gin.Context, err error) {
	code := errors.KindUnexpected
	msg := newLine.ReplaceAllString(err.Error(), " ")
	e, ok := err.(*errors.Error)
	if ok {
		code = e.Kind
		if index := strings.Index(msg, ":"); len(msg) > index+1 {
			msg = strings.TrimSpace(msg[index+1:])
		}
	}

	ctx.AbortWithStatusJSON(code, &HTTPErrorResponse{
		Code:    code,
		Message: msg,
	})
}
