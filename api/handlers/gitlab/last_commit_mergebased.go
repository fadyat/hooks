package gitlab

import (
	"fmt"
	"github.com/fadyat/hooks/api"
	"github.com/fadyat/hooks/api/entities"
	"github.com/fadyat/hooks/api/entities/gitlab"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"net/http"
)

// OnBranchMerge handles the merge event and updates the last commit info
//	@Description	Update last commit info, based on merge event of MR
//	@Accept			json
//	@Param			X-Gitlab-Event	header		string					true	"Gitlab event"
//	@Param			X-Gitlab-Token	header		string					true	"Gitlab token"
//	@Param			body			body		gitlab.PushRequestHook	true	"Gitlab request body"
//	@Success		200				{object}	api.Response
//	@Failure		400				{object}	api.Response
//	@Failure		401				{object}	api.Response
//	@Failure		500				{object}	api.Response
//	@Router			/api/v1/asana/merge [post]
func (h *Handler) OnBranchMerge(c *gin.Context) {
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

	if r.ObjectAttributes.Action != gitlab.MergeRequestActionMerge {
		h.l.Debug().Msgf("unsupported action: %s", r.ObjectAttributes.Action)

		// returning 200 to avoid gitlab retrying the request
		c.JSON(http.StatusOK, api.Response{
			Ok:    false,
			Error: "unsupported action, only merge is supported",
		})
		return
	}

	attr := r.ObjectAttributes
	err := h.tm.UpdateLastCommitInfo(attr.SourceBranch, entities.Message{
		Text: fmt.Sprintf("'%s' is merged into '%s'", attr.SourceBranch, attr.TargetBranch),
		URL:  attr.URL,
	})

	if err == nil {
		c.JSON(http.StatusOK, api.Response{
			Ok:     true,
			Result: "updated",
		})
		return
	}

	c.JSON(api.GetErrStatusCode(err), api.Response{
		Ok:     true,
		Result: "task mentioned in the description",
	})
}
