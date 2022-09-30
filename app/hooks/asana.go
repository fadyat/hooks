package hooks

import (
	"bitbucket.org/mikehouston/asana-go"
	"github.com/fadyat/gitlab-hooks/app"
	"github.com/fadyat/gitlab-hooks/app/entities"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

// GitlabMergeRequestAsana hooking merge request events.
// Sending the merge request URL to the asana task specified in the comment.
func GitlabMergeRequestAsana(c *gin.Context) {
	cfg, exists := c.MustGet("AsanaConfig").(app.ApiConfig)
	if !exists {
		c.JSON(500, gin.H{"error": "Cannot get config for asana"})
		return
	}

	var gitlabRequest entities.GitlabMergeRequestHook
	if err := c.BindJSON(&gitlabRequest); err != nil {
		c.JSON(400, gin.H{"error": "Cannot parse gitlab request"})
		return
	}

	if !slices.Contains(cfg.GitlabSecretTokens, c.GetHeader("X-Gitlab-Token")) {
		c.JSON(403, gin.H{"error": "Invalid token"})
		return
	}

	client := asana.NewClientWithAccessToken(cfg.AsanaApiKey)
	urls := app.GetAsanaURLS(gitlabRequest.ObjectAttributes.Description)
	for _, asanaUrl := range *urls {
		t := &asana.Task{
			ID: asanaUrl.TaskId,
		}

		result, err := t.CreateComment(client, &asana.StoryBase{
			Text: gitlabRequest.ObjectAttributes.URL,
		})

		if err != nil {
			err := err.(asana.Error)
			c.JSON(err.StatusCode, gin.H{"message": err.Message})
			return
		}

		c.JSON(200, result)
	}
}
