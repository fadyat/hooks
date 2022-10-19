package gitlab

import (
	"bitbucket.org/mikehouston/asana-go"
	"errors"
	"fmt"
	"github.com/fadyat/hooks/api"
	"github.com/fadyat/hooks/api/entities"
	"github.com/fadyat/hooks/api/helpers"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
	"net/http"
	"strings"
)

// MergeRequestAsana godoc
// @Summary     Gitlab merge request hook
// @Description Endpoint to set last commit url to custom field in asana task, passed via commit message
// @Tags        gitlab
// @Accept      json
// @Produce     json
// @Param       X-Gitlab-Token header   string                          true "Gitlab token"
// @Param       body           body     entities.GitlabMergeRequestHook true "Gitlab merge request"
// @Success     200            {object} gitlab.SuccessResponse
// @Failure     400            {object} gitlab.ErrorResponse
// @Failure     401            {object} gitlab.ErrorResponse
// @Failure     500            {object} gitlab.ErrorResponse
// @Router      /api/v1/asana/merge [post]
func MergeRequestAsana(c *gin.Context) {
	icfg, exists := c.Get("HTTPAPI")
	if !exists {
		helpers.EndWithError(c, errors.New("apiConfig not found"), http.StatusInternalServerError, &log.Logger)
		return
	}

	cfg := icfg.(*api.HTTPAPI)
	var gitlabRequest entities.GitlabMergeRequestHook
	if err := c.BindJSON(&gitlabRequest); err != nil {
		helpers.EndWithError(c, err, http.StatusBadRequest, &log.Logger)
		return
	}

	logger := log.Logger.With().Str("pr", gitlabRequest.ObjectAttributes.URL).Logger()
	if !slices.Contains(cfg.GitlabSecretTokens, c.GetHeader("X-Gitlab-Token")) {
		helpers.EndWithError(c, errors.New("invalid gitlab token"), http.StatusUnauthorized, &logger)
		return
	}

	if gitlabRequest.ObjectKind != entities.GitlabRequestObjectNameMergeRequest {
		helpers.EndWithError(c, errors.New("invalid event type"), http.StatusBadRequest, &logger)
		return
	}

	const cutset string = "\f\t\r\n "
	lastCommit := gitlabRequest.ObjectAttributes.LastCommit
	lastCommitURL := strings.Trim(lastCommit.URL, cutset)

	urls := helpers.GetAsanaURLS(lastCommit.Message)
	if len(urls) == 0 {
		logger.Info().Msg("No asana URLS found")
	}

	client := asana.NewClientWithAccessToken(cfg.AsanaAPIKey)
	for _, asanaURL := range urls {
		// todo: take project from task fetching
		p := &asana.Project{ID: asanaURL.ProjectID}

		err := p.Fetch(client)
		if err != nil {
			e := err.(*asana.Error)
			logger.Info().Msg(fmt.Sprintf("Failed to fetch asana project %s, %s", asanaURL.ProjectID, e.Message))
			continue
		}

		t := &asana.Task{ID: asanaURL.TaskID}

		lastCommitField, asanaErr := helpers.GetCustomField(p, cfg.LastCommitFieldName)
		filteredMessage := helpers.RemoveAsanaURLS(lastCommit.Message)

		if asanaErr != nil {
			logger.Info().Msg(fmt.Sprintf("Failed to get custom field %s, %s", cfg.LastCommitFieldName, asanaErr.Message))
			comment := fmt.Sprintf("%s\n\n %s", lastCommit.URL, filteredMessage)
			helpers.CreateTaskCommentWithLogs(t, client, &comment, &logger)
			continue
		}

		err = t.Update(client, &asana.UpdateTaskRequest{
			CustomFields: map[string]interface{}{
				lastCommitField.ID: lastCommitURL,
			},
		})

		if err != nil {
			e := err.(*asana.Error)
			logger.Info().Msg(fmt.Sprintf("Failed to update asana task %s, %s", asanaURL.TaskID, e.Message))
			comment := fmt.Sprintf("%s\n\n %s", lastCommit.URL, filteredMessage)
			helpers.CreateTaskCommentWithLogs(t, client, &comment, &logger)
			continue
		}

		logger.Debug().Msg(fmt.Sprintf("Updated asana task %s", asanaURL.TaskID))
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// PushRequestAsana godoc
// @Summary     Gitlab push request hook
// @Description Endpoint to set last commit url to custom field in asana task, passed via commit message
// @Tags        gitlab
// @Accept      json
// @Produce     json
// @Param       X-Gitlab-Token header   string                         true "Gitlab token"
// @Param       body           body     entities.GitlabPushRequestHook true "Gitlab push request"
// @Success     200            {object} gitlab.SuccessResponse
// @Failure     400            {object} gitlab.ErrorResponse
// @Failure     401            {object} gitlab.ErrorResponse
// @Failure     500            {object} gitlab.ErrorResponse
// @Router      /api/v1/asana/push [post]
func PushRequestAsana(c *gin.Context) {
	icfg, exists := c.Get("HTTPAPI")
	if !exists {
		helpers.EndWithError(c, errors.New("apiConfig not found"), http.StatusInternalServerError, &log.Logger)
		return
	}
	cfg := icfg.(*api.HTTPAPI)
	var gitlabRequest entities.GitlabPushRequestHook
	if err := c.BindJSON(&gitlabRequest); err != nil {
		helpers.EndWithError(c, err, http.StatusBadRequest, &log.Logger)
		return
	}

	logger := log.Logger.With().Str("pr", gitlabRequest.Project.WebURL).Logger()
	if !slices.Contains(cfg.GitlabSecretTokens, c.GetHeader("X-Gitlab-Token")) {
		helpers.EndWithError(c, errors.New("invalid gitlab token"), http.StatusUnauthorized, &logger)
		return
	}

	if gitlabRequest.ObjectKind != entities.GitlabRequestObjectNamePushRequest {
		helpers.EndWithError(c, errors.New("invalid event type"), http.StatusBadRequest, &logger)
		return
	}

	// todo: is it really last commit?
	lastCommit := gitlabRequest.Commits[len(gitlabRequest.Commits)-1]

	urls := helpers.GetAsanaURLS(lastCommit.Message)
	if len(urls) == 0 {
		logger.Info().Msg("No asana URLS found")
	}

	client := asana.NewClientWithAccessToken(cfg.AsanaAPIKey)
	for _, asanaURL := range urls {
		// todo: take project from task fetching
		p := &asana.Project{ID: asanaURL.ProjectID}

		err := p.Fetch(client)
		if err != nil {
			e := err.(*asana.Error)
			logger.Info().Msg(fmt.Sprintf("Failed to fetch asana project %s, %s", asanaURL.ProjectID, e.Message))
			continue
		}

		t := &asana.Task{ID: asanaURL.TaskID}
		lastCommitField, asanaErr := helpers.GetCustomField(p, cfg.LastCommitFieldName)
		filteredMessage := helpers.RemoveAsanaURLS(lastCommit.Message)

		if asanaErr != nil {
			logger.Info().Msg(fmt.Sprintf("Failed to get custom field %s, %s", cfg.LastCommitFieldName, asanaErr.Message))
			comment := fmt.Sprintf("%s\n\n %s", lastCommit.URL, filteredMessage)
			helpers.CreateTaskCommentWithLogs(t, client, &comment, &logger)
			continue
		}

		err = t.Update(client, &asana.UpdateTaskRequest{
			CustomFields: map[string]interface{}{
				lastCommitField.ID: lastCommit.URL,
			},
		})

		if err != nil {
			e := err.(*asana.Error)
			logger.Info().Msg(fmt.Sprintf("Failed to update asana task %s, %s", asanaURL.TaskID, e.Message))
			comment := fmt.Sprintf("%s\n\n %s", lastCommit.URL, filteredMessage)
			helpers.CreateTaskCommentWithLogs(t, client, &comment, &logger)
			continue
		}
	}
}
