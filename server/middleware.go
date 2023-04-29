package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggerMiddleware returns a gin.HandlerFunc (middleware) that logs requests using logrus.
//
// Requests with errors are logged using logrus.Error().
// Requests without errors are logged using logrus.Info().
//
// It receives:
//  1. A time package format string (e.g. time.RFC3339).
//  2. A boolean stating whether to use UTC time zone or local.
func LoggerMiddleware(logger logrus.FieldLogger, now func() time.Time, timeFormat string, utc bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		c.Next()

		end := now()
		latency := end.Sub(start)
		if utc {
			end = end.UTC()
		}

		entry := logger.WithFields(logrus.Fields{
			"status":       c.Writer.Status(),
			"method":       c.Request.Method,
			"uri":          c.Request.RequestURI,
			"path":         path,
			"content_type": c.ContentType(),
			"remote-addr":  c.ClientIP(),
			"user-agent":   c.Request.UserAgent(),
			"x-request-id": c.GetHeader("X-Request-Id"),
			"latency":      latency,
			"time":         end.Format(timeFormat),
		})

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			entry.Error(c.Errors.String())
		} else {
			entry.Info()
		}
	}
}
