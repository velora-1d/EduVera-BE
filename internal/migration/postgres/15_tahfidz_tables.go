package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upTahfidz, downTahfidz)
}

func upTahfidz(ctx context.Context, tx *sql.Tx) error {
	// 1. Tahfidz Setoran (Records of Memorization)
	query := `
			CREATE TABLE IF NOT EXISTS sekolah_tahfidz_setoran (
				id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				tenant_id UUID NOT NULL,
				santri_id UUID NOT NULL REFERENCES sekolah_siswa(id),
				ustadz_id UUID REFERENCES sekolah_guru(id),
				tanggal DATE NOT NULL DEFAULT CURRENT_DATE,
				juz INT,
				surah VARCHAR(100),
				ayat_awal INT,
				ayat_akhir INT,
				tipe VARCHAR(50) NOT NULL, -- Ziyadah (Hafalan Baru), Murajaah (Ulang)
				kualitas VARCHAR(50), -- Lancar, Kurang, Ulang
				catatan TEXT,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
			);
			CREATE TRIGGER update_sekolah_tahfidz_setoran_updated_at
				BEFORE UPDATE ON sekolah_tahfidz_setoran
				FOR EACH ROW
				EXECUTE FUNCTION update_updated_at_column();
			`
	if _, err := tx.Exec(query); err != nil {
		return fmt.Errorf("failed to create sekolah_tahfidz_setoran: %w", err)
	}

	return nil
}

func downTahfidz(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.Exec(`DROP TABLE IF EXISTS sekolah_tahfidz_setoran`); err != nil {
		return err
	}
	return nil
}
