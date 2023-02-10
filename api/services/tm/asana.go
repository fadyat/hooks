package tm

import (
	"bitbucket.org/mikehouston/asana-go"
	"errors"
	"github.com/fadyat/hooks/api"
	"github.com/fadyat/hooks/api/config"
	"github.com/fadyat/hooks/api/entities"
	"github.com/fadyat/hooks/api/entities/gitlab"
	"github.com/fadyat/hooks/api/helpers"
	"github.com/rs/zerolog"
)

type AsanaService struct {
	c   *asana.Client
	l   *zerolog.Logger
	cfg *config.HTTPAPI
}

// NewAsanaService creates a new instance of the Asana service
func NewAsanaService(apiKey string, l *zerolog.Logger, cfg *config.HTTPAPI) *AsanaService {
	return &AsanaService{
		l:   l,
		cfg: cfg,
		c:   asana.NewClientWithAccessToken(apiKey),
	}
}

func (a *AsanaService) UpdateCustomField(mention entities.TaskMention, customFieldName string, value string) error {
	t := asana.Task{ID: mention.ID}
	if err := t.Fetch(a.c); err != nil {
		return err
	}

	cf := helpers.FindCustomFieldByName(t.CustomFields, customFieldName)
	if cf == nil {
		return errors.New(api.CustomFieldNotFound)
	}

	err := t.Update(a.c, &asana.UpdateTaskRequest{
		CustomFields: map[string]interface{}{cf.ID: value},
	})

	if err == nil {
		a.l.Debug().Msgf("custom field %s updated for task %s", customFieldName, mention.ID)
	}

	return err
}

func (a *AsanaService) CreateComment(mention entities.TaskMention, value string) error {
	t := asana.Task{ID: mention.ID}
	_, err := t.CreateComment(a.c, &asana.StoryBase{
		Text: value,
	})

	if err == nil {
		a.l.Debug().Msgf("comment created for task %s", mention.ID)
	}

	return err
}

func (a *AsanaService) UpdateLastCommitInfo(branchName string, lastCommit gitlab.Commit) error {
	message := helpers.ConfigureMessageForTaskManager(
		lastCommit.Message,
		lastCommit.URL,
	)

	mentions := append(
		helpers.ParseTaskMentions(branchName),
		helpers.ParseTaskMentions(lastCommit.Message)...,
	)

	if len(mentions) == 0 {
		a.l.Debug().Msgf("no task mentions found in branch name %s or commit message %s", branchName, lastCommit.Message)
		return errors.New(api.NoTaskMentionsFound)
	}

	var wrappedError error = nil
	for _, m := range mentions {
		err := a.UpdateCustomField(m, a.cfg.LastCommitFieldName, lastCommit.URL)
		if err == nil {
			continue
		}

		a.l.Debug().Err(err).Msgf("failed to update custom field %s for task %s", a.cfg.LastCommitFieldName, m)
		if err = a.CreateComment(m, message); err != nil {
			a.l.Debug().Err(err).Msgf("failed to create comment for task %s", m)
			wrappedError = helpers.WrapError(wrappedError, err)
		}
	}

	return wrappedError
}
