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
	"sort"
)

var (
	cachedLastCommits = make(map[string]string, 100000)
	allPrAsanaTasks   = make(map[string][]string, 100000)
)

func clearCache() {
	cacheLen := 100000
	clearLen := cacheLen * 9 / 10
	optimalLen := cacheLen * 5 / 10

	if len(cachedLastCommits) > clearLen {
		for k := range cachedLastCommits {
			delete(cachedLastCommits, k)
			if len(cachedLastCommits) < optimalLen {
				break
			}
		}
	}
}

func clearPrAsanaTasks(prURL string) {
	if _, ok := allPrAsanaTasks[prURL]; ok {
		delete(allPrAsanaTasks, prURL)
	}
}

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

	prLink := gitlabRequest.ObjectAttributes.URL
	logger := log.Logger.With().Str("pr", prLink).Logger()
	if !slices.Contains(cfg.GitlabSecretTokens, c.GetHeader("X-Gitlab-Token")) {
		helpers.EndWithError(c, errors.New("invalid gitlab token"), http.StatusUnauthorized, &logger)
		return
	}

	if gitlabRequest.ObjectKind != entities.GitlabRequestObjectNameMergeRequest {
		helpers.EndWithError(c, errors.New("invalid event type"), http.StatusBadRequest, &logger)
		return
	}

	action := gitlabRequest.ObjectAttributes.Action
	if !slices.Contains(cfg.SupportedMergeRequestActions, action) {
		helpers.EndWithError(c, errors.New("unsupported action: "+action), http.StatusBadRequest, &logger)
		return
	}

	lastCommit := gitlabRequest.ObjectAttributes.LastCommit
	urls := helpers.GetAsanaURLS(lastCommit.Message)
	if len(urls) == 0 {
		logger.Info().Msg("No asana URLS found")
	}

	client := asana.NewClientWithAccessToken(cfg.AsanaAPIKey)
	if action == entities.OpenAction || action == entities.UpdateAction {
		for _, asanaURL := range urls {
			cachedLastCommitID, have := cachedLastCommits[asanaURL.TaskID]
			if have && cachedLastCommitID == lastCommit.ID {
				logger.Debug().Msg(fmt.Sprintf("last commit %s already cached", lastCommit.URL))
				continue
			}

			cachedLastCommits[asanaURL.TaskID] = lastCommit.ID
			helpers.UpdateAsanaTaskLastCommitInfo(
				client,
				&asanaURL,
				lastCommit.Message,
				lastCommit.URL,
				cfg.LastCommitFieldName,
				&logger,
			)

			if _, ok := allPrAsanaTasks[prLink]; !ok {
				allPrAsanaTasks[prLink] = []string{}
			}

			if !slices.Contains(allPrAsanaTasks[prLink], asanaURL.TaskID) {
				allPrAsanaTasks[prLink] = append(allPrAsanaTasks[prLink], asanaURL.TaskID)
			}
		}
	}

	if action == entities.MergeAction {
		for _, taskID := range allPrAsanaTasks[prLink] {
			t := &asana.Task{ID: taskID}
			err := t.Fetch(client)
			if err != nil {
				e := err.(*asana.Error)
				logger.Info().Msg(fmt.Sprintf("Failed to fetch asana task %s, %s", taskID, e.Message))
				return
			}

			mergeCommit := fmt.Sprintf("%s/commits/%s", gitlabRequest.Project.WebURL, gitlabRequest.ObjectAttributes.MergeCommitSha)
			mergeMsg := fmt.Sprintf(
				"%s\n\nMerge branch '%s' into '%s'",
				mergeCommit,
				gitlabRequest.ObjectAttributes.SourceBranch,
				gitlabRequest.ObjectAttributes.TargetBranch,
			)
			helpers.CreateTaskCommentWithLogs(t, client, &mergeMsg, &logger)
		}

		clearPrAsanaTasks(prLink)
	}

	clearCache()
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// PushRequestAsana godoc
// @Summary     Gitlab push request hook
// @Description Endpoint to set last commit url to custom field in asana task, passed via commit message
// @Tags        gitlab
// @Accept      json
// @Produce     json
// @Param       X-Gitlab-Token header   string                          true "Gitlab token"
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

	if !slices.Contains(cfg.GitlabSecretTokens, c.GetHeader("X-Gitlab-Token")) {
		helpers.EndWithError(c, errors.New("invalid gitlab token"), http.StatusUnauthorized, &log.Logger)
		return
	}

	if gitlabRequest.ObjectKind != entities.GitlabRequestObjectNamePush {
		helpers.EndWithError(c, errors.New("invalid event type"), http.StatusBadRequest, &log.Logger)
		return
	}

	if len(gitlabRequest.Commits) == 0 {
		helpers.EndWithError(c, errors.New("no commits found"), http.StatusBadRequest, &log.Logger)
		return
	}

	sort.SliceStable(gitlabRequest.Commits, func(i, j int) bool {
		return gitlabRequest.Commits[i].Timestamp.Before(gitlabRequest.Commits[j].Timestamp)
	})

	client := asana.NewClientWithAccessToken(cfg.AsanaAPIKey)
	lastCommit := gitlabRequest.Commits[len(gitlabRequest.Commits)-1]

	urls := helpers.GetAsanaURLS(lastCommit.Message)

	if len(urls) == 0 {
		log.Logger.Info().Msg("No asana URLS found")
	}

	for _, u := range urls {
		cachedLastCommitID, have := cachedLastCommits[u.TaskID]
		if have && cachedLastCommitID == lastCommit.ID {
			log.Logger.Debug().Msg(fmt.Sprintf("last commit %s already cached", lastCommit.URL))
			continue
		}

		cachedLastCommits[u.TaskID] = lastCommit.ID
		helpers.UpdateAsanaTaskLastCommitInfo(
			client,
			&u,
			lastCommit.Message,
			lastCommit.URL,
			cfg.LastCommitFieldName,
			&log.Logger,
		)
	}

	clearCache()
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
