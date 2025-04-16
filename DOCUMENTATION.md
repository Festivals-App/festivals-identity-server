<h1 align="center">
    FestivalsIdentityAPI Documentation
</h1>

<p align="center">
  <a href="#overview">Overview</a> •
  <a href="#server-status">Server-Status</a> •
  <a href="#users">Users</a> •
  <a href="#validation-key">Validation-Key</a> •
  <a href="#service-keys">Service-Keys</a> •
  <a href="#api-keys">API-Keys</a>
</p>

### Used Languages

* Documentation: `Markdown`, `HTML`
* Server Application: `golang`
* Deployment: `bash`

### Authentication & Authorization

To authenticate to the `FestivalsIdentityAPI` you need to either provide an API/service key via a custom header or a JWT
with your requests authorization header, for login you need to use basic authentication alsongside an API key.

```text
Api-Key: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
Service-Key: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
Authorization: Bearer <JWT>
Authorization: Basic <base64 encoded user:password>
```

If you have the authorization to call the given endpoint is determined by your [user role](./auth/user.go).

#### Making a request with curl

```bash
curl -H "X-Request-ID: <uuid>" -H "Authorization: Bearer <JWT>" --cacert ca.crt --cert client.crt --key client.key https://identity-0.festivalsapp.home/info
```

### Response

For `GET` requests that are handled gracefully by the server will always return the requested ressource directly,
otherwise an `error` field is returned and will always contain a string with the error message.

```json
{
   "error": "An error occured"
}
```

## Overview

