package gitlab

import (
	"github.com/fadyat/hooks/api"
	"github.com/fadyat/hooks/api/entities/gitlab"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"net/http"
)

// SyncMRDescriptionWithAsanaTasks updates the merge request description with the task info
// based on the mentioned tasks in the branch name
//
//	@Description	Update merge request description with the task info
//	@Accept			json
//	@Param			X-Gitlab-Event	header		string					true	"Gitlab event"
//	@Param			X-Gitlab-Token	header		string					true	"Gitlab token"
//	@Param			body			body		gitlab.MergeRequestHook	true	"Gitlab request body"
//	@Success		200				{object}	api.Response
//	@Failure		400				{object}	api.Response
//	@Failure		401				{object}	api.Response
//	@Failure		500				{object}	api.Response
//	@Router			/api/v1/gitlab/sync_description [post]
func (h *Handler) SyncMRDescriptionWithAsanaTasks(c *gin.Context) {
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

		// returning 200 to avoid gitlab retrying the request
		c.JSON(http.StatusOK, api.Response{
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
