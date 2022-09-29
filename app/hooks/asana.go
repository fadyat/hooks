package hooks

import (
	"bitbucket.org/mikehouston/asana-go"
	"github.com/fadyat/gitlab-hooks/app"
	"github.com/fadyat/gitlab-hooks/app/entities"
	"github.com/gin-gonic/gin"
)

// GitlabMergeRequestAsana hooking merge request events.
// Sending the merge request URL to the asana task specified in the comment.
func GitlabMergeRequestAsana(c *gin.Context) {
	// todo: validate the request X-Gitlab-Token

	cfg, exists := c.MustGet("AsanaConfig").(app.AsanaConfig)
	if !exists {
		c.JSON(500, gin.H{"error": "Cannot get config for asana"})
		return
	}

	var gitlabRequest entities.GitlabMergeRequestHook
	if err := c.BindJSON(&gitlabRequest); err != nil {
		c.JSON(400, gin.H{"error": "Cannot parse gitlab request"})
		return
	}

	client := asana.NewClientWithAccessToken(cfg.ApiKey)
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
