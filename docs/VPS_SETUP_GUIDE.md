# Panduan Deploy Backend ke VPS (EduVera)

Panduan ini khusus untuk men-deploy **Backend API** ke VPS. Frontend dan Mobile App nanti akan connect ke IP VPS ini.

## 1. Login ke VPS
Dari terminal laptop Anda:
```bash
ssh root@43.156.132.218
```

## 2. Install Docker (Sekali Saja)
Copy-paste perintah ini ke terminal VPS untuk install Docker & Docker Compose:

```bash
apt update && apt upgrade -y && \
apt install -y git curl && \
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg && \
echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null && \
apt update && apt install -y docker-ce docker-ce-cli containerd.io && \
curl -L "https://github.com/docker/compose/releases/download/v2.24.5/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose && \
chmod +x /usr/local/bin/docker-compose && \
echo "✅ Docker Siap!"
```

## 3. Clone Repository Backend
```bash
cd /opt
# Clone repo EduVera-BE yang sudah ada
git clone https://github.com/velora-1d/EduVera-BE.git eduvera-backend
cd eduvera-backend
```
*(Login pakai Username GitHub `velora-1d` dan Token Anda)*

## 4. Setup Environment
Buat file `.env` di VPS:
```bash
nano .env
```

**Isi file .env (Copy Paste ini):**
```env
# SERVER CONFIG
SERVER_PORT=8000
APP_MODE=release

# DATABASE INTERNAL
POSTGRES_USER=prabogo
POSTGRES_PASSWORD=changeme_secure_pass
POSTGRES_DB=prabogo
REDIS_PASSWORD=changeme_secure_pass
RABBITMQ_USER=prabogo
RABBITMQ_PASSWORD=changeme_secure_pass

# URLs (PENTING: Ganti IP di sini jika perlu)
DATABASE_URL=postgres://prabogo:changeme_secure_pass@postgres:5432/prabogo?sslmode=disable
REDIS_URL=redis://:changeme_secure_pass@redis:6379/0
```
Simpan: `Ctrl+O` -> `Enter` -> `Ctrl+X`.

## 5. Jalankan Backend 🚀
```bash
# Jalankan pakai config production
docker-compose -f docker-compose.prod.yml up -d --build
```

## 6. Cek Status
```bash
docker ps
```
Pastikan `eduvera_backend`, `postgres`, `redis` statusnya **Up**.

---
**Info Koneksi:**
- **Backend API:** `http://43.156.132.218:8000`
- Gunakan IP ini untuk config di Mobile App dan Frontend.
