package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upKalender, downKalender)
}

func upKalender(ctx context.Context, tx *sql.Tx) error {
	query := `
			CREATE TABLE IF NOT EXISTS sekolah_kalender (
				id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				tenant_id UUID NOT NULL,
				title VARCHAR(255) NOT NULL,
				start_date DATE NOT NULL,
				end_date DATE NOT NULL,
				category VARCHAR(50), -- Libur, Ujian, Kegiatan, Penting
				description TEXT,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
			);
			CREATE TRIGGER update_sekolah_kalender_updated_at
				BEFORE UPDATE ON sekolah_kalender
				FOR EACH ROW
				EXECUTE FUNCTION update_updated_at_column();
			`
	if _, err := tx.Exec(query); err != nil {
		return fmt.Errorf("failed to create kalender table: %w", err)
	}
	return nil
}

func downKalender(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE IF EXISTS sekolah_kalender`)
	return err
}
