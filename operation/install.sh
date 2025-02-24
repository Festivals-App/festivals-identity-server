#!/bin/bash
#
# install.sh - FestivalsApp Identity Server Installer Script
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
    echo -e "\033[1;34m🔹  USAGE:\033[0m sudo ./install.sh \033[1;32m<mysql_root_pw> <mysql_backup_pw> <read_write_pw>\033[0m"
    echo -e "\033[1;31m❌  Exiting.\033[0m\n"
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

# ─────────────────────────────────────────────────────────────────────────────
# 📁 Setup Working Directory
# ─────────────────────────────────────────────────────────────────────────────
WORK_DIR="/usr/local/festivals-identity-server/install"
mkdir -p "$WORK_DIR" && cd "$WORK_DIR" || { echo -e "\n\033[1;31m❌  ERROR: Failed to create/access working directory!\033[0m\n"; exit 1; }
echo -e "\n📂  Working directory set to \e[1;34m$WORK_DIR\e[0m"
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🖥  Detect System OS and Architecture
# ─────────────────────────────────────────────────────────────────────────────
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

# ─────────────────────────────────────────────────────────────────────────────
# 📦 Download latest release
# ─────────────────────────────────────────────────────────────────────────────
file_url="https://github.com/Festivals-App/festivals-identity-server/releases/latest/download/festivals-identity-server-$os-$arch.tar.gz"
echo -e "\n📥  Downloading latest FestivalsApp Identity Server release..."
curl --progress-bar -L "$file_url" -o festivals-identity-server.tar.gz
echo -e "📦  Extracting archive..."
tar -xf festivals-identity-server.tar.gz

# ─────────────────────────────────────────────────────────────────────────────
# 📦 Install & Enable & Start MySQL Server
# ─────────────────────────────────────────────────────────────────────────────
echo -e "\n🗂️  Installing MySQL server..."
apt-get install mysql-server -y > /dev/null 2>&1
systemctl enable mysql &>/dev/null && systemctl start mysql &>/dev/null
echo -e "✅  MySQL service is up and running."
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🔐 Install MySQL Backup Credential File
# ─────────────────────────────────────────────────────────────────────────────
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
    echo -e "✅  MySQL backup credential file successfully created at \e[1;34m$credentialsFile\e[0m"
else
    echo -e "🚨  ERROR: Failed to create MySQL credential file. Exiting.\n"
    exit 1
fi
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🔑 Secure MySQL
# ─────────────────────────────────────────────────────────────────────────────
chmod +x secure-mysql.sh
./secure-mysql.sh "$root_password"
echo -e "✅  MySQL security script executed."
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🗄️  Setup Database & Users
# ─────────────────────────────────────────────────────────────────────────────
mysql -e "source $WORK_DIR/create_database.sql"
mysql -e "CREATE USER 'festivals.identity.writer'@'localhost' IDENTIFIED BY '$read_write_password';"
mysql -e "GRANT SELECT, INSERT, UPDATE, DELETE ON festivals_identity_database.* TO 'festivals.identity.writer'@'localhost';"
mysql -e "CREATE USER 'festivals.identity.backup'@'localhost' IDENTIFIED BY '$backup_password';"
mysql -e "GRANT ALL PRIVILEGES ON festivals_identity_database.* TO 'festivals.identity.backup'@'localhost';"
mysql -e "FLUSH PRIVILEGES;"
echo -e "✅  Database and users created."
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 📂 Setup Database Backup Directory
# ─────────────────────────────────────────────────────────────────────────────
mkdir -p /srv/festivals-identity-server/backups
mv backup.sh /srv/festivals-identity-server/backups/backup.sh
chmod +x /srv/festivals-identity-server/backups/backup.sh
echo -e "✅  Database backup directory and script configured."
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# ⏳ Install Cronjob for Daily Backup at 3 AM
# ─────────────────────────────────────────────────────────────────────────────
echo -e "0 3 * * * $WEB_USER /srv/festivals-identity-server/backups/backup.sh" | tee -a /etc/cron.d/festivals_identity_server_backup > /dev/null
echo -e "✅  Cronjob installed! Backup will run daily at \e[1;34m3 AM\e[0m"
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 📦 Install FestivalsApp Identity Server
# ─────────────────────────────────────────────────────────────────────────────
echo -e "\n📥  Installing latest FestivalsApp Identity Server binary..."
mv festivals-identity-server /usr/local/bin/festivals-identity-server || {
    echo -e "\n🚨  ERROR: Failed to install Festivals Identity Server binary. Exiting.\n"
    exit 1
}
echo -e "✅  Installed FestivalsApp Identity Server to \e[1;34m/usr/local/bin/festivals-identity-server\e[0m."
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🛠  Install Server Configuration File
# ─────────────────────────────────────────────────────────────────────────────
mv config_template.toml /etc/festivals-identity-server.conf
if [ -f "/etc/festivals-identity-server.conf" ]; then
    echo -e "✅  Configuration file moved to \e[1;34m/etc/festivals-identity-server.conf\e[0m."
