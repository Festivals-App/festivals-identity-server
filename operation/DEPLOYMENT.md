# Development Deployment on Proxmox

![Proxmox](https://img.shields.io/badge/proxmox-proxmox?style=for-the-badge&logo=proxmox&logoColor=%23E57000&labelColor=%232b2a33&color=%232b2a33)

## Prerequisites

This guide assumes you have already created a Virtual Machine (VM) by following the [VM deployment guide](https://github.com/Festivals-App/festivals-documentation/tree/main/deployment/vm-deployment).

Before starting the installation, ensure you have:

- Created and configured your VM
- SSH access secured and logged in as the new admin user
- Your server's IP address (use `ip a` to check)
- A server name matching the Common Name (CN) for your server certificate (e.g., `identity-0.festivalsapp.home` for a hostname `identity-0`).

I use the development wildcard server certificate (`CN=*festivalsapp.home`) for this guide.

  > **DON'T USE THIS IN PRODUCTION, SEE [festivals-pki](https://github.com/Festivals-App/festivals-pki) FOR SECURITY BEST PRACTICES FOR PRODUCTION**

## Installing the FestivalsApp Identity Server

Run the following commands to download and install the FestivalsApp Identity Server:

```bash
curl -o install.sh https://raw.githubusercontent.com/Festivals-App/festivals-identity-server/master/operation/install.sh
chmod +x install.sh
sudo ./install.sh <mysql_root_pw> <mysql_backup_pw> <read_write_pw>
```

## Copying mTLS Certificates to the VM

Copy the server mTLS certificates from your development machine to the VM:

```bash
scp /opt/homebrew/etc/pki/ca.crt <user>@<ip-address>:.
scp /opt/homebrew/etc/pki/issued/server.crt <user>@<ip-address>:.
scp /opt/homebrew/etc/pki/issued/server.key <user>@<ip-address>:.
```

Once copied, SSH into the VM and move them to the correct location:

```bash
sudo mv ca.crt /usr/local/festivals-identity-server/ca.crt
sudo mv server.crt /usr/local/festivals-identity-server/server.crt
sudo mv server.key /usr/local/festivals-identity-server/server.key
```

Set the correct permissions:

```bash
# Change owner to web user
sudo chown www-data /usr/local/festivals-identity-server/ca.crt
sudo chown www-data /usr/local/festivals-identity-server/server.crt
sudo chown www-data /usr/local/festivals-identity-server/server.key

# Set secure permissions
sudo chmod 640 /usr/local/festivals-identity-server/ca.crt
sudo chmod 640 /usr/local/festivals-identity-server/server.crt
sudo chmod 600 /usr/local/festivals-identity-server/server.key
```

## Configuring JWT Signing Keys

Convert the mTLS server certificate to use it as the authentication key:

  > **DON'T USE THIS IN PRODUCTION, SEE [festivals-pki](https://github.com/Festivals-App/festivals-pki) FOR SECURITY BEST PRACTICES FOR PRODUCTION**

```bash
openssl x509 -in /usr/local/festivals-identity-server/server.crt -out /usr/local/festivals-identity-server/authentication.publickey.pem -outform PEM
openssl rsa -in /usr/local/festivals-identity-server/server.key -text | sudo tee /usr/local/festivals-identity-server/authentication.privatekey.pem
```

Set the correct permissions:

```bash
# Change owner to web user
sudo chown www-data /usr/local/festivals-identity-server/authentication.publickey.pem
sudo chown www-data /usr/local/festivals-identity-server/authentication.privatekey.pem

# Set secure permissions
sudo chmod 640 /usr/local/festivals-identity-server/authentication.publickey.pem
sudo chmod 600 /usr/local/festivals-identity-server/authentication.privatekey.pem
```

## Configuring Root CA

Lets add the Festivals Development Root CA certificate to the system CA's.

```bash
sudo cp /usr/local/festivals-identity-server/ca.crt /usr/local/share/ca-certificates/festivals-dev-ca.crt
sudo update-ca-certificates
```

## Configuring the Festivals Identity Server

Open the configuration file:

```bash
sudo nano /etc/festivals-identity-server.conf
```

Set the IP address, servername and heartbeat endpoint:

```ini
[service]
bind-address = "<ip address>"
bin-host = "<servername>"
# For example: 
# bind-address = "192.168.8.185"
# bind-address = "identity-0.festivalsapp.home"

[heartbeat]
endpoint = "<discovery endpoint>"
#For example: endpoint = "https://discovery.festivalsapp.home/loversear"
```

## Setting Up DNS Resolution

For the services in the FestivalsApp backend to work correctly, proper DNS resolution is required. If you don’t have a DNS server, manually add the necessary entries to `/etc/hosts`:

```bash
sudo nano /etc/hosts
```

Add the following entries:

```ini
<ip address> <server name>
<gateway ip address> <discovery endpoint>

# For example: 
# 192.168.8.185 identity-0.festivalsapp.home
# 192.168.8.186 discovery.festivalsapp.home
```

## Final Steps

If you have not already set up the Festivals Gateway with its discovery service, you may see an error like:

```log
ERR Failed to send heartbeat
```

## **🚀 However, the identity service should now be running successfully. 🚀**

## Testing

Lets login as the default admin user and get the server info:

```bash
curl -H "Api-Key: TEST_API_KEY_001" -u "admin@email.com:we4711" --cert /opt/homebrew/etc/pki/issued/api-client.crt --key /opt/homebrew/etc/pki/private/api-client.key --cacert /opt/homebrew/etc/pki/ca.crt https://identity-0.festivalsapp.dev:22580/login
```

This should return a JWT Token `<Header.<Payload>.<Signatur>`

  > eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.
  > eyJVc2VySUQiOiIxIiwiVXNlclJvbGUiOjQyLCJVc2VyRmVzdGl2YWxzIjpbXSwiVXNlckFydGlzdHMiOltdLCJVc2VyTG9jYXRpb25zIjpbXSwiVXNlckV2ZW50cyI6W10sIlVzZXJMaW5rcyI6W10sIlVzZXJQbGFjZXMiOltdLCJVc2VySW1hZ2VzIjpbXSwiVXNlclRhZ3MiOltdLCJpc3MiOiJpZGVudGl0eS0wLmZlc3RpdmFsc2FwcC5ob21lIiwiZXhwIjoxNzQwMjMxMTQ4fQ.
  > geBq1pxEvqwjnKA5YTHQ8IjJc9mwkpsQIRy1kGc63oNXzyAhPrPJsepICXxr2yVmB0E8oDECXLn4Cy5V_p4UAduWXnc0r8S05ijV8NCfmsEcJg-RRO8POkGykiC2mrn-XR8Nf8OF0fLp7Mhsb0_aqBoTOLdtB9V7IV49-JjWwX5gHl3HuXGOOhe4n_epumc8w8yDxYakWeaBFtEtaRmhFXK_yttexYOLP6Z1BBTL005hBGhO58qVW0cfgf_t5VWBpUnz3zqdC-GFegItqJQbKZ2pmfmXNz_AoJf2JmPtCzpJ4lG6QeSslvdFuwaCdYpDQPOvnMSIORwrAq_FL2m7qw

Use this to make authorized calls to the FestivalsIdentityAPI:

```bash
curl -H "Api-Key: TEST_API_KEY_001" -H "Authorization: <JWT>" --cert /opt/homebrew/etc/pki/issued/api-client.crt --key /opt/homebrew/etc/pki/private/api-client.key --cacert /opt/homebrew/etc/pki/ca.crt https://identity-0.festivalsapp.home:22580/info
```
