package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upPayments, downPayments)
}

func upPayments(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS payments (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
			order_id VARCHAR(100) UNIQUE NOT NULL,
			amount BIGINT NOT NULL,
			status VARCHAR(20) DEFAULT 'pending',
			payment_type VARCHAR(50),
			snap_token VARCHAR(255),
			midtrans_id VARCHAR(100),
			paid_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_payments_tenant_id ON payments(tenant_id);
		CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id);
		CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);
	`)
	return err
}

func downPayments(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
		DROP INDEX IF EXISTS idx_payments_status;
		DROP INDEX IF EXISTS idx_payments_order_id;
		DROP INDEX IF EXISTS idx_payments_tenant_id;
		DROP TABLE IF EXISTS payments;
	`)
	return err
}