else
    echo -e "\n🚨  ERROR: Failed to move configuration file. Exiting.\n"
    exit 1
fi
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 📂  Prepare Log Directory
# ─────────────────────────────────────────────────────────────────────────────
mkdir -p /var/log/festivals-identity-server || {
    echo -e "\n🚨  ERROR: Failed to create log directory. Exiting.\n"
    exit 1
}
echo -e "✅  Log directory created at \e[1;34m/var/log/festivals-identity-server\e[0m."
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🔄 Prepare Remote Update Workflow
# ─────────────────────────────────────────────────────────────────────────────
mv update.sh /usr/local/festivals-identity-server/update.sh
chmod +x /usr/local/festivals-identity-server/update.sh
cp /etc/sudoers /tmp/sudoers.bak
echo "$WEB_USER ALL = (ALL) NOPASSWD: /usr/local/festivals-identity-server/update.sh" >> /tmp/sudoers.bak
# Validate and replace sudoers file if syntax is correct
if visudo -cf /tmp/sudoers.bak &>/dev/null; then
    sudo cp /tmp/sudoers.bak /etc/sudoers
    echo -e "✅  Prepared remote update workflow."
else
    echo -e "\n🚨  ERROR: Could not modify /etc/sudoers file. Please do this manually. Exiting.\n"
    exit 1
fi
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🔥 Enable and Configure Firewall
# ─────────────────────────────────────────────────────────────────────────────
if command -v ufw > /dev/null; then
    echo -e "\n🔥  Configuring UFW firewall..."
    mv ufw_app_profile /etc/ufw/applications.d/festivals-identity-server
    ufw allow festivals-identity-server > /dev/null
    echo -e "✅  Added festivals-identity-server to UFW with port 22580."
    sleep 1
elif ! [ "$(uname -s)" = "Darwin" ]; then
    echo -e "\n🚨  ERROR: No firewall detected and not on macOS. Exiting.\n"
    exit 1
fi

# ─────────────────────────────────────────────────────────────────────────────
# ⚙️  Install Systemd Service
# ─────────────────────────────────────────────────────────────────────────────
if command -v service > /dev/null; then
    echo -e "\n🚀  Configuring systemd service..."
    if ! [ -f "/etc/systemd/system/festivals-identity-server.service" ]; then
        mv service_template.service /etc/systemd/system/festivals-identity-server.service
        echo -e "✅  Created systemd service configuration."
        sleep 1
    fi
    systemctl enable festivals-identity-server > /dev/null
    echo -e "✅  Enabled systemd service for Festivals Identity Server."
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
echo -e "\n🔐  Set Appropriate Permissions."
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🧹 Cleanup Installation Files
# ─────────────────────────────────────────────────────────────────────────────
echo -e "🧹  Cleaning up installation files..."
cd /usr/local/festivals-identity-server || exit
rm -rf /usr/local/festivals-identity-server/install
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🎉 COMPLETE Message
# ─────────────────────────────────────────────────────────────────────────────
echo -e "\n\033[1;32m══════════════════════════════════════════════════════════════════════════\033[0m"
echo -e "\033[1;32m✅  INSTALLATION COMPLETE! 🚀\033[0m"
echo -e "\033[1;32m══════════════════════════════════════════════════════════════════════════\033[0m"
echo -e "\n📂 \033[1;34mBefore starting, you need to:\033[0m"
echo -e "\n   \033[1;34m1. Configure the mTLS certificates.\033[0m"
echo -e "   \033[1;34m2. Configure the JWT signing keys.\033[0m"
echo -e "   \033[1;34m3. Configuring the FestivlasApp Root CA.\033[0m"
echo -e "   \033[1;34m4. Update the configuration file at:\033[0m"
echo -e "\n   \033[1;32m    /etc/festivals-identity-server.conf\033[0m"
echo -e "\n🔹 \033[1;34mThen start the server manually:\033[0m"
echo -e "\n   \033[1;32m    sudo systemctl start festivals-identity-server\033[0m"
echo -e "\n\033[1;32m══════════════════════════════════════════════════════════════════════════\033[0m\n"
sleep 1
