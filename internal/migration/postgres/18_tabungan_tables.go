package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upTabungan, downTabungan)
}

func upTabungan(ctx context.Context, tx *sql.Tx) error {
	query := `
			-- Tabungan Account per Siswa
			CREATE TABLE IF NOT EXISTS sekolah_tabungan (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				tenant_id UUID NOT NULL,
				santri_id UUID NOT NULL REFERENCES sekolah_siswa(id),
				saldo BIGINT DEFAULT 0,
				status VARCHAR(20) DEFAULT 'Wajib', -- Wajib, Sukarela
				created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				UNIQUE(tenant_id, santri_id)
			);
			CREATE TRIGGER update_sekolah_tabungan_updated_at
				BEFORE UPDATE ON sekolah_tabungan
				FOR EACH ROW
				EXECUTE FUNCTION update_updated_at_column();

			-- Mutasi Tabungan (Transactions)
			CREATE TABLE IF NOT EXISTS sekolah_tabungan_mutasi (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				tabungan_id UUID NOT NULL REFERENCES sekolah_tabungan(id) ON DELETE CASCADE,
				tenant_id UUID NOT NULL,
				tipe VARCHAR(20) NOT NULL, -- Debit (Masuk), Kredit (Keluar)
				nominal BIGINT NOT NULL,
				keterangan TEXT,
				petugas VARCHAR(100), -- User who input
				created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
			);
			`
	if _, err := tx.Exec(query); err != nil {
		return fmt.Errorf("failed to create tabungan tables: %w", err)
	}

	return nil
}

func downTabungan(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.Exec(`DROP TABLE IF EXISTS sekolah_tabungan_mutasi`); err != nil {
		return err
	}
	if _, err := tx.Exec(`DROP TABLE IF EXISTS sekolah_tabungan`); err != nil {
		return err
	}
	return nil
}
