# Development Deployment

This deployment guide explains how to deploy the FestivalsApp Identity Server
using certificates intended for development purposes.

## Prerequisites

This guide assumes you have already created a Virtual Machine (VM) by following the [VM deployment guide](https://github.com/Festivals-App/festivals-documentation/tree/main/deployment/vm-deployment).

Before starting the installation, ensure you have:

- Created and configured your VM
- SSH access secured and logged in as the admin user
- Your server's IP address (use `ip a` to check)
- A server name matching the Common Name (CN) for your mTLS certificate
  (e.g., `identity-0.festivalsapp.home` for a hostname `identity-0`).

I use the development wildcard server certificate (`CN=*festivalsapp.home`) for this guide.

  > **DON'T USE THIS IN PRODUCTION, SEE [festivals-pki](https://github.com/Festivals-App/festivals-pki) FOR SECURITY BEST PRACTICES FOR PRODUCTION**

## 1. Installing the FestivalsApp Identity Server

Run the following commands to install the FestivalsApp Identity Server:

```bash
curl -o install.sh https://raw.githubusercontent.com/Festivals-App/festivals-identity-server/master/operation/install.sh
chmod +x install.sh
sudo ./install.sh <mysql_root_pw> <mysql_backup_pw> <read_write_pw>
```

The config file is located at:

  > `/etc/festivals-identity-server.conf`.

You also need to provide certificates in the right format and location:

  > Root CA certificate           `/usr/local/festivals-identity-server/ca.crt`  
  > Server certificate            `/usr/local/festivals-identity-server/server.crt`  
  > Server key                    `/usr/local/festivals-identity-server/server.key`  
  > Authentication certificate    `/usr/local/festivals-identity-server/authentication.publickey.pem`  
  > Authentication key            `/usr/local/festivals-identity-server/authentication.privatekey.pem`  

Where the root CA certificate is required to validate incoming requests, the server certificate and key
is required to make outgoing connections via mTLS and the authentication certificate and key is required
to create and validate JSON Web Token ([JWT](https://de.wikipedia.org/wiki/JSON_Web_Token)) for the authentication API.
For instructions on how to manage and create the certificates see the [festivals-pki](https://github.com/Festivals-App/festivals-pki) repository.

## 2. Copying mTLS Certificates to the VM

Copy the server mTLS certificates from your development machine to the VM:

```bash
scp /opt/homebrew/etc/pki/ca.crt <user>@<ip-address>:.
scp /opt/homebrew/etc/pki/issued/server.crt <user>@<ip-address>:.
scp /opt/homebrew/etc/pki/private/server.key <user>@<ip-address>:.
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
sudo chown www-data:www-data /usr/local/festivals-identity-server/ca.crt
sudo chown www-data:www-data /usr/local/festivals-identity-server/server.crt
sudo chown www-data:www-data /usr/local/festivals-identity-server/server.key
# Set secure permissions
sudo chmod 640 /usr/local/festivals-identity-server/ca.crt
sudo chmod 640 /usr/local/festivals-identity-server/server.crt
sudo chmod 600 /usr/local/festivals-identity-server/server.key
```

## 3. Configuring the  JWT Signing Keys

Convert the mTLS server certificate to use it as the authentication key:

  > **DON'T USE THIS IN PRODUCTION, SEE [festivals-pki](https://github.com/Festivals-App/festivals-pki) FOR SECURITY BEST PRACTICES FOR PRODUCTION**

```bash
sudo openssl x509 -in /usr/local/festivals-identity-server/server.crt -out /usr/local/festivals-identity-server/authentication.publickey.pem -outform PEM
sudo openssl rsa -in /usr/local/festivals-identity-server/server.key -text | sudo tee /usr/local/festivals-identity-server/authentication.privatekey.pem
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

## 4. Configuring the Festivals Identity Server

Open the configuration file:

```bash
sudo nano /etc/festivals-identity-server.conf
```

Set the server name and heartbeat endpoint:

```ini
[service]
bind-host = "<server name>"
# For example: 
# bind-host = "identity-0.festivalsapp.home"

[database]
password = "<festivals.identity.writer password>"
# For example: 
# password = "we4711"

[heartbeat]
endpoint = "<discovery endpoint>"
#For example: endpoint = "https://discovery.festivalsapp.home/loversear"
```

## Optional: Restore database backup

Copy the backup from the old server and copy to the new one

```bash
scp <user>@<host>:/srv/festivals-identity-server/backups/<date>/festivals_identity_database-<datetime>.gz ~/Desktop
scp ~/Desktop/festivals_identity_database-<datetime>.gz <user>@<host>:.
```

Now decompress and import the backuped database into mysql

```bash
gzip -d festivals_identity_database-<datetime>.gz
sudo mysql -uroot -p < festivals_identity_database-<datetime>
```

And now let's start the service:

```bash
sudo systemctl start festivals-identity-server
```

## **🚀 The identity service should now be running successfully. 🚀**

  > You might encounter an `ERR Failed to send heartbeat` error if the discovery service is not yet available.
    However, the service should function correctly.

### Optional: Setting Up DNS Resolution  

For the services in the FestivalsApp backend to function correctly, proper DNS resolution is required.
This is because mTLS is configured to validate the client’s certificate identity based on its DNS hostname.  

If you don’t have a DNS server to manage DNS for your development VMs, you can manually configure DNS resolution
by adding the necessary entries to each server’s `/etc/hosts` file:  

```bash
sudo nano /etc/hosts
```

Add the following entries:  

```ini
<IP address> <server name>  
<Gateway IP address> <discovery endpoint>  
# ...

# Example:  
# 192.168.8.185 identity-0.festivalsapp.home
# 192.168.8.186 discovery.festivalsapp.home
# ...
```

**Keep in mind that you will need to update each machine’s `hosts` file whenever you add a new VM or if any IP addresses change.**

### Testing

Lets login as the default admin user and get the server info:

```bash
curl -H "Api-Key: TEST_API_KEY_001" -u "admin@email.com:we4711" --cert /opt/homebrew/etc/pki/issued/api-client.crt --key /opt/homebrew/etc/pki/private/api-client.key --cacert /opt/homebrew/etc/pki/ca.crt https://identity-0.festivalsapp.home:22580/users/login
```

This should return a JWT Token `<Header.<Payload>.<Signatur>`, use this to make authorized calls to the FestivalsIdentityAPI:

```bash
curl -H "Api-Key: TEST_API_KEY_001" -H "Authorization: Bearer <JWT>" --cert /opt/homebrew/etc/pki/issued/api-client.crt --key /opt/homebrew/etc/pki/private/api-client.key --cacert /opt/homebrew/etc/pki/ca.crt https://identity-0.festivalsapp.home:22580/info
```