[Server-Status](#server-status)

* GET              `/info`
* GET              `/version`
* POST             `/update`
* GET              `/health`
* GET              `/log`
* GET              `/log/trace`

[Users](#users)

* POST             `/users/signup`
* GET              `/users/login`
* GET              `/users/refresh`
* GET              `/users`
* POST             `/users/{objectID}/change-password`
* POST             `/users/{objectID}/suspend`
* POST             `/users/{objectID}/role/{resourceID}`
* POST             `/users/{objectID}/{festival|artist|location}/{resourceID}`
* DELETE           `/users/{objectID}/{festival|artist|location}/{resourceID}`

[Validation-Key](#validation-key)

* GET                         `/validation-key`

[Service-Keys](#service-keys)

* GET, POST, PATCH, DELETE    `/service-keys`

[API-Keys](#api-keys)

* GET, POST, PATCH, DELETE    `/api-keys`

## Server Status

The **server status routes** serve status-related information.

It is commonly used for health checks, CI/CD diagnostics, or runtime introspection. This route uses
a `server-info` object containing metadata about the currently running binary, such as build time,
Git reference, service name, and version.

**`server-info`** object

```json
{
  "BuildTime": "string",
  "GitRef": "string",
  "Service": "string",
  "Version": "string"
}
```

| Field      | Description                                                                 |
|------------|-----------------------------------------------------------------------------|
| `BuildTime` | Timestamp of the binary build. Format: `Sun Apr 13 13:55:44 UTC 2025`       |
| `GitRef`    | Git reference used for the build. Format: `refs/tags/v2.2.0` [🔗 Git Docs](https://git-scm.com/book/en/v2/Git-Internals-Git-References) |
| `Service`   | Service identifier. Matches a defined [Service type](https://github.com/Festivals-App/festivals-server-tools/blob/main/heartbeattools.go) |
| `Version`   | Version tag of the deployed binary. Format: `v2.2.0`                        |

> In production builds, these values are injected at build time and reflect the deployment source and context.

#### GET `/info`

Returns the `server-info`.

Example:  
  `GET https://identity-0.festivalsapp.home/info`

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN`.

**Response**

* `data` or `error` field
* Codes `200`/`40x`/`50x`

------------------------------------------------------------------------------------

#### GET `/version`

Returns the release version of the server.

> In production builds this will have the format `v2.2.0` but
for manual builds this will may be `development`.

Example:  
  `GET https://identity-0.festivalsapp.home:22580/version`

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN`.

**Response**

* Server version as a string `text/plain`.
* Codes `200`/`40x`/`50x`

------------------------------------------------------------------------------------

#### POST `/update`

Updates to the newest release on github and restarts the service.

Example:  
  `POST https://identity-0.festivalsapp.home:22580/update`

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN`.

**Response**

* The current version and the version the server is updated to as a string `text/plain`. Format: `v2.1.3 => v2.2.0`
* Codes `202`/`40x`/`50x`

------------------------------------------------------------------------------------

#### GET `/health`

A simple health check endpoint that returns a `200 OK` status if the service is running and able to respond.

Example:  
  `GET https://identity-0.festivalsapp.home:22580/health`

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN`.

**Response**

* Always returns `200 OK`

------------------------------------------------------------------------------------

#### GET `/log`

Returns the info log file as a string, containing all log messages except trace log entries.
See [loggertools](https://github.com/Festivals-App/festivals-server-tools/blob/main/DOCUMENTATION.md#loggertools) for log format.

Example:  
  `GET https://identity-0.festivalsapp.home:22580/log`

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN`.

**Response**

* Returns a string as `text/plain`
* Codes `200`/`40x`/`50x`

------------------------------------------------------------------------------------

#### GET `/log/trace`

Returns the trace log file as a string, containing all remote calls to the server.
See [loggertools](https://github.com/Festivals-App/festivals-server-tools/blob/main/DOCUMENTATION.md#loggertools) for log format.

Example:  
  `GET https://identity-0.festivalsapp.home:22580/log/trace`

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN`.

**Response**

* Returns a string as `text/plain`
* Codes `200`/`40x`/`50x`

------------------------------------------------------------------------------------

## Users

The **user routes** serve user related endpoints including signup, login and user resources.
This route uses a `user` object containing metadata about a user.

**`user`** object

```json
{
  "user_id": "string",
  "user_email": "string",
  "user_createdat": "string",
  "user_updatedat": "string",
  "user_role": "int"
}
```

| Field            | Description                                                           |
|------------------|-----------------------------------------------------------------------|
| `user_id`        | The users ID.                                                         |
| `user_email`     | The users email address                                               |
| `user_createdat` | The date the user was created. Format: `2024-03-27T01:49:32Z`         |
| `user_updatedat` | The date the user was updated. Format: `2024-03-27T01:49:32Z`         |
| `user_role`      | One of the [user role](./auth/user.go) values.                        |

------------------------------------------------------------------------------------

### POST `/users/signup`

Signup to the festivalsapp backend as a creator.

Example:  
  `POST https://identity-0.festivalsapp.home:22580/users/signup`  
  `BODY: { "email": "your email", "password": "<your password>" }`

**Authorization**
Requires a valid `API-Key`.

**Response**

* Returns `201 CREATED` on success or `error` field on failure.
* Codes `201`/`40x`/`50x`

------------------------------------------------------------------------------------

### GET `/users/login`

Login to the festivalsapp backend.

Examples:  
    `GET https://identity-0.festivalsapp.home:22580/users/login`

**Authorization**
Requires a valid `API-Key` and correct `Basic Auth` credentials.

**Response**

* Returns the raw `JWT` on success or `error` field on failure.
* Codes `200`/`40x`/`50x`

------------------------------------------------------------------------------------

### GET `/users/refresh`

Refreshes the `JWT`. This will only refresh the users claims but not the expiration date of the token.

Examples:  
    `GET https://identity-0.festivalsapp.home:22580/users/refresh`

**Authorization**
Requires a valid `JWT` token with any user role.

**Response**

* Returns the refreshed `JWT` on success or `error` field on failure.
* Codes `200`/`40x`/`50x`

------------------------------------------------------------------------------------

### GET `/users`

Retruns all registered users as a list of `user`s.

Examples:  
    `GET https://identity-0.festivalsapp.home:22580/users`

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN`.

**Response**

* `data` or `error` field
* Codes `200`/`40x`/`50x`

------------------------------------------------------------------------------------

### POST `/users/{objectID}/change-password`

Change the password of the given user.

Examples:  
   `POST https://identity-0.festivalsapp.home:22580/users/3/change-password`
   `BODY: { "old-password": "<your old password>", "new-password": "<your new-password>" }`

**Authorization**
Requires a valid `JWT` token with any user role.

**Response**

* Returns `200 OK` on success and `error` on failure.
* Codes `200`/`40x`/`50x`
  
------------------------------------------------------------------------------------

### POST `/users/{objectID}/suspend`

Suspends the given user.

Examples:  
    `POST https://identity-0.festivalsapp.home:22580/users/3/suspend`

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN`.

**Response**

* Returns `200 OK` on success and `error` on failure.
* Codes `200`/`40x`/`50x`

------------------------------------------------------------------------------------

### POST `/users/{objectID}/role/{resourceID}`

Sets the given user role for the given user. See [here](jwt/user.go) for possible values.

Examples:  
    `POST https://identity-0.festivalsapp.home:22580/users/3/role/42`

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN`.

**Response**

* Returns `200 OK` on success and `error` on failure.
* Codes `200`/`40x`/`50x`

------------------------------------------------------------------------------------

### POST `/users/{objectID}/{festival|artist|location}/{resourceID}`

Associates the given user with the specified festival, artist or location.

Examples:  
    `POST https://identity-0.festivalsapp.home:22580/users/3/artist/134`

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN` or valid `service key`.

**Response**

* Returns `200 OK` on success and `error` on failure.
* Codes `200`/`40x`/`50x`

------------------------------------------------------------------------------------

### DELETE `/users/{objectID}/{festival|artist|location}/{resourceID}`

Removes the association between the given user and the specified festival, artist or location.

Examples:  
    `DELETE https://identity-0.festivalsapp.home:22580/users/3/festival/26`

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN` or valid `service key`.

**Response**

* Returns `200 OK` on success and `error` on failure.
* Codes `200`/`40x`/`50x`

------------------------------------------------------------------------------------

## Validation-Key

The **validation-key route** provides the public key used to sign `JWT`'s issued by this identity service
in order for other services to validate those `JWT`s.

------------------------------------------------------------------------------------

### GET `/validation-key`

Returns the public key used to sign the `JWT`'s issued by this identity service.

Examples:  
    `GET https://identity-0.festivalsapp.home:22580/validation-keys`

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN` or valid `service key`.

**Response**

* Returns the raw validation key as `text/plain` on success or `error` on failure.
* Codes `200`/`40x`/`50x`

------------------------------------------------------------------------------------

## Service-Keys

The **service-key routes** serve service-key related endpoints including retrieving, creating and deleting service-keys.
This route uses a `service-key` object containing metadata about a service-key.

**`service-key`** object

```json
{
  "service_key_id": "int",
  "service_key": "string",
  "service_key_comment": "string"
}
```

| Field                 | Description                                                        |
|-----------------------|--------------------------------------------------------------------|
| `service_key_id`      | The ID of the service key.                                         |
| `service_key`         | The service key.                                                   |
| `service_key_comment` | The comment for the service key                                    |

------------------------------------------------------------------------------------

### GET `/service-keys`

Returns all registered service keys as a list of `service-key`s.

Examples:  
    `GET https://identity-0.festivalsapp.home:22580/service-keys`

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN` or valid `service key`.

**Response**

* `data` or `error` field
* Codes `200`/`40x`/`50x`

------------------------------------------------------------------------------------

### POST `/service-keys`

Registers a new service key.

Examples:  
    `POST https://identity-0.festivalsapp.home:22580/service-keys`
    `BODY: { "service_key": "<service key>", "service_key_comment": "<Comment for the service key>" }`

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN`.

**Response**

* Returns `201 Created` on success and `error` on failure.
* Codes `201`/`40x`/`50x`

------------------------------------------------------------------------------------

### PATCH `/service-keys/{objectID}`

Updates the given service key.
  
Examples:  
    `POST https://identity-0.festivalsapp.home:22580/service-keys/23`
    `BODY: { "service_key": "<service key>", "service_key_comment": "<Comment for the service key>" }`

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN`.

**Response**

* Returns `200 OK` on success and `error` on failure.
* Codes `200`/`40x`/`50x`

------------------------------------------------------------------------------------

### DELETE `/service-keys/{objectID}`

Deletes the given service key.

Examples:  
    `DELETE https://identity-0.festivalsapp.home:22580/service-keys/23`

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN`.

**Response**

* Returns `200 OK` on success and `error` on failure.
* Codes `200`/`40x`/`50x`

------------------------------------------------------------------------------------

## API-Keys

The **api-key routes** serve api-key related endpoints including retrieving, creating and deleting api-keys.
This route uses a `api-key` object containing metadata about a api-key.

**`api-key`** object

```json
{
  "api_key_id": "int",
  "api_key": "string",
  "api_key_comment": "string"
}
```

| Field                 | Description                                                    |
|-----------------------|----------------------------------------------------------------|
| `api_key_id`          | The ID of the api key.                                         |
| `api_key`             | The api key.                                                   |
| `api_key_comment`     | The comment for the api key                                    |

------------------------------------------------------------------------------------

### GET `/api-keys`

Returns all registered API keys as a list of `api-key`s.

Examples:  
    `GET https://identity-0.festivalsapp.home:22580/api-keys`

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN` or valid `service key`.

**Response**

* `data` or `error` field
* Codes `200`/`40x`/`50x`

------------------------------------------------------------------------------------

### POST `/api-keys`

Registers a new API key.

Examples:  
    `POST https://identity-0.festivalsapp.home:22580/api-keys`
    `BODY: { "api_key": "<api key>", "api_key_comment": "<Comment for the api key>" }`

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN`.

**Response**

* Returns `201 Created` on success and `error` on failure.
* Codes `201`/`40x`/`50x`

------------------------------------------------------------------------------------

### PATCH `/api-keys/{objectID}`

Updates the given API key.

Examples:  
    `POST https://identity-0.festivalsapp.home:22580/api-keys/23`
    `BODY: { "api_key": "<api key>", "api_key_comment": "<Comment for the API key>" }`

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN`.

**Response**

* Returns `200 OK` on success and `error` on failure.
* Codes `200`/`40x`/`50x`

------------------------------------------------------------------------------------

### DELETE `/api-keys/{objectID}`

Deletes the given api key.

Examples:  
    `DELETE https://identity-0.festivalsapp.home:22580/api-keys/23`

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN`.

**Response**

* Returns `200 OK` on success and `error` on failure.
* Codes `200`/`40x`/`50x`

------------------------------------------------------------------------------------
