package main

import (
	"github.com/fadyat/hooks/api"
	_ "github.com/fadyat/hooks/api/docs"
	"github.com/fadyat/hooks/api/hooks/gitlab"
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
// @Accept  json
// @Produce json
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
	zerolog.TimeFieldFormat = time.RFC822

	cfg, err := api.LoadConfig()
	if err != nil {
		log.Info().Err(err).Msg("Failed to load config")
	}

	baseAPI := r.Group("/api/v1")
	baseAPI.Use(api.ConfigMiddleware(cfg))

	baseAPI.GET("/ping", ping)
	asanaHooks := baseAPI.Group("/asana")
	asanaHooks.POST("/merge", gitlab.MergeRequestAsana)
	asanaHooks.POST("/push", gitlab.PushRequestAsana)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	if err = r.Run(":80"); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
