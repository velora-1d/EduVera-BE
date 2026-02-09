#!/bin/bash

# Configuration
VPS_USER="ubuntu"
VPS_IP="43.156.132.218"
PROJECT_DIR="~/EduVera"

# Get script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BE_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_ROOT="$(dirname "$BE_DIR")"

echo "🚀 Starting Full Stack Deploy to VPS ($VPS_IP)..."
echo "📁 Deploying from: $PROJECT_ROOT"

# 1. Sync Files (from project root, excluding unnecessary folders)
echo "📦 Syncing files..."
rsync -avz --progress \
  --exclude '.git' \
  --exclude 'node_modules' \
  --exclude '.next' \
  --exclude 'out' \
  --exclude 'EduVera-FE' \
  --exclude 'EduVera Brief' \
  --exclude 'EduVera-BE/tmp' \
  --exclude '.vscode' \
  "$PROJECT_ROOT/" "$VPS_USER@$VPS_IP:$PROJECT_DIR"

echo "🛠 Preparing environment on VPS..."

# 2. Setup and deploy on VPS
ssh "$VPS_USER@$VPS_IP" << EOF
  cd $PROJECT_DIR/EduVera-BE
  
  # Link .env if needed
  if [ -f .env ]; then
    echo "✅ .env found in EduVera-BE"
  fi
  
  # Ensure necessary ports are open or handled by docker
  echo "🏗 Rebuilding and restarting containers..."
  
  # Stop and remove legacy containers
  docker stop eduvera_backend eduvera_postgres eduvera_redis eduvera_rabbitmq eduvera-whatsback 2>/dev/null || true
  docker rm eduvera_backend eduvera_postgres eduvera_redis eduvera_rabbitmq eduvera-whatsback 2>/dev/null || true
  
  # Ensure we have a working docker-compose v2
  if [ ! -f ./docker-compose ]; then
    echo "⬇️ Downloading docker-compose v2..."
    curl -SL https://github.com/docker/compose/releases/download/v2.24.5/docker-compose-linux-x86_64 -o ./docker-compose
    chmod +x ./docker-compose
  fi

  echo "🧹 Cleaning up old containers..."
  docker ps -a --filter "name=eduvera" -q | xargs -r docker rm -f
  
  echo "🏗 Rebuilding and restarting containers..."
  ./docker-compose -f deploy/docker-compose.vps.yml up -d --build --remove-orphans
EOF

echo "✅ Deployment completed successfully!"
echo "🌐 Backend: https://api-eduvera.ve-lora.my.id"
