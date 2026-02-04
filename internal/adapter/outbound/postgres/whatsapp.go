package postgres_outbound_adapter

import (
	"context"
	"database/sql"
	"time"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"

	"github.com/google/uuid"
)

type whatsAppAdapter struct {
	db outbound_port.DatabaseExecutor
}

func NewWhatsAppAdapter(db outbound_port.DatabaseExecutor) outbound_port.WhatsAppDatabasePort {
	return &whatsAppAdapter{db: db}
}

func (a *whatsAppAdapter) GetByTenantID(ctx context.Context, tenantID string) (*model.WhatsAppSession, error) {
	query := `
		SELECT id, tenant_id, instance_name, api_key, status, qr_code, device_info, created_at, updated_at
		FROM tenant_whatsapp_sessions
		WHERE tenant_id = $1
		LIMIT 1
	`

	row := a.db.QueryRow(query, tenantID)
	var session model.WhatsAppSession
	var qrCode sql.NullString
	var deviceInfo sql.NullString

	err := row.Scan(
		&session.ID,
		&session.TenantID,
		&session.InstanceName,
		&session.APIKey,
		&session.Status,
		&qrCode,
		&deviceInfo,
		&session.CreatedAt,
		&session.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	session.QRCode = qrCode.String
	session.DeviceInfo = deviceInfo.String
	return &session, nil
}

func (a *whatsAppAdapter) Save(ctx context.Context, session *model.WhatsAppSession) error {
	// Upsert - insert if not exists, update if exists
	query := `
		INSERT INTO tenant_whatsapp_sessions (id, tenant_id, instance_name, api_key, status, qr_code, device_info, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (tenant_id) DO UPDATE SET
			instance_name = EXCLUDED.instance_name,
			api_key = EXCLUDED.api_key,
			status = EXCLUDED.status,
			qr_code = EXCLUDED.qr_code,
			device_info = EXCLUDED.device_info,
			updated_at = EXCLUDED.updated_at
	`

	if session.ID == "" {
		session.ID = uuid.New().String()
	}
	now := time.Now()
	if session.CreatedAt.IsZero() {
		session.CreatedAt = now
	}
	session.UpdatedAt = now

	_, err := a.db.Exec(query,
		session.ID,
		session.TenantID,
		session.InstanceName,
		session.APIKey,
		session.Status,
		sql.NullString{String: session.QRCode, Valid: session.QRCode != ""},
		sql.NullString{String: session.DeviceInfo, Valid: session.DeviceInfo != ""},
		session.CreatedAt,
		session.UpdatedAt,
	)
	return err
}

func (a *whatsAppAdapter) UpdateStatus(ctx context.Context, sessionID, status string) error {
	query := `UPDATE tenant_whatsapp_sessions SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := a.db.Exec(query, status, time.Now(), sessionID)
	return err
}

func (a *whatsAppAdapter) Delete(ctx context.Context, sessionID string) error {
	query := `DELETE FROM tenant_whatsapp_sessions WHERE id = $1`
	_, err := a.db.Exec(query, sessionID)
	return err
}
