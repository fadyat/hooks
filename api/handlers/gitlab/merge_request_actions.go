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

// OnMergeRequestActions handles the merge request actions and sends the comment to the task
// based on the action.
//
//	@Description	Creates a comment for the merge request actions.
//	@Accept			json
//	@Param			X-Gitlab-Event	header		string					true	"Gitlab event"
//	@Param			X-Gitlab-Token	header		string					true	"Gitlab token"
//	@Param			body			body		gitlab.PushRequestHook	true	"Gitlab request body"
//	@Success		200				{object}	api.Response
//	@Failure		400				{object}	api.Response
//	@Failure		401				{object}	api.Response
//	@Failure		500				{object}	api.Response
//	@Router			/api/v1/asana/merge [post]
func (h *Handler) OnMergeRequestActions(c *gin.Context) {
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
		h.l.Error().Err(err).Msg("invalid request body")
		c.JSON(http.StatusBadRequest, api.Response{
			Ok:    false,
			Error: "invalid request body",
		})
		return
	}

	text := getText(&r)
	if text == UnsupportedMergeRequestAction {
		h.l.Info().Msgf("unsupported action: %s", r.ObjectAttributes.Action)

		// returning 200 to avoid gitlab retrying the request
		c.JSON(http.StatusOK, api.Response{
			Ok:    false,
			Error: fmt.Sprintf("unsupported action: %s", r.ObjectAttributes.Action),
		})
		return
	}

	attr := r.ObjectAttributes
	if err := h.tm.CreateComment(entities.Message{
		Text:       text,
		URL:        attr.URL,
		Author:     r.User.Username,
		BranchName: attr.SourceBranch,
	}); err != nil {
		h.l.Error().Err(err).Msg("failed to update last commit info")
		c.JSON(api.GetErrStatusCode(err), api.Response{
			Ok:     true,
			Result: "task mentioned in the description",
		})
		return
	}

	c.JSON(http.StatusOK, api.Response{
		Ok:     true,
		Result: "updated",
	})
}

func getText(r *gitlab.MergeRequestHook) string {
	switch r.ObjectAttributes.Action {
	case gitlab.MergeRequestActionMerge:
		return fmt.Sprintf("%q is merged into %q", r.ObjectAttributes.SourceBranch, r.ObjectAttributes.TargetBranch)
	case gitlab.MergeRequestActionOpen:
		return fmt.Sprintf("Created merge request to merge %q into %q", r.ObjectAttributes.SourceBranch, r.ObjectAttributes.TargetBranch)
	case gitlab.MergeRequestActionReopen:
		return "The merge request has been reopened."
	case gitlab.MergeRequestActionClose:
		return "The merge request has been closed."
	case gitlab.MergeRequestActionApproved:
		return "The merge request has been approved."
	case gitlab.MergeRequestActionUnapproved:
		return "Changes have been requested in the merge request."
	}

	return UnsupportedMergeRequestAction
}

const (
	UnsupportedMergeRequestAction = "unsupported merge request action"
)
