package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(up37NotificationTemplates, down37NotificationTemplates)
}

func up37NotificationTemplates(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS notification_templates (
			id VARCHAR(36) PRIMARY KEY,
			event_type VARCHAR(50) NOT NULL,
			channel VARCHAR(20) NOT NULL,
			template_name VARCHAR(100) NOT NULL,
			template_content TEXT NOT NULL,
			variables TEXT,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(event_type, channel)
		);

		-- Seed default templates
		INSERT INTO notification_templates (id, event_type, channel, template_name, template_content, variables, is_active) VALUES
		-- Registration templates
		('tmpl-reg-owner-tg', 'registration', 'telegram', 'Registrasi Baru', 
		 E'🎉 *Pendaftaran Baru!*\n\nNama: {{name}}\nInstitusi: {{institution}}\nEmail: {{email}}\nNomor HP: {{phone}}\nSubdomain: {{subdomain}}\n\nWaktu: {{timestamp}}', 
		 'name,institution,email,phone,subdomain,timestamp', true),
		
		('tmpl-reg-owner-wa', 'registration', 'whatsapp_owner', 'Registrasi Baru', 
		 E'🎉 *Pendaftaran Baru!*\n\nNama: {{name}}\nInstitusi: {{institution}}\nEmail: {{email}}\nNomor HP: {{phone}}\nSubdomain: {{subdomain}}\n\nWaktu: {{timestamp}}', 
		 'name,institution,email,phone,subdomain,timestamp', true),
		
		('tmpl-reg-tenant-wa', 'registration', 'whatsapp_tenant', 'Selamat Datang', 
		 E'Assalamu''alaikum {{name}},\n\nSelamat! Pendaftaran Anda di EduVera berhasil.\n\n📱 Akses dashboard: https://{{subdomain}}.eduvera.ve-lora.my.id\n📧 Email: {{email}}\n\nTerima kasih telah bergabung!\n\n- Tim EduVera', 
		 'name,subdomain,email', true),

		-- Upgrade templates
		('tmpl-upg-owner-tg', 'upgrade', 'telegram', 'Upgrade Paket', 
		 E'💎 *Upgrade Paket!*\n\nInstitusi: {{institution}}\nDari: {{from_tier}}\nKe: {{to_tier}}\nNominal: Rp {{amount}}\n\nWaktu: {{timestamp}}', 
		 'institution,from_tier,to_tier,amount,timestamp', true),
		
		('tmpl-upg-owner-wa', 'upgrade', 'whatsapp_owner', 'Upgrade Paket', 
		 E'💎 *Upgrade Paket!*\n\nInstitusi: {{institution}}\nDari: {{from_tier}}\nKe: {{to_tier}}\nNominal: Rp {{amount}}\n\nWaktu: {{timestamp}}', 
		 'institution,from_tier,to_tier,amount,timestamp', true),
		
		('tmpl-upg-tenant-wa', 'upgrade', 'whatsapp_tenant', 'Upgrade Berhasil', 
		 E'Assalamu''alaikum,\n\nUpgrade paket Anda berhasil!\n\n✅ Paket: {{to_tier}}\n📅 Berlaku hingga: {{expires_at}}\n\nFitur premium sudah aktif. Terima kasih!\n\n- Tim EduVera', 
		 'to_tier,expires_at', true),

		-- Payment templates (SPP)
		('tmpl-pay-owner-tg', 'payment', 'telegram', 'Pembayaran SPP', 
		 E'💰 *Pembayaran SPP Masuk!*\n\nInstitusi: {{institution}}\nSiswa: {{student_name}}\nNominal: Rp {{amount}}\nBulan: {{month}}\n\nWaktu: {{timestamp}}', 
		 'institution,student_name,amount,month,timestamp', true),
		
		('tmpl-pay-owner-wa', 'payment', 'whatsapp_owner', 'Pembayaran SPP', 
		 E'💰 *Pembayaran SPP Masuk!*\n\nInstitusi: {{institution}}\nSiswa: {{student_name}}\nNominal: Rp {{amount}}\nBulan: {{month}}\n\nWaktu: {{timestamp}}', 
		 'institution,student_name,amount,month,timestamp', true),
		
		('tmpl-pay-tenant-wa', 'payment', 'whatsapp_tenant', 'Konfirmasi Pembayaran', 
		 E'Assalamu''alaikum,\n\nPembayaran SPP berhasil diterima.\n\n✅ Siswa: {{student_name}}\n💰 Nominal: Rp {{amount}}\n📅 Bulan: {{month}}\n\nTerima kasih.\n\n- {{institution}}', 
		 'student_name,amount,month,institution', true)
		
		ON CONFLICT (event_type, channel) DO NOTHING;
	`)
	return err
}

func down37NotificationTemplates(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS notification_templates`)
	return err
}
