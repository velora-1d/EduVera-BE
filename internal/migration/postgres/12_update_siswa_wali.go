package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upSekolahSiswa, downSekolahSiswa)
}

func upSekolahSiswa(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS sekolah_siswa (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL,
			nis VARCHAR(50),
			nisn VARCHAR(50),
			nama VARCHAR(255) NOT NULL,
			kelas_id UUID REFERENCES sekolah_kelas(id),
			alamat TEXT,
			nama_wali VARCHAR(100),
			no_hp_wali VARCHAR(20),
			status VARCHAR(20) DEFAULT 'Aktif',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_sekolah_siswa_tenant ON sekolah_siswa(tenant_id);
		CREATE INDEX IF NOT EXISTS idx_sekolah_siswa_kelas ON sekolah_siswa(kelas_id);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_sekolah_siswa_nis ON sekolah_siswa(tenant_id, nis) WHERE nis IS NOT NULL AND nis != '';
		
		CREATE TRIGGER update_sekolah_siswa_updated_at
			BEFORE UPDATE ON sekolah_siswa
			FOR EACH ROW
			EXECUTE FUNCTION update_updated_at_column();
	`)
	return err
}

func downSekolahSiswa(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		DROP TABLE IF EXISTS sekolah_siswa;
	`)
	return err
}
