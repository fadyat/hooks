package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"time"
)

func AsanaApiMiddleware(config ApiConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("ApiConfig", config)
		c.Next()
	}
}

func LoggerMiddleware(logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		params := gin.LogFormatterParams{
			Latency:      time.Now().Sub(start),
			Method:       c.Request.Method,
			StatusCode:   c.Writer.Status(),
			ErrorMessage: c.Errors.String(),
			BodySize:     c.Writer.Size(),
			Path:         fmt.Sprintf("%s?%s", c.Request.URL.Path, c.Request.URL.RawQuery),
		}

		var logEvent *zerolog.Event
		if params.StatusCode >= 500 {
			logEvent = logger.Error()
		} else if params.StatusCode >= 400 {
			logEvent = logger.Warn()
		} else {
			logEvent = logger.Info()
		}

		logEvent.Str("method", params.Method).
			Str("path", params.Path).
			Int("status", params.StatusCode).
			Int("size", params.BodySize).
			Dur("latency", params.Latency).
			Msg(params.ErrorMessage)
	}
}
