#!/bin/bash
#
# Removes the firewall configuration, uninstalls go, git and the festivals-identity-server and stops and removes it as a service.
#
# (c)2020-2023 Simon Gaus
#

# Move to working directory
#
cd /usr/local || exit

# Stop the service
#
systemctl stop festivals-identity-server >/dev/null
echo "Stopped festivals-identity-server"
sleep 1

# Remove systemd configuration
#
systemctl disable festivals-identity-server >/dev/null
rm /etc/systemd/system/festivals-identity-server.service
echo "Removed systemd service"
sleep 1

# Remove the firewall configuration.
# This step is skipped under macOS.
#
if command -v ufw > /dev/null; then

  ufw delete allow 22580/tcp >/dev/null
  echo "Removed ufw configuration"
  sleep 1

elif ! [ "$(uname -s)" = "Darwin" ]; then
  echo "No firewall detected and not on macOS. Exiting."
  exit 1
fi

# Remove go
#
apt-get --purge remove golang -y
apt-get autoremove -y
echo "Removed go"
sleep 1

# Remove festivals-identity-server
#
rm /usr/local/bin/festivals-identity-server
rm /etc/festivals-identity-server.conf
rm -R /var/log/festivals-identity-server
rm -R /usr/local/festivals-identity-server
echo "Removed festivals-identity-server"
sleep 1

echo "Done"
