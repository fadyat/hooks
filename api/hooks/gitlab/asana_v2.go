package gitlab

import (
	"bitbucket.org/mikehouston/asana-go"
	"errors"
	"github.com/fadyat/hooks/api"
	"github.com/fadyat/hooks/api/entities"
	"github.com/fadyat/hooks/api/helpers"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
	"net/http"
	"sort"
)

// PushRequestAsanaV2 godoc
// @Summary     Gitlab push request hook
// @Description Endpoint to set last commit url to custom field in asana task, passed via branch name
// @Tags        gitlab
// @Accept      json
// @Produce     json
// @Param       X-Gitlab-Token header   string                         true "Gitlab token"
// @Param       body           body     entities.GitlabPushRequestHook true "Gitlab push request"
// @Success     200            {object} gitlab.SuccessResponse
// @Failure     400            {object} gitlab.ErrorResponse
// @Failure     401            {object} gitlab.ErrorResponse
// @Failure     500            {object} gitlab.ErrorResponse
// @Router      /api/v2/asana/push [post]
func PushRequestAsanaV2(c *gin.Context) {
	icfg, exists := c.Get("HTTPAPI")
	if !exists {
		helpers.EndWithError(c, errors.New("apiConfig not found"), http.StatusInternalServerError, &log.Logger)
		return
	}

	cfg := icfg.(*api.HTTPAPI)
	if !slices.Contains(cfg.GitlabSecretTokens, c.GetHeader("X-Gitlab-Token")) {
		helpers.EndWithError(c, errors.New("invalid gitlab token"), http.StatusUnauthorized, &log.Logger)
		return
	}

	var hook entities.GitlabPushRequestHook
	if err := c.ShouldBindJSON(&hook); err != nil {
		helpers.EndWithError(c, err, http.StatusBadRequest, &log.Logger)
		return
	}

	if hook.ObjectKind != entities.GitlabRequestObjectNamePush {
		helpers.EndWithError(c, errors.New("invalid event type"), http.StatusBadRequest, &log.Logger)
		return
	}

	if len(hook.Commits) == 0 {
		helpers.EndWithError(c, errors.New("no commits found"), http.StatusBadRequest, &log.Logger)
		return
	}

	source := helpers.GetBranchNameFromRef(hook.Ref)
	tID := helpers.GetAsanaTaskID(source)
	if tID == "" {
		helpers.EndWithError(c, errors.New("brName doesn't contain asana task_id"), http.StatusBadRequest, &log.Logger)
		return
	}

	client := asana.NewClientWithAccessToken(cfg.AsanaAPIKey)
	sort.SliceStable(hook.Commits, func(i, j int) bool {
		return hook.Commits[i].Timestamp.Before(hook.Commits[j].Timestamp)
	})
	lastCommit := hook.Commits[len(hook.Commits)-1]
	taskStatus := entities.NewTaskStatus(cfg.LastCommitFieldName, lastCommit.Message, lastCommit.URL)

	err := helpers.UpdateAsanaTaskStatus(client, tID, taskStatus)
	if err != nil {
		helpers.EndWithError(c, err, http.StatusInternalServerError, &log.Logger)
		return
	}

	log.Info().Msgf("task %s updated", tID)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
