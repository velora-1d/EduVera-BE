package content

import (
	"context"

	"github.com/palantir/stacktrace"

	"eduvera/internal/model"
	outbound_port "eduvera/internal/port/outbound"
)

type ContentDomain interface {
	Upsert(ctx context.Context, input *model.ContentInput) (*model.Content, error)
	Get(ctx context.Context, key string) (*model.Content, error)
}

type contentDomain struct {
	databasePort outbound_port.DatabasePort
}

func NewContentDomain(databasePort outbound_port.DatabasePort) ContentDomain {
	return &contentDomain{
		databasePort: databasePort,
	}
}

func (d *contentDomain) Upsert(ctx context.Context, input *model.ContentInput) (*model.Content, error) {
	content := &model.Content{
		Key:   input.Key,
		Value: input.Value,
		Type:  input.Type,
	}

	err := d.databasePort.Content().Upsert(content)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to upsert content")
	}

	return content, nil
}

func (d *contentDomain) Get(ctx context.Context, key string) (*model.Content, error) {
	content, err := d.databasePort.Content().FindByKey(key)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to get content")
	}
	return content, nil
}
