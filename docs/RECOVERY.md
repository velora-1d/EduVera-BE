# Disaster Recovery Guide — EduVera

> **PENTING**: Simpan dokumen ini di tempat aman (bukan hanya di server). Print jika perlu.

## 1. Lokasi Backup
Backup otomatis berjalan setiap hari jam 02:00 WIB.
- **Lokasi di VPS**: `/root/backups/eduvera/`
- **Format**: `backup_YYYYMMDD_HHMMSS.sql.gz` (Kompresi GZIP)
- **Retensi**: Menyimpan 7 hari terakhir.

## 2. Cara Restore Database (Jika Data Hilang/Corrupt)

### Langkah 1: Pastikan Database Bersih
Jika database rusak parah, drop dulu (HATI-HATI!):
```bash
docker exec -it eduvera_postgres psql -U prabogo -d telegram_bot
# Di dalam psql console:
DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
\q
```

### Langkah 2: Restore dari File Backup
Pilih file backup terakhir yang valid, misal `backup_20260211.sql.gz`.

```bash
# 1. Unzip file (piping on the fly)
gunzip -c /root/backups/eduvera/backup_20260211.sql.gz | docker exec -i eduvera_postgres psql -U prabogo -d prabogo
```

### Langkah 3: Verify
Cek apakah tabel dan data sudah kembali:
```bash
docker exec -it eduvera_postgres psql -U prabogo -d prabogo -c "\dt"
```

## 3. Full Server Recovery (Jika VPS Meledak)

1. **Setup VPS Baru**: Install Docker & Docker Compose.
2. **Clone Repo**: `git clone ...`
3. **Setup ENV**: Copy file `.env` (pastikan Anda punya backup `.env` di local laptop!).
4. **Copy Backup Data**: SCP file `.sql.gz` dari backup storage ke VPS baru.
5. **Start Services**: `docker-compose -f deploy/docker-compose.vps.yml up -d`
6. **Restore DB**: Lakukan langkah Restore di atas.
