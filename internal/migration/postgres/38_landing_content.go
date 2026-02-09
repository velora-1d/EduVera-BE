package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upLandingContent, downLandingContent)
}

func upLandingContent(ctx context.Context, tx *sql.Tx) error {
	// Create table
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS landing_content (
			key VARCHAR(255) PRIMARY KEY,
			value JSONB NOT NULL,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}

	// Seed data
	// Pricing Plans - MATCHING Frontend Interface PlanPricing
	_, err = tx.Exec(`
		INSERT INTO landing_content (key, value)
		SELECT 'pricing_plans', '{
			"sekolah": {
				"basic": {"monthly": 299000, "annual": 2999000},
				"premium": {"monthly": 499000, "annual": 4999000}
			},
			"pesantren": {
				"basic": {"monthly": 299000, "annual": 2999000},
				"premium": {"monthly": 499000, "annual": 4999000}
			},
			"hybrid": {
				"basic": {"monthly": 449000, "annual": 4499000},
				"premium": {"monthly": 699000, "annual": 6999999}
			}
		}'::jsonb
		WHERE NOT EXISTS (SELECT 1 FROM landing_content WHERE key = 'pricing_plans');
	`)
	if err != nil {
		return err
	}

	// Features School
	_, err = tx.Exec(`
		INSERT INTO landing_content (key, value)
		SELECT 'features_school', '[{"title": "Dashboard Akademik", "desc": "Ringkasan data & status progres.", "icon": "LayoutDashboard"}, {"title": "Data Akademik", "desc": "Profil lengkap & pemetaan Mapel.", "icon": "Users"}, {"title": "Pembelajaran", "desc": "Penugasan guru & jadwal harian.", "icon": "School"}, {"title": "Kurikulum Nasional", "desc": "Support K13 & Merdeka.", "icon": "BookMarked"}, {"title": "E-Rapor Nasional", "desc": "Otomasi generate PDF rapor.", "icon": "ClipboardCheck"}, {"title": "SDM & HR", "desc": "Struktur & absensi pegawai.", "icon": "UserCheck"}, {"title": "ERP Keuangan", "desc": "SPP, BOS & Penggajian.", "icon": "Wallet"}, {"title": "Kalender", "desc": "Reminder kegiatan otomatis.", "icon": "Calendar"}, {"title": "Pusat Laporan", "desc": "Rekap data komprehensif.", "icon": "BarChart3"}, {"title": "Config", "desc": "Role & fitur management.", "icon": "Settings"}]'::jsonb
		WHERE NOT EXISTS (SELECT 1 FROM landing_content WHERE key = 'features_school');
	`)
	if err != nil {
		return err
	}

	// Features Pesantren
	_, err = tx.Exec(`
		INSERT INTO landing_content (key, value)
		SELECT 'features_pesantren', '[{"title": "Dashboard Pondok", "desc": "Monitor santri & kas harian.", "icon": "LayoutDashboard"}, {"title": "Database Santri", "desc": "Profil & riwayat mukim.", "icon": "Users"}, {"title": "Asrama", "desc": "Plotting & absensi asrama.", "icon": "Bed"}, {"title": "Kedisiplinan", "desc": "Poin & perizinan terpadu.", "icon": "Scale"}, {"title": "Tahfidz", "desc": "Target & setoran harian.", "icon": "GraduationCap"}, {"title": "Kitab", "desc": "Halaqah & pemahaman kitab.", "icon": "BookOpen"}, {"title": "E-Rapor Pondok", "desc": "Rapor gabungan PDF.", "icon": "ClipboardCheck"}, {"title": "SDM Ustadz", "desc": "Honor & insentif ustadz.", "icon": "UserCheck"}, {"title": "Keuangan ERP", "desc": "Pemasukan & pengeluaran.", "icon": "Wallet"}, {"title": "Hijriah", "desc": "Kalender & kegiatan haflah.", "icon": "Calendar"}, {"title": "Audit Report", "desc": "Laporan kepesantrenan.", "icon": "FileText"}]'::jsonb
		WHERE NOT EXISTS (SELECT 1 FROM landing_content WHERE key = 'features_pesantren');
	`)
	if err != nil {
		return err
	}

	return nil
}

func downLandingContent(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE IF EXISTS landing_content;")
	return err
}
