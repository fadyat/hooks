package handlers

import (
	"github.com/fadyat/hooks/api"
	"github.com/fadyat/hooks/api/config"
	"github.com/fadyat/hooks/api/entities/gitlab"
	"github.com/fadyat/hooks/api/helpers"
	"github.com/fadyat/hooks/api/services/tm"
	"github.com/fadyat/hooks/api/services/vcs"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"golang.org/x/exp/slices"
	"net/http"
	"sort"
)

type GitlabHandler struct {
	cfg *config.HTTPAPI
	l   *zerolog.Logger
	tm  tm.ITaskManager
	vcs vcs.IVCS
}

func NewGitlabHandler(cfg *config.HTTPAPI, l *zerolog.Logger, m tm.ITaskManager, v vcs.IVCS) *GitlabHandler {
	return &GitlabHandler{
		cfg: cfg,
		l:   l,
		tm:  m,
		vcs: v,
	}
}

// UpdateLastCommitInfo updates the last commit info of a task
// @Description Update last commit info, in custom field or creating a comment
// @Accept      json
// @Param       X-Gitlab-Event header   string                 true "Gitlab event"
// @Param       X-Gitlab-Token header   string                 true "Gitlab token"
// @Param       body           body     gitlab.PushRequestHook true "Gitlab request body"
// @Success     200            {object} api.Response
// @Failure     400            {object} api.Response
// @Failure     401            {object} api.Response
// @Failure     500            {object} api.Response
// @Router      /api/v1/asana/push [post]
func (h *GitlabHandler) UpdateLastCommitInfo(c *gin.Context) {
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
		h.l.Debug().Err(err).Msg("invalid request body")
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

// UpdateMergeRequestDescription updates the merge request description with the task info
// @Description Update merge request description with the task info
// @Accept      json
// @Param       X-Gitlab-Event header   string                  true "Gitlab event"
// @Param       X-Gitlab-Token header   string                  true "Gitlab token"
// @Param       body           body     gitlab.MergeRequestHook true "Gitlab request body"
// @Success     200            {object} api.Response
// @Failure     400            {object} api.Response
// @Failure     401            {object} api.Response
// @Failure     500            {object} api.Response
// @Router      /api/v1/gitlab/merge [post]
func (h *GitlabHandler) UpdateMergeRequestDescription(c *gin.Context) {
	if c.Request.Header.Get("X-Gitlab-Event") != gitlab.MergeEvent {
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

	var r gitlab.MergeRequestHook
	if err := c.ShouldBindJSON(&r); err != nil {
		h.l.Debug().Err(err).Msg("invalid request body")
		c.JSON(http.StatusBadRequest, api.Response{
			Ok:    false,
			Error: "invalid request body",
		})
		return
	}

	if r.ObjectAttributes.Action != gitlab.MergeRequestActionOpen {
		h.l.Debug().Msgf("unsupported action: %s", r.ObjectAttributes.Action)
		c.JSON(http.StatusBadRequest, api.Response{
			Ok:    false,
			Error: "unsupported action, only open is supported",
		})
		return
	}

	attr := r.ObjectAttributes
	err := h.vcs.UpdatePRDescription(r.Project.ID, attr.Iid, attr.SourceBranch, attr.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.Response{
			Ok:    false,
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, api.Response{
		Ok:     true,
		Result: "task mentioned in the description",
	})
}
