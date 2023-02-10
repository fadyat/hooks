package handlers

import (
	"github.com/fadyat/hooks/api"
	"github.com/fadyat/hooks/api/config"
	"github.com/fadyat/hooks/api/entities/gitlab"
	"github.com/fadyat/hooks/api/helpers"
	"github.com/fadyat/hooks/api/services/tm"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"net/http"
	"sort"
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
	if c.Request.Header.Get("X-Gitlab-Event") != gitlab.PushEvent {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if !slices.Contains(h.cfg.GitlabSecretTokens, c.Request.Header.Get("X-Gitlab-Token")) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var r gitlab.PushRequestHook
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, api.Response{
			Ok:    false,
			Error: "invalid request body",
		})
		return
	}

	if len(r.Commits) == 0 {
		c.JSON(http.StatusBadRequest, api.Response{
			Ok:    false,
			Error: "no commits found",
		})
		return
	}

	sort.Slice(r.Commits, func(i, j int) bool {
		return r.Commits[i].Timestamp.After(r.Commits[j].Timestamp)
	})

	err := h.tm.UpdateLastCommitInfo(helpers.GetBranchNameFromRef(r.Ref), r.Commits[0])
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.Response{
			Ok:    false,
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, api.Response{
		Ok:     true,
		Result: "updated",
	})
}
