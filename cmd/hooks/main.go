package main

import (
	"github.com/fadyat/hooks/api"
	"github.com/fadyat/hooks/api/config"
	_ "github.com/fadyat/hooks/api/docs"
	"github.com/fadyat/hooks/api/handlers"
	"github.com/fadyat/hooks/api/services/tm"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"time"
)

// PingExample godoc
// @Tags    example
// @Success 200 {string} string "pong"
// @Router  /api/v1/ping [get]
func ping(g *gin.Context) {
	g.JSON(http.StatusOK, "pong")
}

// @title          Hooks API
// @version        1.0.0
// @description    This is a sample server Hooks API.
// @termsOfService https://swagger.io/terms/
// @contact.name   Fadeyev Artyom
// @contact.url    https://github.com/fadyat
// @license.name   MIT
// @license.url    https://github.com/MartinHeinz/go-project-blueprint/blob/master/LICENSE
// @schemes        http https
// @host           localhost:80
// @BasePath       /api/v1
func main() {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(api.LoggerMiddleware(&log.Logger))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	cfg, err := config.Load()
	if err != nil {
		log.Info().Err(err).Msg("Failed to load config")
	}

	setupLogger()
	setupApiV1(r, cfg)

	if err = r.Run(":80"); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}

func setupLogger() {
	zerolog.TimeFieldFormat = time.RFC822
}

func setupApiV1(r *gin.Engine, cfg *config.HTTPAPI) {
	v1 := r.Group("/api/v1")
	v1.GET("/ping", ping)

	gh := handlers.NewGitlabHandler(cfg, tm.NewAsanaService(cfg.AsanaAPIKey))
	v1.POST("/asana/push", gh.UpdateLastCommitInfo)
}
