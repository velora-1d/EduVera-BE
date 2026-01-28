package outbound_port

import "eduvera/internal/model"

//go:generate mockgen -source=content.go -destination=./../../../tests/mocks/port/mock_content.go
type ContentDatabasePort interface {
	Upsert(content *model.Content) error
	FindByKey(key string) (*model.Content, error)
	Delete(key string) error
}
