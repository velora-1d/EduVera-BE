package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upKepesantrenan, downKepesantrenan)
}

func upKepesantrenan(ctx context.Context, tx *sql.Tx) error {
	// 0. Sekolah Guru (Dependency for Perizinan)
	if _, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS sekolah_guru (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL,
			nip VARCHAR(50),
			nama VARCHAR(255) NOT NULL,
			jenis VARCHAR(50), -- Guru Mapel, Guru Kelas
			status VARCHAR(50), -- PNS, Honorer
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_sekolah_guru_tenant ON sekolah_guru(tenant_id);
		CREATE TRIGGER update_sekolah_guru_updated_at BEFORE UPDATE ON sekolah_guru FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
	`); err != nil {
		return fmt.Errorf("failed to create sekolah_guru: %w", err)
	}

	// 1. Pelanggaran Aturan (Master Data)
	query := `
			CREATE TABLE IF NOT EXISTS sekolah_pelanggaran_aturan (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				tenant_id UUID NOT NULL,
				judul VARCHAR(255) NOT NULL,
				kategori VARCHAR(100) NOT NULL,
				poin INT NOT NULL DEFAULT 0,
				level VARCHAR(50) NOT NULL, -- Ringan, Sedang, Berat
				created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
			);
			CREATE TRIGGER update_sekolah_pelanggaran_aturan_updated_at
				BEFORE UPDATE ON sekolah_pelanggaran_aturan
				FOR EACH ROW
				EXECUTE FUNCTION update_updated_at_column();
			`
	if _, err := tx.Exec(query); err != nil {
		return fmt.Errorf("failed to create sekolah_pelanggaran_aturan: %w", err)
	}

	// 2. Pelanggaran Siswa (Records)
	query = `
			CREATE TABLE IF NOT EXISTS sekolah_pelanggaran_siswa (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				tenant_id UUID NOT NULL,
				santri_id UUID NOT NULL REFERENCES sekolah_siswa(id),
				aturan_id UUID REFERENCES sekolah_pelanggaran_aturan(id),
				tanggal TIMESTAMP WITH TIME ZONE NOT NULL,
				poin INT NOT NULL, -- Snapshot of points at the time
				keterangan TEXT,
				status VARCHAR(50) DEFAULT 'Pending', -- Pending, Diproses, Selesai
				sanksi VARCHAR(255),
				created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
			);
			CREATE TRIGGER update_sekolah_pelanggaran_siswa_updated_at
				BEFORE UPDATE ON sekolah_pelanggaran_siswa
				FOR EACH ROW
				EXECUTE FUNCTION update_updated_at_column();
			`
	if _, err := tx.Exec(query); err != nil {
		return fmt.Errorf("failed to create sekolah_pelanggaran_siswa: %w", err)
	}

	// 3. Perizinan (Permissions)
	query = `
			CREATE TABLE IF NOT EXISTS sekolah_perizinan (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				tenant_id UUID NOT NULL,
				santri_id UUID NOT NULL REFERENCES sekolah_siswa(id),
				tipe VARCHAR(50) NOT NULL, -- Izin Pulang, Izin Keluar, Izin Sakit
				alasan TEXT,
				dari TIMESTAMP WITH TIME ZONE NOT NULL,
				sampai TIMESTAMP WITH TIME ZONE NOT NULL,
				status VARCHAR(50) DEFAULT 'Pending', -- Pending, Disetujui, Ditolak
				penyetuju_id UUID REFERENCES sekolah_guru(id),
				created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
			);
			CREATE TRIGGER update_sekolah_perizinan_updated_at
				BEFORE UPDATE ON sekolah_perizinan
				FOR EACH ROW
				EXECUTE FUNCTION update_updated_at_column();
			`
	if _, err := tx.Exec(query); err != nil {
		return fmt.Errorf("failed to create sekolah_perizinan: %w", err)
	}

	return nil
}

func downKepesantrenan(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.Exec(`DROP TABLE IF EXISTS sekolah_perizinan`); err != nil {
		return err
	}
	if _, err := tx.Exec(`DROP TABLE IF EXISTS sekolah_pelanggaran_siswa`); err != nil {
		return err
	}
	if _, err := tx.Exec(`DROP TABLE IF EXISTS sekolah_pelanggaran_aturan`); err != nil {
		return err
	}
	return nil
}
