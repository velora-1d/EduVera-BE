package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upWhatsAppSessions, downWhatsAppSessions)
}

func upWhatsAppSessions(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS tenant_whatsapp_sessions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			instance_name VARCHAR(100) NOT NULL,
			api_key VARCHAR(255),
			status VARCHAR(50) DEFAULT 'disconnected',
			qr_code TEXT,
			device_info TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);

		CREATE UNIQUE INDEX IF NOT EXISTS idx_whatsapp_sessions_tenant_id ON tenant_whatsapp_sessions(tenant_id);
		CREATE INDEX IF NOT EXISTS idx_whatsapp_sessions_status ON tenant_whatsapp_sessions(status);

		DROP TRIGGER IF EXISTS update_tenant_whatsapp_sessions_updated_at ON tenant_whatsapp_sessions;
		CREATE TRIGGER update_tenant_whatsapp_sessions_updated_at BEFORE UPDATE ON tenant_whatsapp_sessions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
	`)
	return err
}

func downWhatsAppSessions(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
		DROP TABLE IF EXISTS tenant_whatsapp_sessions;
	`)
	return err
}
