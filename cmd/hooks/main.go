package main

import (
	"fmt"
	"github.com/fadyat/hooks/api"
	"github.com/fadyat/hooks/api/config"
	_ "github.com/fadyat/hooks/api/docs"
	"github.com/fadyat/hooks/api/handlers/gitlab"
	"github.com/fadyat/hooks/api/services/tm"
	"github.com/fadyat/hooks/api/services/vcs"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os"
	"strings"
	"time"
)

// PingExample godoc
//	@Tags		example
//	@Success	200	{string}	string	"pong"
//	@Router		/api/v1/ping [get]
func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

//	@title			Hooks API
//	@version		1.0.0
//	@description	This is a sample server Hooks API.
//	@termsOfService	https://swagger.io/terms/
//	@contact.name	Fadeyev Artyom
//	@contact.url	https://github.com/fadyat
//	@license.name	MIT
//	@license.url	https://github.com/MartinHeinz/go-project-blueprint/blob/master/LICENSE
//	@schemes		http https
//	@BasePath		/api/v1
func main() {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(api.LoggerMiddleware(&log.Logger))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	featureFlags, err := config.LoadFeatureFlags()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load feature flags")
	}

	setupLogger()
	setupAPIV1(r, cfg, featureFlags)

	if featureFlags.IsRepresentSecretsEnabled {
		blurSecrets(&log.Logger)
	}

	addr := fmt.Sprintf(":%s", cfg.Port)
	if err = r.Run(addr); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}

func setupLogger() {
	zerolog.TimeFieldFormat = time.RFC822
}

func setupAPIV1(r *gin.Engine, cfg *config.HTTPAPI, featureFlags *config.FeatureFlags) {
	v1 := r.Group("/api/v1")
	v1.GET("/ping", ping)

	as := tm.NewAsanaService(cfg.AsanaAPIKey, &log.Logger, cfg, featureFlags)
	gs := vcs.NewGitlabService(cfg.GitlabAPIKey, cfg.GitlabDomain, &log.Logger, as)
	gh := gitlab.NewHandler(cfg, &log.Logger, as, gs)

	toAsana := v1.Group("/asana")
	toAsana.POST("/push", gh.UpdateLastCommitInfo)
	toAsana.POST("/merge", gh.OnBranchMerge)

	toGitlab := v1.Group("/gitlab")
	toGitlab.POST("/sync_description", gh.SyncMRDescriptionWithAsanaTasks)
}

func min(a, b int) int {
	if a > b {
		return b
	}

	return a
}

func blurSecrets(lg *zerolog.Logger) {
	blur := "***************"

	for _, env := range os.Environ() {
		sp := strings.Split(env, "=")
		k, v := sp[0], sp[1]
		idx := min(len(v), 3)
		lg.Debug().Msg(fmt.Sprintf("%s=%s", k, v[:idx]+blur[idx:]))
	}
}
