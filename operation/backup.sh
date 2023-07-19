#!/bin/bash
#
# backup-mysql.sh 1.0.5
#
# Dumps all databases to seperate files.
# All files are created in a folder named by the current date.
# Folders exceeding the defined hold time are purged automatically.
#
# (c)2015-2019 Harald Schneider
#

# Setup.start
#
HOLD_DAYS=30
TIMESTAMP=$(date +"%F")
BACKUP_DIR="/srv/festivals-identity-server/backups"
CREDENTIALS_FILE="/usr/local/festivals-identity-server/mysql.conf"

#
# Setup.end
# Check and auto-repair all databases first
#
echo
echo "Checking all databases - this can take a while ..."
sleep 1
echo
mysqlcheck --defaults-extra-file=$CREDENTIALS_FILE --auto-repair --all-databases
sleep 1

# Backup
#
echo "Starting backup ..."
sleep 1
mkdir -p "$BACKUP_DIR/$TIMESTAMP"
mysqldump --defaults-extra-file=$CREDENTIALS_FILE --force --opt --no-tablespaces --databases 'festivals_identity_database' | gzip > "$BACKUP_DIR/$TIMESTAMP/festivals_identity_database-$(date "+%F-%H-%M-%S").gz"

# Cleaning up
#
echo "Cleaning up ..."
find $BACKUP_DIR -maxdepth 1 -mindepth 1 -type d -mtime +$HOLD_DAYS -exec rm -rf {} \;
sleep 1
echo "-- DONE!"
sleep 1
