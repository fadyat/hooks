package gitlab

import (
	"net/http"
	"slices"
	"sort"

	"github.com/fadyat/hooks/api"
	"github.com/fadyat/hooks/api/config"
	"github.com/fadyat/hooks/api/entities"
	"github.com/fadyat/hooks/api/entities/gitlab"
	"github.com/fadyat/hooks/api/helpers"
	"github.com/fadyat/hooks/api/services/tm"
	"github.com/fadyat/hooks/api/services/vcs"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Handler struct {
	cfg *config.HTTPAPI
	l   *zerolog.Logger
	tm  tm.ITaskManager
	vcs vcs.IVCS
}

func NewHandler(cfg *config.HTTPAPI, l *zerolog.Logger, t tm.ITaskManager, v vcs.IVCS) *Handler {
	return &Handler{
		cfg: cfg,
		l:   l,
		tm:  t,
		vcs: v,
	}
}

// UpdateLastCommitInfo updates the last commit info of a task
//
//	@Description	Update last commit info, in custom field or creating a comment
//	@Accept			json
//	@Param			X-Gitlab-Event	header		string					true	"Gitlab event"
//	@Param			X-Gitlab-Token	header		string					true	"Gitlab token"
//	@Param			body			body		gitlab.PushRequestHook	true	"Gitlab request body"
//	@Success		200				{object}	api.Response
//	@Failure		400				{object}	api.Response
//	@Failure		401				{object}	api.Response
//	@Failure		500				{object}	api.Response
//	@Router			/api/v1/asana/push [post]
func (h *Handler) UpdateLastCommitInfo(c *gin.Context) {
	if c.Request.Header.Get("X-Gitlab-Event") != gitlab.PushEvent {
		c.JSON(http.StatusBadRequest, api.Response{
			Ok:    false,
			Error: "invalid event",
		})
		return
	}

	if !slices.Contains(h.cfg.GitlabSecretTokens, c.Request.Header.Get("X-Gitlab-Token")) {
		c.JSON(http.StatusUnauthorized, api.Response{
			Ok:    false,
			Error: "invalid token",
		})
		return
	}

	var r gitlab.PushRequestHook
	if err := c.ShouldBindJSON(&r); err != nil {
		h.l.Error().Err(err).Msg("invalid request body")
		c.JSON(http.StatusBadRequest, api.Response{
			Ok:    false,
			Error: "invalid request body",
		})
		return
	}

	if len(r.Commits) == 0 {
		h.l.Info().Msg("no commits found")
		c.JSON(http.StatusOK, api.Response{
			Ok:    true,
			Error: "no commits found",
		})
		return
	}

	sort.Slice(r.Commits, func(i, j int) bool {
		return r.Commits[i].Timestamp.After(r.Commits[j].Timestamp)
	})

	lastCommit := r.Commits[0]
	if err := h.tm.UpdateLastCommitInfo(entities.Message{
		Text:       lastCommit.Message,
		URL:        lastCommit.URL,
		Author:     lastCommit.Author.Name,
		BranchName: helpers.GetBranchNameFromRef(r.Ref),
	}); err != nil {
		h.l.Error().Err(err).Msg("failed to update last commit info")
		c.JSON(api.GetErrStatusCode(err), api.Response{
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
