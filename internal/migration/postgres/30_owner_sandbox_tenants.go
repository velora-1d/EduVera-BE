package migrations

import (
	"context"
	"database/sql"
	"os"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upOwnerSandboxTenants, downOwnerSandboxTenants)
}

func upOwnerSandboxTenants(ctx context.Context, tx *sql.Tx) error {
	// Add sandbox columns to tenants table
	_, err := tx.Exec(`
		ALTER TABLE tenants ADD COLUMN IF NOT EXISTS is_sandbox BOOLEAN DEFAULT FALSE;
		ALTER TABLE tenants ADD COLUMN IF NOT EXISTS owner_id UUID REFERENCES users(id);
		
		CREATE INDEX IF NOT EXISTS idx_tenants_is_sandbox ON tenants(is_sandbox);
		CREATE INDEX IF NOT EXISTS idx_tenants_owner_id ON tenants(owner_id);
	`)
	if err != nil {
		return err
	}

	// Get owner email from environment or use default
	ownerEmail := os.Getenv("OWNER_EMAIL")
	if ownerEmail == "" {
		ownerEmail = "nawawimahinutsman@gmail.com"
	}

	// Create sandbox tenants for owner
	// First, get owner user ID
	var ownerID string
	err = tx.QueryRow(`SELECT id FROM users WHERE email = $1`, ownerEmail).Scan(&ownerID)
	if err != nil {
		// Owner might not exist yet, skip creating sandbox tenants
		// They will be created when owner first logs in
		return nil
	}

	// Create sandbox tenants with PREMIUM tier and ACTIVE status (no expiration)
	sandboxTenants := []struct {
		name      string
		subdomain string
		planType  string
	}{
		{"Sandbox Sekolah EduVera", "sandbox-sekolah", "sekolah"},
		{"Sandbox Pesantren EduVera", "sandbox-pesantren", "pesantren"},
		{"Sandbox Hybrid EduVera", "sandbox-hybrid", "hybrid"},
	}

	for _, t := range sandboxTenants {
		_, err = tx.Exec(`
			INSERT INTO tenants (name, subdomain, plan_type, subscription_tier, status, is_sandbox, owner_id, created_at, updated_at)
			VALUES ($1, $2, $3, 'premium', 'active', true, $4, NOW(), NOW())
			ON CONFLICT (subdomain) DO UPDATE SET
				is_sandbox = true,
				owner_id = $4,
				subscription_tier = 'premium',
				status = 'active',
				updated_at = NOW()
		`, t.name, t.subdomain, t.planType, ownerID)
		if err != nil {
			return err
		}
	}

	return nil
}

func downOwnerSandboxTenants(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
		DELETE FROM tenants WHERE is_sandbox = true;
		
		DROP INDEX IF EXISTS idx_tenants_owner_id;
		DROP INDEX IF EXISTS idx_tenants_is_sandbox;
		
		ALTER TABLE tenants DROP COLUMN IF EXISTS owner_id;
		ALTER TABLE tenants DROP COLUMN IF EXISTS is_sandbox;
	`)
	return err
}
