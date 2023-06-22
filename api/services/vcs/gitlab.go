package vcs

import (
	"errors"
	"github.com/fadyat/hooks/api"
	"github.com/fadyat/hooks/api/entities"
	"github.com/fadyat/hooks/api/helpers"
	"github.com/fadyat/hooks/api/services/tm"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/xanzy/go-gitlab"
)

type GitlabService struct {
	c  *gitlab.Client
	l  *zerolog.Logger
	tm tm.ITaskManager
}

func NewGitlabService(apiKey string, l *zerolog.Logger, t tm.ITaskManager) *GitlabService {
	c, err := gitlab.NewClient(apiKey)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to create gitlab client")
	}

	return &GitlabService{
		l:  l,
		tm: t,
		c:  c,
	}
}

func (g *GitlabService) UpdatePRDescription(pid, prID int, branch, desc string) error {
	mentions := helpers.ParseTaskMentions(branch)
	if len(mentions) == 0 {
		return errors.New(api.NoTaskMentionsFound)
	}

	shortlinks := make([]entities.TaskMentionHidden, len(mentions))
	for i, mention := range mentions {
		sh, err := g.tm.GetTaskShortLink(mention)
		if err != nil {
			log.Debug().Err(err).Msgf("failed to get short link for task %s", mention.ID)
			continue
		}

		name, err := g.tm.GetTaskName(mention)
		if err != nil {
			log.Debug().Err(err).Msgf("failed to get name for task %s", mention.ID)
			name = mention.ID
		}

		shortlinks[i] = entities.TaskMentionHidden{
			ID:        mention.ID,
			Name:      name,
			ShortLink: sh,
		}
	}

	updatedDesc := entities.WrapInMarkDownLinks(shortlinks, "\n")
	updatedDesc += desc

	_, _, err := g.c.MergeRequests.UpdateMergeRequest(pid, prID, &gitlab.UpdateMergeRequestOptions{
		Description: &updatedDesc,
	})

	return err
}
