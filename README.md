<p align="center">
   <a href="https://github.com/festivals-app/festivals-identity-server/commits/" title="Last Commit"><img src="https://img.shields.io/github/last-commit/festivals-app/festivals-identity-server?style=flat"></a>
   <a href="https://github.com/festivals-app/festivals-identity-server/issues" title="Open Issues"><img src="https://img.shields.io/github/issues/festivals-app/festivals-identity-server?style=flat"></a>
   <a href="./LICENSE" title="License"><img src="https://img.shields.io/github/license/festivals-app/festivals-identity-server.svg"></a>
</p>

<h1 align="center">
  <br/><br/>
    Festivals Identity Server
  <br/><br/>
</h1>

<p align="center">
  <a href="#development">Development</a> •
  <a href="#deployment">Deployment</a> •
  <a href="#festivalsidentityapi">FestivalsIdentityAPI</a> •
  <a href="#architecture">Architecture</a> •
  <a href="#engage">Engage</a> •
  <a href="#licensing">Licensing</a>
</p>

A lightweight go server app providing a simple RESTful API using [go-chi/chi](https://github.com/go-chi/chi) and [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql).

![Figure 1: Architecture Overview Highlighted](https://github.com/Festivals-App/festivals-documentation/blob/main/images/architecture/overview_id.png "Figure 1: Architecture Overview Highlighted")

## Development

1. Make server run
2. Make server serves tls
3. Make other server serve tls
4. implement jwt to validate requests to other partys (especialy for admin requests)


TBA

### Requirements

-  go 1.20

### Setup development

TBA

## Deployment

The install, update and uninstall scripts should work with any system that uses *systemd* and *firewalld* or *ufw*. 
Additionally the scripts will somewhat work under macOS but won't configure the firewall or launch service.

Installing
```bash
curl -o install.sh https://raw.githubusercontent.com/Festivals-App/festivals-identity-server/master/operation/install.sh
chmod +x install.sh
sudo ./install.sh <mysql_root_pw> <mysql_backup_pw> <read_write_pw>
sudo nano /etc/mysql/mysql.conf.d/mysqld.cnf          // edit bind-address=<private-ip>
```

Updating
```bash
curl -o update.sh https://raw.githubusercontent.com/Festivals-App/festivals-identity-server/master/operation/update.sh
chmod +x update.sh
sudo ./update.sh
```

### Build and Run manually
```bash
cd $GOPATH/src/github.com/Festivals-App/festivals-fileserver
go build main.go
./main

# Default API Endpoint : http://localhost:1910
```

### FestivalsIdentityAPI

The FestivalsIdentityAPI is documented in detail [here](./DOCUMENTATION.md).

## Architecture

The general documentation for the Festivals App is in the [festivals-documentation](https://github.com/festivals-app/festivals-documentation) repository. 
The documentation repository contains architecture information, general deployment documentation, templates and other helpful documents.

## Engage

I welcome every contribution, whether it is a pull request or a fixed typo. The best place to discuss questions and suggestions regarding the festivals-identity-server is the [issues](https://github.com/festivals-app/festivals-identity-server/issues/) section. More general information and a good starting point if you want to get involved is the [festival-documentation](https://github.com/Festivals-App/festivals-documentation) repository.

The following channels are available for discussions, feedback, and support requests:

| Type                     | Channel                                                |
| ------------------------ | ------------------------------------------------------ |
| **General Discussion**   | <a href="https://github.com/festivals-app/festivals-documentation/issues/new/choose" title="General Discussion"><img src="https://img.shields.io/github/issues/festivals-app/festivals-documentation/question.svg?style=flat-square"></a> </a>   |
| **Other Requests**    | <a href="mailto:simon.cay.gaus@gmail.com" title="Email me"><img src="https://img.shields.io/badge/email-Simon-green?logo=mail.ru&style=flat-square&logoColor=white"></a>   |

## Licensing

Copyright (c) 2020-2023 Simon Gaus.

Licensed under the **GNU Lesser General Public License v3.0** (the "License"); you may not use this file except in compliance with the License.

You may obtain a copy of the License at https://www.gnu.org/licenses/lgpl-3.0.html.

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the [LICENSE](./LICENSE) for the specific language governing permissions and limitations under the License.