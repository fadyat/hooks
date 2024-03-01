package tm

import (
	"bitbucket.org/mikehouston/asana-go"
	"errors"
	"fmt"
	"github.com/fadyat/hooks/api"
	"github.com/fadyat/hooks/api/config"
	"github.com/fadyat/hooks/api/entities"
	"github.com/fadyat/hooks/api/helpers"
	"github.com/rs/zerolog"
)

type AsanaService struct {
	c             *asana.Client
	l             *zerolog.Logger
	cfg           *config.HTTPAPI
	featureFlags  *config.FeatureFlags
	messageParser *helpers.MessageParser
}

// NewAsanaService creates a new instance of the Asana service
func NewAsanaService(
	apiKey string,
	l *zerolog.Logger,
	cfg *config.HTTPAPI,
	ff *config.FeatureFlags,
) *AsanaService {
	return &AsanaService{
		l:             l,
		cfg:           cfg,
		c:             asana.NewClientWithAccessToken(apiKey),
		featureFlags:  ff,
		messageParser: helpers.NewMessageParser(ff),
	}
}

func (a *AsanaService) UpdateCustomField(mention *entities.TaskMention, customFieldName, value string) error {
	t := asana.Task{ID: mention.ID}
	if err := t.Fetch(a.c); err != nil {
		return err
	}

	cf := helpers.Find[asana.CustomFieldValue](t.CustomFields, func(cf *asana.CustomFieldValue) bool {
		return cf.Name == customFieldName
	})

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

func (a *AsanaService) createComment(mention *entities.TaskMention, value string) error {
	t := asana.Task{ID: mention.ID}
	_, err := t.CreateComment(a.c, &asana.StoryBase{
		Text: value,
	})

	if err == nil {
		a.l.Debug().Msgf("comment created for task %s", mention.ID)
	}

	return err
}

func (a *AsanaService) CreateComment(msg entities.Message) error {
	message, mentions, e := a.messageParser.GetTaskMentions(msg)
	if e != nil {
		return e
	}

	var wrappedError error
	for _, m := range mentions {
		if err := a.createComment(m, message); err != nil {
			a.l.Debug().Err(err).Msgf("failed to create comment for task %s", m)
			wrappedError = errors.Join(wrappedError, err)
		}
	}

	return wrappedError
}

func (a *AsanaService) UpdateLastCommitInfo(msg entities.Message) error {
	message, mentions, e := a.messageParser.GetTaskMentions(msg)
	if e != nil {
		return e
	}

	var wrappedError error
	for _, m := range mentions {
		err := a.UpdateCustomField(m, a.cfg.LastCommitFieldName, msg.URL)
		if err == nil {
			continue
		}

		a.l.Debug().Err(err).Msgf("failed to update custom field %s for task %s", a.cfg.LastCommitFieldName, m)
		if err = a.createComment(m, message); err != nil {
			a.l.Debug().Err(err).Msgf("failed to create comment for task %s", m)
			wrappedError = errors.Join(wrappedError, err)
		}
	}

	return wrappedError
}

func (a *AsanaService) GetTaskShortLink(mention *entities.TaskMention) (string, error) {
	return fmt.Sprintf("https://app.asana.com/0/0/%s/f", mention.ID), nil
}

func (a *AsanaService) GetTaskName(mention *entities.TaskMention) (string, error) {
	t := asana.Task{ID: mention.ID}
	if err := t.Fetch(a.c); err != nil {
		return "", err
	}

	return t.Name, nil
}
