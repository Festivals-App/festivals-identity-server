#!/bin/bash
#
# install.sh - FestivalsApp Identity Server Installer Scrip
#
# Enables the firewall, installs the latest version of MySQL, starts it as a service, 
# configures it as the database server for FestivalsApp Identity Server, and sets up backup routines.
#
# (c)2020-2025 Simon Gaus
#

# ─────────────────────────────────────────────────────────────────────────────
# 🛑 Check if all parameters are supplied
# ─────────────────────────────────────────────────────────────────────────────
if [ $# -ne 3 ]; then
    echo -e "\n\033[1;31m🚨  ERROR: Missing parameters!\033[0m"
    echo -e "\n\033[1;34m🔹  USAGE:\033[0m sudo ./install.sh \033[1;32m<mysql_root_pw> <mysql_backup_pw> <read_write_pw>\033[0m"
    echo -e "\n\033[1;34m⚠️  REQUIREMENTS:\033[0m Run as \033[1;33mroot\033[0m or with \033[1;33msudo\033[0m."
    echo -e "\n\033[1;31m❌  Exiting.\033[0m\n"
    exit 1
fi

# ─────────────────────────────────────────────────────────────────────────────
# 🎯 Store parameters in variables
# ─────────────────────────────────────────────────────────────────────────────
root_password="$1"
backup_password="$2"
read_write_password="$3"

# ─────────────────────────────────────────────────────────────────────────────
# 🔍 Detect Web Server User
# ─────────────────────────────────────────────────────────────────────────────
WEB_USER="www-data"
if ! id -u "$WEB_USER" &>/dev/null; then
    WEB_USER="www"
    if ! id -u "$WEB_USER" &>/dev/null; then
        echo -e "\n\033[1;31m❌  ERROR: Web server user not found! Exiting.\033[0m\n"
        exit 1
    fi
fi

echo -e "\n👤  Web server user detected: \e[1;34m$WEB_USER\e[0m"
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 📁 Setup Working Directory
# ─────────────────────────────────────────────────────────────────────────────
WORK_DIR="/usr/local/festivals-identity-server/install"
mkdir -p "$WORK_DIR" && cd "$WORK_DIR" || { echo -e "\n\033[1;31m❌  ERROR: Failed to create/access working directory!\033[0m\n"; exit 1; }

echo -e "\n📂  Working directory set to \e[1;34m$WORK_DIR\e[0m\n"
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 📦 Install MySQL Server
# ─────────────────────────────────────────────────────────────────────────────
echo -e "\n\n\n🚀  Installing MySQL server..."
apt-get install mysql-server -y > /dev/null 2>&1
echo -e "\n✅  MySQL server installed.\n"
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🔄 Enable & Start MySQL Service
# ─────────────────────────────────────────────────────────────────────────────
echo -e "\n\n\n▶️  Enabling and starting MySQL service..."
systemctl enable mysql &>/dev/null && systemctl start mysql &>/dev/null

echo -e "\n✅  MySQL service is up and running.\n"
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🔐 Install MySQL Credential File
# ─────────────────────────────────────────────────────────────────────────────

echo -e "\n\n\n📂  Installing MySQL credential file to project directory..."
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

if [ -f "$credentialsFile" ]; then
    echo -e "\n✅  MySQL credential file successfully created at \e[1;34m$credentialsFile\e[0m\n"
else
    echo -e "\n🚨  ERROR: Failed to create MySQL credential file. Exiting.\n"
    exit 1
fi
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🔑 Secure MySQL & Create Users
# ─────────────────────────────────────────────────────────────────────────────
echo -e "\n\n\n🔐  Securing MySQL..."
curl --progress-bar -L -o secure-mysql.sh https://raw.githubusercontent.com/Festivals-App/festivals-identity-server/master/operation/secure-mysql.sh
chmod +x secure-mysql.sh
./secure-mysql.sh "$root_password"
echo -e "\n✅  MySQL secured.\n"
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🗄️  Setup Database & Users
# ─────────────────────────────────────────────────────────────────────────────
echo -e "\n\n\n📊  Creating database and users..."
curl --progress-bar -L -o create_database.sql https://raw.githubusercontent.com/Festivals-App/festivals-identity-server/master/database/create_database.sql
mysql -e "source $WORK_DIR/create_database.sql"
mysql -e "CREATE USER 'festivals.identity.writer'@'localhost' IDENTIFIED BY '$read_write_password';"
mysql -e "GRANT SELECT, INSERT, UPDATE, DELETE ON festivals_identity_database.* TO 'festivals.identity.writer'@'localhost';"
mysql -e "CREATE USER 'festivals.identity.backup'@'localhost' IDENTIFIED BY '$backup_password';"
mysql -e "GRANT ALL PRIVILEGES ON festivals_identity_database.* TO 'festivals.identity.backup'@'localhost';"
mysql -e "FLUSH PRIVILEGES;"
echo -e "\n✅  Database and users created.\n"
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 📂 Setup Backup Directory
# ─────────────────────────────────────────────────────────────────────────────
echo -e "\n\n\n💾  Setting up backup directory..."
mkdir -p /srv/festivals-identity-server/backups
curl --progress-bar -L -o /srv/festivals-identity-server/backups/backup.sh https://raw.githubusercontent.com/Festivals-App/festivals-identity-server/main/operation/backup.sh
chmod +x /srv/festivals-identity-server/backups/backup.sh
echo -e "\n✅  Backup directory and script configured.\n"
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# ⏳ Install Cronjob for Daily Backup at 3 AM
# ─────────────────────────────────────────────────────────────────────────────

echo -e "\n\n\n🕒  Installing a cronjob to periodically run a backup..."
sleep 1
echo -e "0 3 * * * $WEB_USER /srv/festivals-identity-server/backups/backup.sh" | tee -a /etc/cron.d/festivals_identity_server_backup > /dev/null
echo -e "\n✅  Cronjob successfully installed! Backup will run daily at \e[1;34m3 AM\e[0m\n"

# ─────────────────────────────────────────────────────────────────────────────
# 🖥  Detect System OS and Architecture
# ─────────────────────────────────────────────────────────────────────────────

echo -e "\n\n\n🔍  Detecting system OS and architecture..."
sleep 1

if [ "$(uname -s)" = "Darwin" ]; then
    os="darwin"
elif [ "$(uname -s)" = "Linux" ]; then
    os="linux"
else
    echo -e "\n🚨  ERROR: Unsupported OS. Exiting.\n"
    exit 1
fi

if [ "$(uname -m)" = "x86_64" ]; then
    arch="amd64"
elif [ "$(uname -m)" = "arm64" ]; then
    arch="arm64"
else
    echo -e "\n🚨  ERROR: Unsupported CPU architecture. Exiting.\n"
    exit 1
fi

echo -e "\n✅  Detected OS: \e[1;34m$os\e[0m, Architecture: \e[1;34m$arch\e[0m."
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 📦 Install FestivalsApp Identity Server
# ─────────────────────────────────────────────────────────────────────────────

file_url="https://github.com/Festivals-App/festivals-identity-server/releases/latest/download/festivals-identity-server-$os-$arch.tar.gz"

echo -e "\n📥  Downloading latest FestivalsApp Identity Server binary..."
curl --progress-bar -L "$file_url" -o festivals-identity-server.tar.gz

echo -e "\n📦  Extracting binary..."
tar -xf festivals-identity-server.tar.gz

mv festivals-identity-server /usr/local/bin/festivals-identity-server || {
    echo -e "\n🚨  ERROR: Failed to install Festivals Identity Server binary. Exiting.\n"
    exit 1
}

echo -e "\n✅  Installed FestivalsApp Identity Server to \e[1;34m/usr/local/bin/festivals-identity-server\e[0m.\n"
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🛠  Install Server Configuration File
# ─────────────────────────────────────────────────────────────────────────────

echo -e "\n\n\n📂  Moving default configuration file..."
mv config_template.toml /etc/festivals-identity-server.conf

if [ -f "/etc/festivals-identity-server.conf" ]; then
    echo -e "\n✅  Configuration file moved to \e[1;34m/etc/festivals-identity-server.conf\e[0m.\n"
else
    echo -e "\n🚨  ERROR: Failed to move configuration file. Exiting.\n"
    exit 1
fi
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 📂  Prepare Log Directory
# ─────────────────────────────────────────────────────────────────────────────

echo -e "\n\n\n📁  Creating log directory..."
mkdir -p /var/log/festivals-identity-server || {
    echo -e "\n🚨  ERROR: Failed to create log directory. Exiting.\n"
    exit 1
}

echo -e "\n✅  Log directory created at \e[1;34m/var/log/festivals-identity-server\e[0m.\n"
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🔄 Prepare Remote Update Workflow
# ─────────────────────────────────────────────────────────────────────────────

echo -e "\n\n\n⚙️  Preparing remote update workflow..."
sleep 1

mv update.sh /usr/local/festivals-identity-server/update.sh
chmod +x /usr/local/festivals-identity-server/update.sh

cp /etc/sudoers /tmp/sudoers.bak
echo "$WEB_USER ALL = (ALL) NOPASSWD: /usr/local/festivals-identity-server/update.sh" >> /tmp/sudoers.bak

# Validate and replace sudoers file if syntax is correct
if visudo -cf /tmp/sudoers.bak &>/dev/null; then
    sudo cp /tmp/sudoers.bak /etc/sudoers
    echo -e "\n✅  Updated sudoers file successfully."
else
    echo -e "\n🚨  ERROR: Could not modify /etc/sudoers file. Please do this manually. Exiting.\n"
    exit 1
fi
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🔥 Enable and Configure Firewall
# ─────────────────────────────────────────────────────────────────────────────

if command -v ufw > /dev/null; then
    echo -e "\n\n\n🚀  Configuring UFW firewall..."
    mv ufw_app_profile /etc/ufw/applications.d/festivals-identity-server
    ufw allow festivals-identity-server > /dev/null
    echo -e "\n✅  Added festivals-identity-server to UFW with port 22580."
    sleep 1
elif ! [ "$(uname -s)" = "Darwin" ]; then
    echo -e "\n🚨  ERROR: No firewall detected and not on macOS. Exiting.\n"
    exit 1
fi

# ─────────────────────────────────────────────────────────────────────────────
# ⚙️  Install Systemd Service
# ─────────────────────────────────────────────────────────────────────────────

if command -v service > /dev/null; then
    echo -e "\n\n\n🚀  Configuring systemd service..."
    if ! [ -f "/etc/systemd/system/festivals-identity-server.service" ]; then
        mv service_template.service /etc/systemd/system/festivals-identity-server.service
        echo -e "\n✅  Created systemd service configuration."
        sleep 1
    fi
    systemctl enable festivals-identity-server > /dev/null
    echo -e "\n✅  Enabled systemd service for Festivals Identity Server."
    sleep 1
elif ! [ "$(uname -s)" = "Darwin" ]; then
    echo -e "\n🚨  ERROR: Systemd is missing and not on macOS. Exiting.\n"
    exit 1
fi

# ─────────────────────────────────────────────────────────────────────────────
# 🔑 Set Appropriate Permissions
# ─────────────────────────────────────────────────────────────────────────────

chown -R "$WEB_USER":"$WEB_USER" /usr/local/festivals-identity-server
chown -R "$WEB_USER":"$WEB_USER" /var/log/festivals-identity-server
chown -R "$WEB_USER":"$WEB_USER" /srv/festivals-identity-server
chown "$WEB_USER":"$WEB_USER" /etc/festivals-identity-server.conf

# ─────────────────────────────────────────────────────────────────────────────
# 🧹 Cleanup Installation Files
# ─────────────────────────────────────────────────────────────────────────────

echo -e "\n🧹  Cleaning up installation files..."
cd /usr/local/festivals-identity-server || exit
rm -rf /usr/local/festivals-identity-server/install
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🎉 Final Message
# ─────────────────────────────────────────────────────────────────────────────

echo -e "\n\n\n\n\033[1;32m══════════════════════════════════════════════════════════════════════════\033[0m"
echo -e "\033[1;32m✅  INSTALLATION COMPLETE! 🚀\033[0m"
echo -e "\033[1;32m══════════════════════════════════════════════════════════════════════════\033[0m"
sleep 1

echo -e "\n🔹 \033[1;34mTo start the server manually, run:\033[0m"
echo -e "\n   \033[1;32msudo systemctl start festivals-identity-server\033[0m"

echo -e "\n📂 \033[1;34mBefore starting, update the configuration file at:\033[0m"
echo -e "\n   \033[1;34m/etc/festivals-identity-server.conf\033[0m"

echo -e "\n\033[1;32m══════════════════════════════════════════════════════════════════════════\033[0m\n"
sleep 1
