package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upSPPPaymentType, downSPPPaymentType)
}

func upSPPPaymentType(ctx context.Context, tx *sql.Tx) error {
	// Add payment_type to differentiate SPP (sekolah) vs Syahriah (pesantren)
	alterations := []string{
		// payment_type: spp (sekolah), syahriah (pesantren), other
		`ALTER TABLE spp_transactions ADD COLUMN IF NOT EXISTS payment_type VARCHAR(20) DEFAULT 'spp'`,

		// Reference to unified students table (optional, for linking to student)
		`ALTER TABLE spp_transactions ALTER COLUMN student_id TYPE UUID USING student_id::uuid`,

		// Index for payment type queries
		`CREATE INDEX IF NOT EXISTS idx_spp_payment_type ON spp_transactions(tenant_id, payment_type)`,

		// Index for student queries
		`CREATE INDEX IF NOT EXISTS idx_spp_student ON spp_transactions(student_id) WHERE student_id IS NOT NULL`,
	}

	for _, query := range alterations {
		_, err := tx.ExecContext(ctx, query)
		if err != nil {
			// Log but continue - some alterations may fail if column exists
			continue
		}
	}

	return nil
}

func downSPPPaymentType(ctx context.Context, tx *sql.Tx) error {
	dropStatements := []string{
		`DROP INDEX IF EXISTS idx_spp_payment_type`,
		`DROP INDEX IF EXISTS idx_spp_student`,
		`ALTER TABLE spp_transactions DROP COLUMN IF EXISTS payment_type`,
	}

	for _, query := range dropStatements {
		_, _ = tx.ExecContext(ctx, query)
	}

	return nil
}
