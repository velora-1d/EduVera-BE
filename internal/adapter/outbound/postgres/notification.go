package postgres_outbound_adapter

import (
	"context"
	"database/sql"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"
)

type notificationAdapter struct {
	db *sql.DB
}

func NewNotificationAdapter(db *sql.DB) outbound_port.NotificationDatabasePort {
	return &notificationAdapter{
		db: db,
	}
}

func (r *notificationAdapter) GetAll(ctx context.Context) ([]model.Notification, error) {
	query := `
		SELECT id, type, recipient, subject, message, status, error_message, created_at
		FROM notification_logs
		ORDER BY created_at DESC
		LIMIT 100
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []model.Notification
	for rows.Next() {
		var n model.Notification
		var subject, errorMsg sql.NullString
		if err := rows.Scan(
			&n.ID, &n.Type, &n.Recipient, &subject, &n.Message, &n.Status, &errorMsg, &n.CreatedAt,
		); err != nil {
			return nil, err
		}
		if subject.Valid {
			n.Subject = subject.String
		}
		if errorMsg.Valid {
			n.ErrorMsg = errorMsg.String
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

func (r *notificationAdapter) GetStats(ctx context.Context) (*model.NotificationStats, error) {
	query := `
		SELECT 
			COUNT(CASE WHEN status = 'sent' THEN 1 END) as total_sent,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as total_failed
		FROM notification_logs
	`

	var stats model.NotificationStats
	err := r.db.QueryRowContext(ctx, query).Scan(&stats.TotalSent, &stats.TotalFailed)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *notificationAdapter) Create(ctx context.Context, notification *model.Notification) error {
	query := `
		INSERT INTO notification_logs (type, recipient, subject, message, status, error_message, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query,
		notification.Type,
		notification.Recipient,
		notification.Subject,
		notification.Message,
		notification.Status,
		notification.ErrorMsg,
	).Scan(&notification.ID)

	return err
}
