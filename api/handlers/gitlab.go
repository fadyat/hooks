package handlers

import (
	"github.com/fadyat/hooks/api"
	"github.com/fadyat/hooks/api/config"
	"github.com/fadyat/hooks/api/services/tm"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GitlabHandler struct {
	cfg *config.HTTPAPI
	tm  tm.ITaskManager
}

func NewGitlabHandler(cfg *config.HTTPAPI, tm tm.ITaskManager) *GitlabHandler {
	return &GitlabHandler{
		cfg: cfg,
		tm:  tm,
	}
}

func (h *GitlabHandler) UpdateLastCommitInfo(c *gin.Context) {
	c.JSON(http.StatusOK, &api.Response{
		Ok:     true,
		Result: "ok",
	})
}
