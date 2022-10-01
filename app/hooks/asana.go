package hooks

import (
	"bitbucket.org/mikehouston/asana-go"
	"errors"
	"fmt"
	"github.com/fadyat/gitlab-hooks/app"
	"github.com/fadyat/gitlab-hooks/app/entities"
	"github.com/fadyat/gitlab-hooks/app/helpers"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
	"net/http"
	"strings"
)

// endWithError ends request with error
func endWithError(c *gin.Context, err error, httpCode int, l *zerolog.Logger) {
	l.Info().Msg(err.Error())
	c.JSON(httpCode, gin.H{"error": err.Error()})
}

// itsIncorrectAsanaURL adds incorrect asana url to incorrectAsanaURLs
func itsIncorrectAsanaURL(
	incorrectAsanaURLs *[]entities.IncorrectAsanaURL,
	asanaURL entities.AsanaURL,
	err error,
) {
	*incorrectAsanaURLs = append(*incorrectAsanaURLs, entities.IncorrectAsanaURL{
		AsanaURL: asanaURL,
		Error:    err,
	})
}

// itsCorrectAsanaURL adds correct asana url to updatedAsanaTasks
func itsCorrectAsanaURL(
	updatedAsanaTasks *[]entities.UpdatedAsanaTask,
	asanaURL entities.AsanaURL,
) {
	*updatedAsanaTasks = append(*updatedAsanaTasks, entities.UpdatedAsanaTask{
		AsanaTaskID: asanaURL.TaskID,
	})
}

// GitlabMergeRequestAsana handles gitlab merge request hook
func GitlabMergeRequestAsana(c *gin.Context) {
	icfg, exists := c.Get("APIConfig")
	if !exists {
		endWithError(c, errors.New("apiConfig not found"), http.StatusInternalServerError, &log.Logger)
		return
	}

	cfg := icfg.(*app.APIConfig)
	var gitlabRequest entities.GitlabMergeRequestHook
	if err := c.BindJSON(&gitlabRequest); err != nil {
		endWithError(c, err, http.StatusBadRequest, &log.Logger)
		return
	}

	if !slices.Contains(cfg.GitlabSecretTokens, c.GetHeader("X-Gitlab-Token")) {
		endWithError(c, errors.New("invalid gitlab token"), http.StatusUnauthorized, &log.Logger)
		return
	}

	const cutset string = "\f\t\r\n "
	lastCommit := gitlabRequest.ObjectAttributes.LastCommit
	lastCommitURL := strings.Trim(lastCommit.URL, cutset)
	message := strings.Trim(lastCommit.Message, cutset)

	logger := log.Logger.With().
		Str("pr", gitlabRequest.ObjectAttributes.URL).
		Logger()

	urls := helpers.GetAsanaURLS(lastCommit.Message)
	if len(urls) == 0 {
		endWithError(c, errors.New("no asana urls found"), http.StatusBadRequest, &logger)
		return
	}

	client := asana.NewClientWithAccessToken(cfg.AsanaAPIKey)
	var incorrectAsanaURLs = make([]entities.IncorrectAsanaURL, 0)
	var updatedAsanaTasks = make([]entities.UpdatedAsanaTask, 0)
	for _, asanaURL := range urls {
		p := &asana.Project{
			ID: asanaURL.ProjectID,
		}

		err := p.Fetch(client)
		if err != nil {
			logger.Info().Msg(fmt.Sprintf("Failed to fetch asana project %s", asanaURL.ProjectID))
			itsIncorrectAsanaURL(&incorrectAsanaURLs, asanaURL, err)
			continue
		}

		lastCommitField, asanaErr := helpers.GetOrCreateCustomField(client, p, cfg.LastCommitFieldName)
		if asanaErr != nil {
			logger.Info().Msg(fmt.Sprintf("Failed to get or create '%s' custom field, %s", cfg.LastCommitFieldName, p.Name))
			itsIncorrectAsanaURL(&incorrectAsanaURLs, asanaURL, asanaErr)
			continue
		}
		messageField, asanaErr := helpers.GetOrCreateCustomField(client, p, cfg.MessageCommitFieldName)
		if asanaErr != nil {
			logger.Info().Msg(fmt.Sprintf("Failed to get or create '%s' custom field, %s", cfg.MessageCommitFieldName, p.Name))
			itsIncorrectAsanaURL(&incorrectAsanaURLs, asanaURL, asanaErr)
			continue
		}

		t := &asana.Task{
			ID: asanaURL.TaskID,
		}

		err = t.Update(client, &asana.UpdateTaskRequest{
			CustomFields: map[string]interface{}{
				lastCommitField.ID: lastCommitURL,
				messageField.ID:    message,
			},
		})

		if err != nil {
			logger.Info().Msg(fmt.Sprintf("Failed to update asana task %s", asanaURL.TaskID))
			itsIncorrectAsanaURL(&incorrectAsanaURLs, asanaURL, asanaErr)
			continue
		}

		itsCorrectAsanaURL(&updatedAsanaTasks, asanaURL)
	}

	code := http.StatusOK
	if len(incorrectAsanaURLs) > 0 {
		code = http.StatusBadRequest
	}

	c.JSON(code, gin.H{"incorrect": incorrectAsanaURLs, "updated": updatedAsanaTasks})
}
