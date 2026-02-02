#!/bin/bash

# Configuration
VPS_USER="ubuntu"
VPS_IP="43.156.132.218"
PROJECT_DIR="~/EduVera/EduVera-BE"
LOCAL_DIR="."

echo "🚀 Starting Quick Deploy to VPS ($VPS_IP)..."

# 1. Sync Files (RSYNC)
# Mengirim file yang berubah saja (Hemat Bandwidth & Waktu)
# Exclude: .git (backup manual), .env (biar prod config aman), main (binary)
rsync -avz \
  --exclude '.git' \
  --exclude '.env' \
  --exclude 'main' \
  --exclude 'tmp' \
  "$LOCAL_DIR" "$VPS_USER@$VPS_IP:$PROJECT_DIR"

echo "🛠 Rebuilding Backend on VPS..."

# 2. Rebuild & Restart Container
# Hanya rebuild service 'app' (backend)
ssh "$VPS_USER@$VPS_IP" "cd $PROJECT_DIR && docker-compose -f docker-compose.prod.yml up -d --build app"

echo "✅ Backend Deployed Successfully!"
echo "ℹ️  Jangan lupa 'git push' manual untuk backup code ke GitHub."
