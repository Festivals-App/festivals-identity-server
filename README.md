<p align="center">
   <a href="https://github.com/festivals-app/festivals-identity-server/commits/" title="Last Commit"><img src="https://img.shields.io/github/last-commit/festivals-app/festivals-identity-server?style=flat" alt="Commits Shield"></a>
   <a href="https://github.com/festivals-app/festivals-identity-server/issues" title="Open Issues"><img src="https://img.shields.io/github/issues/festivals-app/festivals-identity-server?style=flat" alt="Issues Shield"></a>
   <a href="./LICENSE" title="License"><img src="https://img.shields.io/github/license/festivals-app/festivals-identity-server.svg" alt="License Shield"></a>
</p>

<h1 align="center">
  <br/><br/>
    FestivalsApp Identity Server
  <br/><br/>
</h1>

A lightweight Go server application providing the [FestivalsIdentityAPI](DOCUMENTATION.md), a RESTful API that handles
all authentication and authorization needs for FestivalsApp components.

![Figure 1: Architecture Overview Highlighted](https://github.com/Festivals-App/festivals-documentation/blob/main/images/architecture/export/architecture_overview_identity.svg "Figure 1: Architecture Overview Highlighted")

<hr/>
<p align="center">
  <a href="#development">Development</a> •
  <a href="#deployment">Deployment</a> •
  <a href="#engage">Engage</a>
</p>
<hr/>

The FestivalsApp backend is secured using three different mechanisms to ensure both secure communication and controlled access:

1. **Mutual TLS (mTLS)** – Every party must have a valid client certificate issued by the FestivalsApp Root CA
  to establish secure communication with other services. This prevents unauthorized access at the transport layer.
  For more details, refer to the [festivals-pki](https://github.com/Festivals-App/festivals-pki) repository.  
2. **API Keys** – Required for accessing read-only parts of the FestivalsAPI. These keys provide a simple way
  to authenticate services and users that do not require full access.
3. **JSON Web Tokens (JWTs)** – Used for all other interactions. JWTs enable role-based access control (RBAC),
  ensuring users are authorized to access specific functions based on their assigned roles and also implement
  resource access based on user identity. The system verifies JWTs on every request to enforce access restrictions dynamically.

In addition to these mechanisms, the backend enforces strict firewall rules and network segmentation
to minimize exposure to unauthorized access.

## Development

The FestivalsApp Identity Server follows a modular structure for clarity and maintainability. The `database` directory
for managing the database, while `auth` handles core authentication logic, `server` manages API routes and middleware
and `operation` documents deployment and environment. GitHub Actions are in `.github`, and `.vscode` provides recommended
settings. The entry point is main.go, with dependencies in go.mod and go.sum.
Refer to [FestivalsIdentityAPI Documentation](DOCUMENTATION.md) for details on available endpoints.

### Requirements

- [Golang](https://go.dev/) Version 1.24.1+
- [Visual Studio Code](https://code.visualstudio.com/download) 1.98.2+
  - Plugin recommendations are managed via [workspace recommendations](https://code.visualstudio.com/docs/editor/extension-marketplace#_recommended-extensions).
- [Bash script](https://en.wikipedia.org/wiki/Bash_(Unix_shell)) friendly environment

## Deployment

The Go binaries are able to run without system dependencies so there are not many requirements for the system
to run the festivals-identity-server binary, just follow the [**deployment guide**](./operation/DEPLOYMENT.md) for
deploying it inside a virtual machine or the [**local deployment guide**](./operation/local/README.md) for
running it on your macOS developer machine.

## Engage

I welcome every contribution, whether it is a pull request or a fixed typo. The best place to discuss questions
and suggestions regarding the festivals-identity-server is the [issues](https://github.com/festivals-app/festivals-identity-server/issues/) section.
More general information and a good starting point if you want to get involved is
the [festival-documentation](https://github.com/Festivals-App/festivals-documentation) repository.

The following channels are available for discussions, feedback, and support requests:

| Type                     | Channel                                                |
| ------------------------ | ------------------------------------------------------ |
| **General Discussion**   | <a href="https://github.com/festivals-app/festivals-documentation/issues/new/choose" title="General Discussion"><img src="https://img.shields.io/github/issues/festivals-app/festivals-documentation/question.svg?style=flat-square"></a> </a>   |
| **Other Requests**    | <a href="mailto:simon@festivalsapp.org" title="Email me"><img src="https://img.shields.io/badge/email-Simon-green?logo=mail.ru&style=flat-square&logoColor=white"></a>   |

### Licensing

Copyright (c) 2020-2025 Simon Gaus. Licensed under the [**GNU Lesser General Public License v3.0**](./LICENSE)
