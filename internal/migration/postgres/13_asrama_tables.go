package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAsramaTables, downAsramaTables)
}

func upAsramaTables(ctx context.Context, tx *sql.Tx) error {
	// Table: sekolah_asrama (Buildings)
	_, err := tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS sekolah_asrama (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL,
			nama VARCHAR(100) NOT NULL,
			jenis VARCHAR(50) NOT NULL, -- Putra, Putri
			musyrif_id UUID, -- FK to sekolah_guru (Optional)
			status VARCHAR(50) DEFAULT 'Aktif',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_sekolah_asrama_tenant ON sekolah_asrama(tenant_id);
	`)
	if err != nil {
		return err
	}

	// Table: sekolah_kamar (Rooms)
	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS sekolah_kamar (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL,
			asrama_id UUID NOT NULL REFERENCES sekolah_asrama(id) ON DELETE CASCADE,
			nomor VARCHAR(50) NOT NULL,
			kapasitas INT NOT NULL DEFAULT 0,
			status VARCHAR(50) DEFAULT 'Tersedia', -- Tersedia, Penuh, Perbaikan
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_sekolah_kamar_tenant ON sekolah_kamar(tenant_id);
		CREATE INDEX IF NOT EXISTS idx_sekolah_kamar_asrama ON sekolah_kamar(asrama_id);
	`)
	if err != nil {
		return err
	}

	// Table: sekolah_penempatan (Student Placement)
	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS sekolah_penempatan (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL,
			santri_id UUID NOT NULL REFERENCES sekolah_siswa(id) ON DELETE CASCADE,
			kamar_id UUID NOT NULL REFERENCES sekolah_kamar(id) ON DELETE CASCADE,
			tanggal_masuk DATE NOT NULL DEFAULT CURRENT_DATE,
			tanggal_keluar DATE, -- Null if currently active
			status VARCHAR(50) DEFAULT 'Aktif', -- Aktif, Pindah, Keluar
			keterangan TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_sekolah_penempatan_tenant ON sekolah_penempatan(tenant_id);
		CREATE INDEX IF NOT EXISTS idx_sekolah_penempatan_santri ON sekolah_penempatan(santri_id);
		CREATE INDEX IF NOT EXISTS idx_sekolah_penempatan_kamar ON sekolah_penempatan(kamar_id);
	`)
	if err != nil {
		return err
	}

	// Trigger for Updated At
	_, err = tx.ExecContext(ctx, `
		CREATE TRIGGER update_sekolah_asrama_updated_at BEFORE UPDATE ON sekolah_asrama FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
		CREATE TRIGGER update_sekolah_kamar_updated_at BEFORE UPDATE ON sekolah_kamar FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
		CREATE TRIGGER update_sekolah_penempatan_updated_at BEFORE UPDATE ON sekolah_penempatan FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
	`)
	return err
}

func downAsramaTables(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		DROP TABLE IF EXISTS sekolah_penempatan;
		DROP TABLE IF EXISTS sekolah_kamar;
		DROP TABLE IF EXISTS sekolah_asrama;
	`)
	return err
}
