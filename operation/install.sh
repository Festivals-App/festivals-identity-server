#!/bin/bash
#
# Enables the firewall, installs the newest festivals-identity-server and starts it as a service.
#
# (c)2020-2023 Simon Gaus
#

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

# Move to working dir
#
mkdir -p /usr/local/festivals-identity-server/install || { echo "Failed to create working directory. Exiting." ; exit 1; }
cd /usr/local/festivals-identity-server/install || { echo "Failed to access working directory. Exiting." ; exit 1; }

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
mkdir /var/log/festivals-identity-server || { echo "Failed to create log directory. Exiting." ; exit 1; }
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

# Remving unused files
#
echo "Cleanup..."
cd /usr/local/festivals-identity-server || exit
rm -R /usr/local/festivals-identity-server/install
sleep 1

echo "Done!"
sleep 1

echo "You can start the server manually by running 'systemctl start festivals-identity-server' after you updated the configuration file at '/etc/festivals-identity-server.conf'"
sleep 1
