package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upFixStudentFK, downFixStudentFK)
}

func upFixStudentFK(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		-- Add Foreign Key constraint to students table
		-- This ensures that when a tenant is deleted, their students are also deleted (Cascade)
		-- And prevents creating students for non-existent tenants
		ALTER TABLE students
		ADD CONSTRAINT fk_students_tenant
		FOREIGN KEY (tenant_id)
		REFERENCES tenants(id)
		ON DELETE CASCADE;
	`)
	return err
}

func downFixStudentFK(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		ALTER TABLE students
		DROP CONSTRAINT IF EXISTS fk_students_tenant;
	`)
	return err
}
