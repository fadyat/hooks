package main

import (
	"github.com/fadyat/gitlab-hooks/app"
	"github.com/fadyat/gitlab-hooks/app/hooks"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
)

func AsanaApiMiddleware(config app.AsanaConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("AsanaConfig", config)
		c.Next()
	}
}

func JSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()
	r.Use(JSONMiddleware())

	api := r.Group("/api/v1")
	api.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	asanaHooks := api.Group("/asana")
	{
		var cfg app.AsanaConfig
		if err := envconfig.Process("", &cfg); err != nil {
			log.Fatal("Error loading .env file")
		}

		asanaHooks.Use(AsanaApiMiddleware(cfg))
		asanaHooks.POST("/merge_request", hooks.GitlabMergeRequestAsana)
	}

	if err := r.Run(":80"); err != nil {
		log.Fatal(err)
	}
}
