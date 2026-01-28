package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upDashboardTables, downDashboardTables)
}

func upDashboardTables(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		-- SPP Transactions (payments from parents to school)
		CREATE TABLE IF NOT EXISTS spp_transactions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL,
			student_id UUID,
			student_name VARCHAR(255) NOT NULL,
			amount BIGINT NOT NULL,
			payment_method VARCHAR(50),
			status VARCHAR(20) NOT NULL DEFAULT 'pending',
			gateway_ref VARCHAR(100),
			description TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		);

		-- Disbursement requests (withdrawal from tenants)
		CREATE TABLE IF NOT EXISTS disbursements (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL,
			amount BIGINT NOT NULL,
			bank_name VARCHAR(100) NOT NULL,
			account_number VARCHAR(50) NOT NULL,
			account_holder VARCHAR(255) NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'pending',
			notes TEXT,
			admin_notes TEXT,
			requested_at TIMESTAMP NOT NULL DEFAULT NOW(),
			processed_at TIMESTAMP
		);

		-- Notification logs
		CREATE TABLE IF NOT EXISTS notification_logs (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			type VARCHAR(50) NOT NULL,
			recipient VARCHAR(255) NOT NULL,
			subject VARCHAR(255),
			message TEXT NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'sent',
			error_message TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		);

		-- Create indexes
		CREATE INDEX IF NOT EXISTS idx_spp_transactions_tenant ON spp_transactions(tenant_id);
		CREATE INDEX IF NOT EXISTS idx_spp_transactions_status ON spp_transactions(status);
		CREATE INDEX IF NOT EXISTS idx_disbursements_tenant ON disbursements(tenant_id);
		CREATE INDEX IF NOT EXISTS idx_disbursements_status ON disbursements(status);
		CREATE INDEX IF NOT EXISTS idx_notification_logs_type ON notification_logs(type);
	`)
	return err
}

func downDashboardTables(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		DROP TABLE IF EXISTS notification_logs;
		DROP TABLE IF EXISTS disbursements;
		DROP TABLE IF EXISTS spp_transactions;
	`)
	return err
}
