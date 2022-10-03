package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// EndWithError ends request with error
func EndWithError(c *gin.Context, err error, httpCode int, l *zerolog.Logger) {
	l.Info().Msg(err.Error())
	c.JSON(httpCode, gin.H{"error": err.Error()})
}
