package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationNoTxContext(upHybridMultiName, downHybridMultiName)
}

// upHybridMultiName adds school_name and pesantren_name columns to tenants table
// for hybrid packages that need different names for school and pesantren
func upHybridMultiName(ctx context.Context, db *sql.DB) error {
	queries := []struct {
		name  string
		query string
	}{
		{
			name:  "add school_name to tenants",
			query: `ALTER TABLE tenants ADD COLUMN IF NOT EXISTS school_name VARCHAR(255)`,
		},
		{
			name:  "add pesantren_name to tenants",
			query: `ALTER TABLE tenants ADD COLUMN IF NOT EXISTS pesantren_name VARCHAR(255)`,
		},
		{
			name:  "add index on school_name",
			query: `CREATE INDEX IF NOT EXISTS idx_tenants_school_name ON tenants(school_name) WHERE school_name IS NOT NULL`,
		},
		{
			name:  "add index on pesantren_name",
			query: `CREATE INDEX IF NOT EXISTS idx_tenants_pesantren_name ON tenants(pesantren_name) WHERE pesantren_name IS NOT NULL`,
		},
	}

	for _, q := range queries {
		if _, err := db.ExecContext(ctx, q.query); err != nil {
			fmt.Printf("Warning: %s: %v\n", q.name, err)
		}
	}

	return nil
}

func downHybridMultiName(ctx context.Context, db *sql.DB) error {
	queries := []string{
		`DROP INDEX IF EXISTS idx_tenants_pesantren_name`,
		`DROP INDEX IF EXISTS idx_tenants_school_name`,
		`ALTER TABLE tenants DROP COLUMN IF EXISTS pesantren_name`,
		`ALTER TABLE tenants DROP COLUMN IF EXISTS school_name`,
	}

	for _, q := range queries {
		if _, err := db.ExecContext(ctx, q); err != nil {
			fmt.Printf("Warning during rollback: %v\n", err)
		}
	}

	return nil
}
