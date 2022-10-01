package main

import (
	"github.com/fadyat/hooks/api"
	"github.com/fadyat/hooks/api/hooks/gitlab"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"time"
)

func main() {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(api.LoggerMiddleware(&log.Logger))
	zerolog.TimeFieldFormat = time.RFC822

	cfg, err := api.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	baseAPI := r.Group("/api/v1")
	baseAPI.Use(api.ConfigMiddleware(cfg))

	baseAPI.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	asanaHooks := baseAPI.Group("/asana")
	asanaHooks.POST("/merge", gitlab.MergeRequestAsana)

	if err = r.Run(":80"); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
