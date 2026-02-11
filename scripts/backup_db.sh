#!/bin/bash

# Configuration
BACKUP_DIR="/root/backups/eduvera"
CONTAINER_NAME="eduvera_postgres"
DB_USER="prabogo"
DB_NAME="prabogo"
DATE=$(date +%Y%m%d_%H%M%S)
FILENAME="backup_${DATE}.sql.gz"

# Create backup directory if not exists
mkdir -p $BACKUP_DIR

# Dump and compress
echo "Starting backup for $DB_NAME..."
docker exec $CONTAINER_NAME pg_dump -U $DB_USER $DB_NAME | gzip > "$BACKUP_DIR/$FILENAME"

# Verify backup success
if [ $? -eq 0 ]; then
  echo "Backup successful: $FILENAME"
  
  # Retention: Delete backups older than 7 days
  find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +7 -delete
  echo "Old backups cleaned up."
else
  echo "Backup FAILED!"
  exit 1
fi
