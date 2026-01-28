package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upRapor, downRapor)
}

func upRapor(ctx context.Context, tx *sql.Tx) error {
	// 1. Rapor Periode (e.g. "Ganjil 2023/2024")
	query := `
			CREATE TABLE IF NOT EXISTS sekolah_rapor_periode (
				id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				tenant_id UUID NOT NULL,
				nama VARCHAR(100) NOT NULL, -- e.g. "Ganjil 2023/2024"
				tanggal_mulai DATE,
				tanggal_akhir DATE,
				is_active BOOLEAN DEFAULT false,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
			);
			CREATE TRIGGER update_sekolah_rapor_periode_updated_at
				BEFORE UPDATE ON sekolah_rapor_periode
				FOR EACH ROW
				EXECUTE FUNCTION update_updated_at_column();

			-- 2. Rapor (The Report Card Header per Santri per Periode)
			CREATE TABLE IF NOT EXISTS sekolah_rapor (
				id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				tenant_id UUID NOT NULL,
				periode_id UUID NOT NULL REFERENCES sekolah_rapor_periode(id),
				santri_id UUID NOT NULL REFERENCES sekolah_siswa(id),
				status VARCHAR(50) DEFAULT 'Draft', -- Draft, Published, Archived
				catatan_wali_kelas TEXT,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				UNIQUE(periode_id, santri_id)
			);
			CREATE TRIGGER update_sekolah_rapor_updated_at
				BEFORE UPDATE ON sekolah_rapor
				FOR EACH ROW
				EXECUTE FUNCTION update_updated_at_column();

			-- 3. Rapor Nilai (The specific grades inside the rapor)
			CREATE TABLE IF NOT EXISTS sekolah_rapor_nilai (
				id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				rapor_id UUID NOT NULL REFERENCES sekolah_rapor(id) ON DELETE CASCADE,
				kategori VARCHAR(50) NOT NULL, -- Tahfidz, Diniyah, Musyrif, Akademik
				jenis VARCHAR(100), -- e.g. "Hifdzul Quran", "Fiqih", "Akhlak"
				nilai VARCHAR(50), -- Can be numeric "90" or qualitative "A", "Mumtaz"
				keterangan TEXT,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
			);
			CREATE TRIGGER update_sekolah_rapor_nilai_updated_at
				BEFORE UPDATE ON sekolah_rapor_nilai
				FOR EACH ROW
				EXECUTE FUNCTION update_updated_at_column();
			`
	if _, err := tx.Exec(query); err != nil {
		return fmt.Errorf("failed to create rapor tables: %w", err)
	}

	return nil
}

func downRapor(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.Exec(`DROP TABLE IF EXISTS sekolah_rapor_nilai`); err != nil {
		return err
	}
	if _, err := tx.Exec(`DROP TABLE IF EXISTS sekolah_rapor`); err != nil {
		return err
	}
	if _, err := tx.Exec(`DROP TABLE IF EXISTS sekolah_rapor_periode`); err != nil {
		return err
	}
	return nil
}
