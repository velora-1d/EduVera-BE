package postgres_outbound_adapter

import (
	"context"
	"database/sql"
	"time"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"

	"github.com/google/uuid"
)

type notificationTemplateAdapter struct {
	db outbound_port.DatabaseExecutor
}

func NewNotificationTemplateAdapter(db outbound_port.DatabaseExecutor) outbound_port.NotificationTemplateDatabasePort {
	return &notificationTemplateAdapter{db: db}
}

func (a *notificationTemplateAdapter) GetAll(ctx context.Context) ([]model.NotificationTemplate, error) {
	query := `
		SELECT id, event_type, channel, template_name, template_content, variables, is_active, created_at, updated_at
		FROM notification_templates
		ORDER BY event_type, channel
	`

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []model.NotificationTemplate
	for rows.Next() {
		var t model.NotificationTemplate
		var variables sql.NullString
		err := rows.Scan(
			&t.ID,
			&t.EventType,
			&t.Channel,
			&t.TemplateName,
			&t.TemplateContent,
			&variables,
			&t.IsActive,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		t.Variables = variables.String
		templates = append(templates, t)
	}

	return templates, nil
}

func (a *notificationTemplateAdapter) GetByID(ctx context.Context, id string) (*model.NotificationTemplate, error) {
	query := `
		SELECT id, event_type, channel, template_name, template_content, variables, is_active, created_at, updated_at
		FROM notification_templates
		WHERE id = $1
	`

	var t model.NotificationTemplate
	var variables sql.NullString
	err := a.db.QueryRow(query, id).Scan(
		&t.ID,
		&t.EventType,
		&t.Channel,
		&t.TemplateName,
		&t.TemplateContent,
		&variables,
		&t.IsActive,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	t.Variables = variables.String
	return &t, nil
}

func (a *notificationTemplateAdapter) GetByEventAndChannel(ctx context.Context, eventType, channel string) (*model.NotificationTemplate, error) {
	query := `
		SELECT id, event_type, channel, template_name, template_content, variables, is_active, created_at, updated_at
		FROM notification_templates
		WHERE event_type = $1 AND channel = $2 AND is_active = true
	`

	var t model.NotificationTemplate
	var variables sql.NullString
	err := a.db.QueryRow(query, eventType, channel).Scan(
		&t.ID,
		&t.EventType,
		&t.Channel,
		&t.TemplateName,
		&t.TemplateContent,
		&variables,
		&t.IsActive,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	t.Variables = variables.String
	return &t, nil
}

func (a *notificationTemplateAdapter) GetActiveByEvent(ctx context.Context, eventType string) ([]model.NotificationTemplate, error) {
	query := `
		SELECT id, event_type, channel, template_name, template_content, variables, is_active, created_at, updated_at
		FROM notification_templates
		WHERE event_type = $1 AND is_active = true
		ORDER BY channel
	`

	rows, err := a.db.Query(query, eventType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []model.NotificationTemplate
	for rows.Next() {
		var t model.NotificationTemplate
		var variables sql.NullString
		err := rows.Scan(
			&t.ID,
			&t.EventType,
			&t.Channel,
			&t.TemplateName,
			&t.TemplateContent,
			&variables,
			&t.IsActive,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		t.Variables = variables.String
		templates = append(templates, t)
	}

	return templates, nil
}

func (a *notificationTemplateAdapter) Save(ctx context.Context, template *model.NotificationTemplate) error {
	if template.ID == "" {
		template.ID = uuid.New().String()
	}
	now := time.Now()
	template.CreatedAt = now
	template.UpdatedAt = now

	query := `
		INSERT INTO notification_templates (id, event_type, channel, template_name, template_content, variables, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := a.db.Exec(query,
		template.ID,
		template.EventType,
		template.Channel,
		template.TemplateName,
		template.TemplateContent,
		sql.NullString{String: template.Variables, Valid: template.Variables != ""},
		template.IsActive,
		template.CreatedAt,
		template.UpdatedAt,
	)
	return err
}

func (a *notificationTemplateAdapter) Update(ctx context.Context, template *model.NotificationTemplate) error {
	template.UpdatedAt = time.Now()

	query := `
		UPDATE notification_templates
		SET template_name = $1, template_content = $2, variables = $3, is_active = $4, updated_at = $5
		WHERE id = $6
	`

	_, err := a.db.Exec(query,
		template.TemplateName,
		template.TemplateContent,
		sql.NullString{String: template.Variables, Valid: template.Variables != ""},
		template.IsActive,
		template.UpdatedAt,
		template.ID,
	)
	return err
}

func (a *notificationTemplateAdapter) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM notification_templates WHERE id = $1`
	_, err := a.db.Exec(query, id)
	return err
}
