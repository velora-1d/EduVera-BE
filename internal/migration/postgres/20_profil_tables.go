package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upProfil, downProfil)
}

func upProfil(ctx context.Context, tx *sql.Tx) error {
	query := `
			CREATE TABLE IF NOT EXISTS sekolah_profil (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				tenant_id UUID NOT NULL UNIQUE,
				jenis_pesantren VARCHAR(50), -- Salaf, Khalaf, Terpadu
				deskripsi TEXT,
				website VARCHAR(100),
				email_kontak VARCHAR(100),
				no_telp_kontak VARCHAR(50),
				logo_url TEXT,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
			);
			CREATE TRIGGER update_sekolah_profil_updated_at
				BEFORE UPDATE ON sekolah_profil
				FOR EACH ROW
				EXECUTE FUNCTION update_updated_at_column();
			`
	if _, err := tx.Exec(query); err != nil {
		return fmt.Errorf("failed to create profil table: %w", err)
	}
	return nil
}

func downProfil(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE IF EXISTS sekolah_profil`)
	return err
}
