<p align="center">
   <a href="https://github.com/festivals-app/festivals-identity-server/commits/" title="Last Commit"><img src="https://img.shields.io/github/last-commit/festivals-app/festivals-identity-server?style=flat"></a>
   <a href="https://github.com/festivals-app/festivals-identity-server/issues" title="Open Issues"><img src="https://img.shields.io/github/issues/festivals-app/festivals-identity-server?style=flat"></a>
   <a href="./LICENSE" title="License"><img src="https://img.shields.io/github/license/festivals-app/festivals-identity-server.svg"></a>
</p>

<h1 align="center">
  <br/><br/>
    Festivals App Identity Server
  <br/><br/>
</h1>

A lightweight go server app providing a RESTful API, called FestivalsIdentityAPI. The FestivalsIdentityAPI exposes all authorization and authentication functions needed by the FestivalsApp components.

![Figure 1: Architecture Overview Highlighted](https://github.com/Festivals-App/festivals-documentation/blob/main/images/architecture/architecture_overview_identity.svg "Figure 1: Architecture Overview Highlighted")

<hr/>
<p align="center">
  <a href="#development">Development</a> •
  <a href="#deployment">Deployment</a> •
  <a href="#festivalsidentityapi">FestivalsIdentityAPI</a> •
  <a href="#architecture">Architecture</a> •
  <a href="#engage">Engage</a>
</p>
<hr/>

## Development

1. Make server run ✅
2. Make server serves tls ✅
3. Make other server serve tls ✅
4. implement jwt to validate requests to other partys (especialy for admin requests) ✅

### Requirements

- [Golang](https://go.dev/) Version 1.23.0+
- [Visual Studio Code](https://code.visualstudio.com/download) 1.92.2+
  - Plugin recommendations are managed via [workspace recommendations](https://code.visualstudio.com/docs/editor/extension-marketplace#_recommended-extensions).
- [Bash script](https://en.wikipedia.org/wiki/Bash_(Unix_shell)) friendly environment

## Deployment

The Go binaries are able to run without system dependencies so there are not many requirements for the system to run the festivals-identity-server binary.
The config file needs to be placed at `/etc/festivals-identity-server.conf` or the template config file needs to be present in the directory the binary runs in.

You also need to provide certificates in the right format and location:

- The default path to the root CA certificate is          `/usr/local/festivals-identity-server/ca.crt`
- The default path to the server certificate is           `/usr/local/festivals-identity-server/server.crt`
- The default path to the corresponding key is            `/usr/local/festivals-identity-server/server.key`
- The default path to the authentication certificate is   `/usr/local/festivals-identity-server/authentication.pem`
- The default path to the corresponding key is            `/usr/local/festivals-identity-server/authentication-key.pem`

Where the root CA certificate is required to validate incoming requests, the server certificate and key is required to make outgoing connections
and the authentication certificate and key is required to create and validate JSON Web Token ([JWT](https://de.wikipedia.org/wiki/JSON_Web_Token)) for the authentication API.
For instructions on how to manage and create the certificates see the [festivals-pki](https://github.com/Festivals-App/festivals-pki) repository.

### VM

```bash
#Installing
curl -o install.sh https://raw.githubusercontent.com/Festivals-App/festivals-identity-server/master/operation/install.sh
chmod +x install.sh
sudo ./install.sh <mysql_root_pw> <mysql_backup_pw> <read_write_pw>
sudo nano /etc/mysql/mysql.conf.d/mysqld.cnf          // edit bind-address=<private-ip>

#Updating
curl -o update.sh https://raw.githubusercontent.com/Festivals-App/festivals-identity-server/master/operation/update.sh
chmod +x update.sh
sudo ./update.sh

#To see if the server is running use:
sudo systemctl status festivals-identity-server
```

#### Build and run using make

```bash
make build
make run
# Default API Endpoint : http://localhost:22580
```

### FestivalsIdentityAPI

The FestivalsIdentityAPI is documented in detail [here](./DOCUMENTATION.md).

## Architecture

There are a three diffrent security mechanisms to secure the festivalsapp backend, at first every party needs a valid client certificate from the FestivalsApp Root CA to communicate with other partys via mTLS, for more information [see the festivals-pki repository](https://github.com/Festivals-App/festivals-pki). After secure communication is established, clients need either an API key for the read-only parts of the FestivalsAPI or an JSON Web Token ([JWT](https://de.wikipedia.org/wiki/JSON_Web_Token)) for everything else. The JWT is used to implement a role-based access control ([RBAC](https://de.wikipedia.org/wiki/Role_Based_Access_Control)) to decide whether the user is authorized to access the given function.

The general documentation for the Festivals App is in the [festivals-documentation](https://github.com/festivals-app/festivals-documentation) repository. 
The documentation repository contains architecture information, general deployment documentation, templates and other helpful documents.

## Engage

I welcome every contribution, whether it is a pull request or a fixed typo. The best place to discuss questions and suggestions regarding the festivals-identity-server is the [issues](https://github.com/festivals-app/festivals-identity-server/issues/) section. More general information and a good starting point if you want to get involved is the [festival-documentation](https://github.com/Festivals-App/festivals-documentation) repository.

The following channels are available for discussions, feedback, and support requests:

| Type                     | Channel                                                |
| ------------------------ | ------------------------------------------------------ |
| **General Discussion**   | <a href="https://github.com/festivals-app/festivals-documentation/issues/new/choose" title="General Discussion"><img src="https://img.shields.io/github/issues/festivals-app/festivals-documentation/question.svg?style=flat-square"></a> </a>   |
| **Other Requests**    | <a href="mailto:simon.cay.gaus@gmail.com" title="Email me"><img src="https://img.shields.io/badge/email-Simon-green?logo=mail.ru&style=flat-square&logoColor=white"></a>   |

#### Licensing

Copyright (c) 2020-2024 Simon Gaus. Licensed under the [**GNU Lesser General Public License v3.0**](./LICENSE)