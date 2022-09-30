package main

import (
	"github.com/fadyat/gitlab-hooks/app"
	"github.com/fadyat/gitlab-hooks/app/hooks"
	"github.com/gin-gonic/gin"
	zerolog "github.com/rs/zerolog/log"
)

func main() {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(app.LoggerMiddleware(&zerolog.Logger))

	cfg, err := app.LoadConfig()
	if err != nil {
		zerolog.Fatal().Err(err).Msg("Failed to load config")
	}

	api := r.Group("/api/v1")
	api.Use(app.AsanaApiMiddleware(*cfg))

	api.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	asanaHooks := api.Group("/asana")
	{
		asanaHooks.POST("/merge", hooks.GitlabMergeRequestAsana)
	}

	if err = r.Run(":80"); err != nil {
		zerolog.Fatal().Err(err).Msg("Failed to start server")
	}
}
