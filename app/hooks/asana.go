package hooks

import (
	"bitbucket.org/mikehouston/asana-go"
	"errors"
	"github.com/fadyat/gitlab-hooks/app"
	"github.com/fadyat/gitlab-hooks/app/entities"
	"github.com/fadyat/gitlab-hooks/app/helpers"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"net/http"
	"strings"
)

const cutset string = "\f\t\r\n "

func endWithError(c *gin.Context, err error, httpCode int) {
	err = c.Error(err)
	c.JSON(httpCode, gin.H{"error": err.Error()})
}

func GitlabMergeRequestAsana(c *gin.Context) {
	icfg, exists := c.Get("ApiConfig")
	if !exists {
		endWithError(c, errors.New("apiConfig not found"), http.StatusInternalServerError)
		return
	}

	cfg := icfg.(app.ApiConfig)
	var gitlabRequest entities.GitlabMergeRequestHook
	if err := c.BindJSON(&gitlabRequest); err != nil {
		endWithError(c, err, http.StatusBadRequest)
		return
	}

	if !slices.Contains(cfg.GitlabSecretTokens, c.GetHeader("X-Gitlab-Token")) {
		endWithError(c, errors.New("invalid gitlab token"), http.StatusUnauthorized)
		return
	}

	lastCommit := gitlabRequest.ObjectAttributes.LastCommit
	lastCommitURL := strings.Trim(lastCommit.URL, cutset)
	message := strings.Trim(lastCommit.Message, cutset)

	urls := helpers.GetAsanaURLS(lastCommit.Message)
	if len(urls) == 0 {
		endWithError(c, errors.New("no asana urls found"), http.StatusBadRequest)
		return
	}

	client := asana.NewClientWithAccessToken(cfg.AsanaApiKey)
	for _, asanaUrl := range urls {
		p := &asana.Project{
			ID: asanaUrl.ProjectId,
		}

		err := p.Fetch(client)
		if err != nil {
			e := err.(*asana.Error)
			endWithError(c, e, e.StatusCode)
			return
		}

		lastCommitField, err := helpers.GetCustomField(p, cfg.LastCommitFieldName)
		messageField, err := helpers.GetCustomField(p, cfg.MessageCommitFieldName)
		if err != nil {
			// todo: create custom field instead of returning error
			e := err.(*asana.Error)
			endWithError(c, e, e.StatusCode)
			return
		}

		t := &asana.Task{
			ID: asanaUrl.TaskId,
		}

		err = t.Update(client, &asana.UpdateTaskRequest{
			CustomFields: map[string]interface{}{
				lastCommitField.ID: lastCommitURL,
				messageField.ID:    message,
			},
		})

		if err != nil {
			// todo: think about situation when cannot create some of the tasks
			endWithError(c, err, http.StatusInternalServerError)
			return
		}

		c.JSON(200, gin.H{"result": "ok"})
	}
}
