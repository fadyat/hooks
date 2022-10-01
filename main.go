package main

import (
	"github.com/fadyat/gitlab-hooks/app"
	"github.com/fadyat/gitlab-hooks/app/hooks"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	logs "github.com/rs/zerolog/log"
	"time"
)

func main() {
	r := gin.New()
	r.Use(gin.Recovery())
	zerolog.TimeFieldFormat = time.RFC822
	r.Use(app.LoggerMiddleware(&logs.Logger))

	cfg, err := app.LoadConfig()
	if err != nil {
		logs.Fatal().Err(err).Msg("Failed to load config")
	}

	api := r.Group("/api/v1")
	api.Use(app.ConfigMiddleware(cfg))

	api.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	asanaHooks := api.Group("/asana")
	asanaHooks.POST("/merge", hooks.GitlabMergeRequestAsana)

	if err = r.Run(":80"); err != nil {
		logs.Fatal().Err(err).Msg("Failed to start server")
	}
}
