# Public-Key-Infrastruktur

To secure communication between microservices the `FestivalsAPI` uses [mTLS](https://www.cloudflare.com/learning/access-management/what-is-mutual-tls/) with self signed certificates. 
First we need to create a certificate authority (CA) to issue certificates, then we create a certificate for each service and each administrative user.
You can read more about the exact procedures in the [Certification Practice Statements](CERTIFICATIONPRACTICE.md) document.

## Development

This project uses [easy-rsa](https://github.com/OpenVPN/easy-rsa), which is maintained by the wonderfull community of [OpenVPN](https://openvpn.net/community/), to build and manage the FestivalsApp Root CA.
Even tho the details are quite complex, at the basis of a certificate authority stands a single root certificate. The root certificate is self-signed, meaning that we create it ourself.
The idea is that every party that is communicating with each other needs a certificate signes with 

1. First we need to install `easy-rsa` and create the FestivalsApp Root CA.
```bash
# installing the easy-rsa
brew install easy-rsa

# init the pki and create our CA
easyrsa init-pki
easyrsa build-ca
```
This will create all neccessary file at `/opt/homebrew/etc/pki`

2. To create a certificate/key pair for inter service communication we first create a certificate request with the name of the service node and then sign the request.
```bash
# create signing request
easyrsa gen-req <UNIQUE_SERVER_NAME> nopass
# Enter <UNIQUE_SERVER_DOMAIN_NAME>
Common Name (eg: your user, host, or server name): <UNIQUE_SERVER_DOMAIN_NAME>
# sign the request
easyrsa sign-req serverClient <UNIQUE_SERVER_NAME>
```

3. Copy the certificate/key pair to server
```bash
scp <path/to/cert/key> <user>@<server>:/home/<user>
```

## Local Development
If you want to test on your local machine

sudo nano /etc/hosts
```bash
# local development on this machine
127.0.0.1       festivalsapp.dev
127.0.0.1       www.festivalsapp.dev
127.0.0.1       website.festivalsapp.dev
127.0.0.1       gateway.festivalsapp.dev
127.0.0.1       discovery.festivalsapp.dev
127.0.0.1       api.festivalsapp.dev
127.0.0.1       files.festivalsapp.dev
```

## Development on a test server
If you have an development server in your private network

```bash
# local development server for festivalsapp
<identity server ip address>	identity-0.festivalsapp.home
<festivals server 0 ip address>	festivals-0.festivalsapp.home
<festivals server 1 ip address>	festivals-1.festivalsapp.home
<database ip address>	        database-0.festivalsapp.home
<fileserver ip address>	        fileserver-0.festivalsapp.home
<website ip address>	        website-0.festivalsapp.home

<gateway ip address>            festivalsapp.home
<gateway ip address>            www.festivalsapp.home
<gateway ip address>            website.festivalsapp.home
<gateway ip address>            gateway.festivalsapp.home
<gateway ip address>            discovery.festivalsapp.home
<gateway ip address>            api.festivalsapp.home
<gateway ip address>            files.festivalsapp.home
```

# Deployment

Add the FestivalsaApp Root CA certificate to the trusted root certificates:

### macOS
```bash

sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain ~/new-root-certificate.crt
```
### Linux (Ubuntu)
```bash
sudo cp new-root-certificate.crt /usr/local/share/ca-certificates/new-root-certificate.crt
sudo update-ca-certificates
```







https://docs.bigchaindb.com/projects/server/en/v1.1.0/production-deployment-template/easy-rsa.html#how-to-install-configure-easy-rsa
openssl rsa -in /opt/homebrew/etc/pki/private/gateway-server.key -out /opt/homebrew/etc/pki/private/gateway-server-unencrypted.key