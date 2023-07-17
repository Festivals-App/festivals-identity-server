#!/bin/bash
#
# Updates the festivals-identity-server and restarts it.
#
# (c)2020-2023 Simon Gaus
#

# Move to working dir
#
mkdir /usr/local/festivals-identity-server/install || { echo "Failed to create working directory. Exiting." ; exit 1; }
cd /usr/local/festivals-identity-server/install || { echo "Failed to access working directory. Exiting." ; exit 1; }

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

# Updating festivals-identity-server to the newest binary release
#
echo "Downloading newest festivals-identity-server binary release..."
curl -L "$file_url" -o festivals-identity-server.tar.gz
tar -xf festivals-identity-server.tar.gz
mv festivals-identity-server /usr/local/bin/festivals-identity-server || { echo "Failed to install festivals-identity-server binary. Exiting." ; exit 1; }
echo "Updated festivals-identity-server binary."
sleep 1

# Removing unused files
#
echo "Cleanup..."
cd /usr/local/festivals-identity-server
rm -r /usr/local/festivals-identity-server/install
sleep 1

# Restarted the festivals-identity-server
#
systemctl restart festivals-identity-server
echo "Restarted the festivals-identity-server"
sleep 1

echo "Done!"
sleep 1