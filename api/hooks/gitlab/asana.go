package gitlab

import (
	"bitbucket.org/mikehouston/asana-go"
	"errors"
	"fmt"
	"github.com/fadyat/hooks/api"
	"github.com/fadyat/hooks/api/entities"
	"github.com/fadyat/hooks/api/helpers"
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

// MergeRequestAsana handles gitlab merge request hook
func MergeRequestAsana(c *gin.Context) {
	icfg, exists := c.Get("HTTPAPI")
	if !exists {
		endWithError(c, errors.New("apiConfig not found"), http.StatusInternalServerError, &log.Logger)
		return
	}

	cfg := icfg.(*api.HTTPAPI)
	var gitlabRequest entities.GitlabMergeRequestHook
	if err := c.BindJSON(&gitlabRequest); err != nil {
		endWithError(c, err, http.StatusBadRequest, &log.Logger)
		return
	}

	logger := log.Logger.With().Str("pr", gitlabRequest.ObjectAttributes.URL).Logger()
	if !slices.Contains(cfg.GitlabSecretTokens, c.GetHeader("X-Gitlab-Token")) {
		endWithError(c, errors.New("invalid gitlab token"), http.StatusUnauthorized, &log.Logger)
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
		p := &asana.Project{ID: asanaURL.ProjectID}

		err := p.Fetch(client)
		if err != nil {
			e := err.(*asana.Error)
			logger.Info().Msg(fmt.Sprintf("Failed to fetch asana project %s, %s", asanaURL.ProjectID, e.Message))
			continue
		}

		t := &asana.Task{ID: asanaURL.TaskID}

		lastCommitField, asanaErr := helpers.GetCustomField(p, cfg.LastCommitFieldName)

		if asanaErr != nil {
			logger.Info().Msg(fmt.Sprintf("Failed to get custom field %s, %s", cfg.LastCommitFieldName, asanaErr.Message))
			comment := fmt.Sprintf("%s\n\n %s", lastCommit.URL, lastCommit.Message)
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
			comment := fmt.Sprintf("%s\n\n %s", lastCommit.URL, lastCommit.Message)
			helpers.CreateTaskCommentWithLogs(t, client, &comment, &logger)
			continue
		}

		logger.Debug().Msg(fmt.Sprintf("Updated asana task %s", asanaURL.TaskID))
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
