#!/bin/bash
#
# install.sh 1.0.0
#
# Enables the firewall, installs the newest mysql, starts it as a service,
# configures it to be used as the database server for the FestivalsIdentityServer and setup
# the backup routines.
#
# (c)2020-2023 Simon Gaus
#

# Check if all passwords are supplied
#
if [ $# -ne 3 ]; then
    echo "$0: usage: sudo ./install.sh <mysql_root_pw> <mysql_backup_pw> <read_write_pw>"
    exit 1
fi

# Store passwords in variables
#
root_password=$1
backup_password=$2
read_write_password=$3
echo "All necessary passwords are provided and valid."
sleep 1

# Test for web server user
#
WEB_USER="www-data"
id -u "$WEB_USER" &>/dev/null;
if [ $? -ne 0 ]; then
  WEB_USER="www"
  if [ $? -ne 0 ]; then
    echo "Failed to find user to run web server. Exiting."
    exit 1
  fi
fi

# Store username in variable
#
# Usign this because whoami would return root if the script is called with sudo!
#
current_user=$(who mom likes | awk '{print $1}')

# Create and move to project directory
#
echo "Creating project directory"
sleep 1
mkdir -p /usr/local/festivals-identity-server || { echo "Failed to create project directory. Exiting." ; exit 1; }
cd /usr/local/festivals-identity-server || { echo "Failed to access project directory. Exiting." ; exit 1; }
chown -R "$current_user":"$current_user" .
chmod -R 761 .

# Install mysql if needed.
#
echo "Installing mysql-server..."
apt-get install mysql-server -y > /dev/null;

# Enables and configures the firewall.
# Supported firewalls: ufw and firewalld
# This step is skipped under macOS.
#
ufw allow mysql
echo "Added mysql service to ufw rules"
sleep 1

# Launch mysql on startup
#
systemctl enable mysql > /dev/null
systemctl start mysql > /dev/null
echo "Enabled and started mysql systemd service."
sleep 1

# Install mysql credential file
#
echo "Installing mysql credential file"
sleep 1
credentialsFile=/usr/local/festivals-identity-server/mysql.conf
cat << EOF > $credentialsFile
# festivals-identity-server configuration file v1.0
# TOML 1.0.0-rc.2+

[client]
user = 'festivals.identity.backup'
password = '$backup_password'
host = 'localhost'
EOF

chown -R "$current_user":"$current_user" /usr/local/festivals-identity-server/mysql.conf
chmod -R 761 /usr/local/festivals-identity-server/mysql.conf

# Download and run mysql secure script
#
echo "Downloading database security script"
curl --progress-bar -L -o secure-mysql.sh https://raw.githubusercontent.com/Festivals-App/festivals-identity-server/master/operation/secure-mysql.sh
chmod +x secure-mysql.sh
./secure-mysql.sh "$root_password"

# Download database creation script
#
echo "Downloading database creation script..."
curl --progress-bar -L -o create_database.sql https://raw.githubusercontent.com/Festivals-App/festivals-identity-server/master/database/create_database.sql

# Run database creation script and configure users
#
echo "Configuring mysql"
sleep 1
mysql -e "source /usr/local/festivals-identity-server/create_database.sql"
echo "Creating local backup user..."
mysql -e "CREATE USER 'festivals.identity.backup'@'localhost' IDENTIFIED BY '$backup_password';"
mysql -e "GRANT ALL PRIVILEGES ON festivals_identity_database.* TO 'festivals.identity.backup'@'localhost';"
sleep 1
echo "Creating local read/write user..."
mysql -e "CREATE USER 'festivals.identity.writer'@'localhost' IDENTIFIED BY '$read_write_password';"
mysql -e "GRANT SELECT, INSERT, UPDATE, DELETE ON festivals_identity_database.* TO 'festivals.identity.writer'@'localhost';"
sleep 1
mysql -e "FLUSH PRIVILEGES;"

# Create the backup directory
#
echo "Create backup directory"
sleep 1
mkdir -p /srv/festivals-identity-server/backups || { echo "Failed to create backup directory. Exiting." ; exit 1; }
cd /srv/festivals-identity-server/backups || { echo "Failed to access backup directory. Exiting." ; exit 1; }
chown -R "$current_user":"$current_user" /srv/festivals-identity-server
chmod -R 761 /srv/festivals-identity-server

# Download the backup script
#
echo "Downloading database backup script"
curl --progress-bar -L -o backup.sh https://raw.githubusercontent.com/Festivals-App/festivals-identity-server/main/operation/backup.sh
chown -R "$current_user":"$current_user" /srv/festivals-identity-server/backups/backup.sh
chmod -R 761 /srv/festivals-identity-server/backups/backup.sh
chmod +x /srv/festivals-identity-server/backups/backup.sh

# Installing a cronjob to run the backup every day at 3 pm.
#
echo "Installing a cronjob to periodically run a backup"
sleep 1
echo "0 3 * * * $current_user /srv/festivals-identity-server/backups/backup.sh" | sudo tee -a /etc/cron.d/festivals_identity_server_backup

# Installing festivals-identity-server binary
#
echo "Installing festivals-identity-server using port 22580."
sleep 1

# Get system os
#
if [ "$(uname -s)" = "Darwin" ]; then
  os="darwin"
elif [ "$(uname -s)" = "Linux" ]; then
  os="linux"
else
  echo "System is not Darwin or Linux. Exiting."
  exit 1
fi

# Get systems cpu architecture
#
if [ "$(uname -m)" = "x86_64" ]; then
  arch="amd64"
elif [ "$(uname -m)" = "arm64" ]; then
  arch="arm64"
else
  echo "System is not x86_64 or arm64. Exiting."
  exit 1
fi

# Build url to latest binary for the given system
#
file_url="https://github.com/Festivals-App/festivals-identity-server/releases/latest/download/festivals-identity-server-$os-$arch.tar.gz"
echo "The system is $os on $arch."
sleep 1

# Install festivals-identity-server to /usr/local/bin/festivals-identity-server. TODO: Maybe just link to /usr/local/bin?
#
echo "Downloading newest festivals-identity-server binary release..."
curl -L "$file_url" -o festivals-identity-server.tar.gz
tar -xf festivals-identity-server.tar.gz
mv festivals-identity-server /usr/local/bin/festivals-identity-server || { echo "Failed to install festivals-identity-server binary. Exiting." ; exit 1; }
echo "Installed the festivals-identity-server binary to '/usr/local/bin/festivals-identity-server'."
sleep 1

## Install server config file
#
mv config_template.toml /etc/festivals-identity-server.conf
echo "Moved default festivals-identity-server config to '/etc/festivals-identity-server.conf'."
sleep 1

## Prepare log directory
#
mkdir -p /var/log/festivals-identity-server || { echo "Failed to create log directory. Exiting." ; exit 1; }
chown "$WEB_USER":"$WEB_USER" /var/log/festivals-identity-server
echo "Create log directory at '/var/log/festivals-identity-server'."

## Prepare update workflow
#
mv update.sh /usr/local/festivals-identity-server/update.sh
cp /etc/sudoers /tmp/sudoers.bak
echo "$WEB_USER ALL = (ALL) NOPASSWD: /usr/local/festivals-identity-server/update.sh" >> /tmp/sudoers.bak
# Check syntax of the backup file to make sure it is correct.
visudo -cf /tmp/sudoers.bak
if [ $? -eq 0 ]; then
  # Replace the sudoers file with the new only if syntax is correct.
  sudo cp /tmp/sudoers.bak /etc/sudoers
else
  echo "Could not modify /etc/sudoers file. Please do this manually." ; exit 1;
fi

# Enable and configure the firewall.
#
if command -v ufw > /dev/null; then

  ufw allow 22580/tcp >/dev/null
  echo "Added festivals-identity-server to ufw using port 22580."
  sleep 1

elif ! [ "$(uname -s)" = "Darwin" ]; then
  echo "No firewall detected and not on macOS. Exiting."
  exit 1
fi

# Install systemd service
#
if command -v service > /dev/null; then

  if ! [ -f "/etc/systemd/system/festivals-identity-server.service" ]; then
    mv service_template.service /etc/systemd/system/festivals-identity-server.service
    echo "Created systemd service."
    sleep 1
  fi

  systemctl enable festivals-identity-server > /dev/null
  echo "Enabled systemd service."
  sleep 1

elif ! [ "$(uname -s)" = "Darwin" ]; then
  echo "Systemd is missing and not on macOS. Exiting."
  exit 1
fi

echo "Done!"
sleep 1

echo "You can start the server manually by running 'systemctl start festivals-identity-server' after you updated the configuration file at '/etc/festivals-identity-server.conf'"
sleep 1


# Cleanup
#
echo "Cleanup"
cd /usr/local/festivals-identity-server || exit
rm secure-mysql.sh
rm create_database.sql
sleep 1

echo "Done."
sleep 1
