package postgres_outbound_adapter

import (
	"database/sql"
	"time"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"

	"github.com/palantir/stacktrace"
)

type contentAdapter struct {
	db outbound_port.DatabaseExecutor
}

func NewContentAdapter(db outbound_port.DatabaseExecutor) outbound_port.ContentDatabasePort {
	return &contentAdapter{
		db: db,
	}
}

func (a *contentAdapter) Upsert(content *model.Content) error {
	query := `
		INSERT INTO contents (key, value, type, updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (key) DO UPDATE
		SET value = EXCLUDED.value,
			type = EXCLUDED.type,
			updated_at = EXCLUDED.updated_at
	`
	_, err := a.db.Exec(query, content.Key, content.Value, content.Type, time.Now())
	if err != nil {
		return stacktrace.Propagate(err, "failed to upsert content")
	}
	return nil
}

func (a *contentAdapter) FindByKey(key string) (*model.Content, error) {
	query := `SELECT key, value, type, updated_at FROM contents WHERE key = $1`
	row := a.db.QueryRow(query, key)

	var content model.Content
	err := row.Scan(&content.Key, &content.Value, &content.Type, &content.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, stacktrace.Propagate(err, "content not found")
		}
		return nil, stacktrace.Propagate(err, "failed to scan content")
	}

	return &content, nil
}

func (a *contentAdapter) Delete(key string) error {
	query := `DELETE FROM contents WHERE key = $1`
	_, err := a.db.Exec(query, key)
	if err != nil {
		return stacktrace.Propagate(err, "failed to delete content")
	}
	return nil
}
