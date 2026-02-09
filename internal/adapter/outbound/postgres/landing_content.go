package postgres_outbound_adapter

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	outbound_port "prabogo/internal/port/outbound"
)

type landingContentAdapter struct {
	db outbound_port.DatabaseExecutor
}

func NewLandingContentAdapter(db outbound_port.DatabaseExecutor) outbound_port.LandingContentDatabasePort {
	return &landingContentAdapter{db: db}
}

func (a *landingContentAdapter) Get(ctx context.Context, key string) (map[string]interface{}, error) {
	query := `SELECT value FROM landing_content WHERE key = $1`

	var valueJSON []byte
	err := a.db.QueryRow(query, key).Scan(&valueJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil if key not found
		}
		return nil, fmt.Errorf("failed to get landing content: %w", err)
	}

	var raw interface{}
	if err := json.Unmarshal(valueJSON, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse landing content JSON: %w", err)
	}

	// If it's a map, return as is
	if m, ok := raw.(map[string]interface{}); ok {
		return m, nil
	}

	// If it's anything else (array, primitive), wrap it
	return map[string]interface{}{"value": raw}, nil
}

func (a *landingContentAdapter) Set(ctx context.Context, key string, value map[string]interface{}) error {
	// If the value contains a "value" wrapper key and we are storing an array, extract it.
	// However, the port interface takes map[string]interface{}.
	// Let's just serialize whatever is passed.

	valueJSON, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal landing content: %w", err)
	}

	query := `
		INSERT INTO landing_content (key, value, updated_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		ON CONFLICT (key) DO UPDATE SET
			value = EXCLUDED.value,
			updated_at = EXCLUDED.updated_at
	`

	_, err = a.db.Exec(query, key, valueJSON)
	if err != nil {
		return fmt.Errorf("failed to set landing content: %w", err)
	}

	return nil
}
