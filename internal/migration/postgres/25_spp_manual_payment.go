package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationNoTxContext(upSPPManualPayment, downSPPManualPayment)
}

// upSPPManualPayment adds fields for manual payment confirmation workflow
func upSPPManualPayment(ctx context.Context, db *sql.DB) error {
	queries := []struct {
		name  string
		query string
	}{
		{
			name:  "add payment_proof to spp_transactions",
			query: `ALTER TABLE spp_transactions ADD COLUMN IF NOT EXISTS payment_proof VARCHAR(500)`,
		},
		{
			name:  "add confirmed_by to spp_transactions",
			query: `ALTER TABLE spp_transactions ADD COLUMN IF NOT EXISTS confirmed_by VARCHAR(50)`,
		},
		{
			name:  "add paid_at to spp_transactions",
			query: `ALTER TABLE spp_transactions ADD COLUMN IF NOT EXISTS paid_at TIMESTAMP`,
		},
		{
			name:  "add due_date to spp_transactions",
			query: `ALTER TABLE spp_transactions ADD COLUMN IF NOT EXISTS due_date TIMESTAMP`,
		},
		{
			name:  "add period to spp_transactions",
			query: `ALTER TABLE spp_transactions ADD COLUMN IF NOT EXISTS period VARCHAR(20)`,
		},
		{
			name:  "create idx_spp_period_status",
			query: `CREATE INDEX IF NOT EXISTS idx_spp_period_status ON spp_transactions(tenant_id, period, status)`,
		},
		{
			name:  "create idx_spp_due_date",
			query: `CREATE INDEX IF NOT EXISTS idx_spp_due_date ON spp_transactions(due_date) WHERE status = 'pending'`,
		},
	}

	for _, q := range queries {
		if _, err := db.ExecContext(ctx, q.query); err != nil {
			fmt.Printf("Warning: %s: %v\n", q.name, err)
			// Continue on error for idempotent migrations
		}
	}

	return nil
}

func downSPPManualPayment(ctx context.Context, db *sql.DB) error {
	queries := []string{
		`DROP INDEX IF EXISTS idx_spp_due_date`,
		`DROP INDEX IF EXISTS idx_spp_period_status`,
		`ALTER TABLE spp_transactions DROP COLUMN IF EXISTS payment_proof`,
		`ALTER TABLE spp_transactions DROP COLUMN IF EXISTS confirmed_by`,
		`ALTER TABLE spp_transactions DROP COLUMN IF EXISTS paid_at`,
		`ALTER TABLE spp_transactions DROP COLUMN IF EXISTS due_date`,
		`ALTER TABLE spp_transactions DROP COLUMN IF EXISTS period`,
	}

	for _, q := range queries {
		if _, err := db.ExecContext(ctx, q); err != nil {
			fmt.Printf("Warning during rollback: %v\n", err)
		}
	}

	return nil
}
