package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"time"
)

// LoggerMiddleware is a middleware that injects the logger into the context
func LoggerMiddleware(logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		row := c.Request.URL.RawQuery
		if row != "" {
			path = fmt.Sprintf("%s?%s", path, row)
		}

		logger.Info().
			Str("method", c.Request.Method).
			Str("path", path).
			Msg("Request started")

		c.Next()

		params := gin.LogFormatterParams{
			Latency:      time.Since(start),
			Method:       c.Request.Method,
			StatusCode:   c.Writer.Status(),
			ErrorMessage: c.Errors.String(),
			Path:         path,
		}

		var logEvent *zerolog.Event
		switch code := params.StatusCode; {
		case code >= 500:
			logEvent = logger.Error()
		default:
			logEvent = logger.Info()
		}

		logEvent.Str("method", params.Method).
			Int("status", params.StatusCode).
			Str("path", params.Path).
			Dur("latency", params.Latency).
			Msg(params.ErrorMessage)
	}
}
