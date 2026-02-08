package migrations

import (
	"context"
	"database/sql"
	"os"

	"github.com/pressly/goose/v3"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	goose.AddMigrationContext(upSeedOwner, downSeedOwner)
}

func upSeedOwner(ctx context.Context, tx *sql.Tx) error {
	ownerEmail := os.Getenv("OWNER_EMAIL")
	ownerPassword := os.Getenv("OWNER_PASSWORD")

	if ownerEmail == "" || ownerPassword == "" {
		// Skip seeding if env vars not set
		return nil
	}

	// Check if owner already exists
	var exists bool
	err := tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", ownerEmail).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		// Owner already exists, just ensure is_owner flag is set
		_, err = tx.ExecContext(ctx, "UPDATE users SET is_owner = true, is_active = true WHERE email = $1", ownerEmail)
		return err
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(ownerPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Insert the owner user
	_, err = tx.ExecContext(ctx, `
		INSERT INTO users (name, email, password_hash, role, is_owner, is_active, created_at, updated_at)
		VALUES ('Owner', $1, $2, 'super_admin', true, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, ownerEmail, string(hashedPassword))

	return err
}

func downSeedOwner(ctx context.Context, tx *sql.Tx) error {
	ownerEmail := os.Getenv("OWNER_EMAIL")
	if ownerEmail == "" {
		return nil
	}

	_, err := tx.ExecContext(ctx, "DELETE FROM users WHERE email = $1 AND is_owner = true", ownerEmail)
	return err
}
