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
Authorization: Bearer <Header>.<Payload>.<Signatur>
Authorization: Basic <base64 encoded user:password>
```

If you have the authorization to call the given endpoint is determined by your [user role](./auth/user.go).

### Requests

The FestivalsIdentityAPI supports the HTTP `GET`, `POST`, `PATCH` and `DELETE` methods.

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

The `server-info` object is providing build time, git ref, service and version of the binary that is running.

```json
{
   "BuildTime":  "string",
   "GitRef":     "string",
   "Service":    "string",
   "Version":    "string"
}
```

In production builds the

* `BuildTime` field will have the format `Sun Apr 13 13:55:44 UTC 2025`.
* `GitRef` field will have the format `refs/tags/v2.2.0`, see "[Git References](https://git-scm.com/book/en/v2/Git-Internals-Git-References)" in the Git documentation.
* `Service` field will reference a [Service](https://github.com/Festivals-App/festivals-server-tools/blob/main/heartbeattools.go) type.
* `Version` field will have the format `v2.2.0`.

#### GET `/info`

This path will return a `server-info` object or an error.

>Authorization: `JWT` with user role set to `ADMIN`

Returns

* Returns the info object
* Codes `200`/`40x`/`50x`
* `data` or `error` field

------------------------------------------------------------------------------------

#### GET `/version`

Returns the release version of the server running. In production builds this will have the format `v1.0.2` but
for manual builds this will may be `development`.

>Authorization: `JWT` with user role set to `ADMIN`

Returns:

* The version of the server application.
* Codes `200`/`40x`/`50x`
* server version as a string `text/plain`

------------------------------------------------------------------------------------

#### POST `/update`

Tries to update to the newest release on github and then restart the service.

>Authorization: `JWT` with user role set to `ADMIN`

Returns:

* The version of the server application.
* Codes `202`/`40x`/`50x`
* server version as a string `text/plain`

------------------------------------------------------------------------------------

#### GET `/health`

>Authorization: `JWT` with user role set to `ADMIN`

Returns:

* Always returns HTTP status code 200
* Code `200`
* empty `text/plain`

------------------------------------------------------------------------------------

#### GET `/log`

Returns the info log as a string. The log format is defined [here](https://github.com/Festivals-App/festivals-server-tools/blob/main/DOCUMENTATION.md#loggertools).

>Authorization: `JWT` with user role set to `ADMIN`

Returns:

* Returns a string
* Codes `200`/`40x`/`50x`
* empty or `text/plain`

------------------------------------------------------------------------------------

#### GET `/log/trace`

Returns the trace log as a string. The log format is defined [here](https://github.com/Festivals-App/festivals-server-tools/blob/main/DOCUMENTATION.md#loggertools).

>Authorization: `JWT` with user role set to `ADMIN`

Returns:

* Returns a `string`
* Codes `200`/`40x`/`50x`
* empty or `text/plain`

## Users

### POST `/users/signup`

Signup to the festivalsapp backend.

>Authorization: `API Token`

Example:  
      `POST https://localhost:22580/users/signup`  
      `BODY: { "email": "your email", "password": "<your password>" }`

Returns:

* Returns nothing but a 201 status code.
* Codes `201`/`40x`/`50x`
* Nothing or `error` field

------------------------------------------------------------------------------------

### GET `/users/login`

Login to the festivalsapp backend.

>Authorization: `API Token` & `Basic Auth`

Examples:  
    `GET https://localhost:22580/users/login`

Returns:

* Returns the JWT on success.
* Codes `200`/`40x`/`50x`
* The raw JWT or `error` field

------------------------------------------------------------------------------------

### GET `/users/refresh`

Refresh the JWT to the festivalsapp backend. This will only refresh the users claims but not the expiration date. 

>Authorization: `JWT` with any user role

Examples:  
    `GET https://localhost:22580/users/refresh`

Returns:

* Returns the JWT on success.
* Codes `200`/`40x`/`50x`
* The raw JWT or `error` field

------------------------------------------------------------------------------------

### GET `/users`

Retruns all registered users.

>Authorization: `JWT` with user role set to `ADMIN`

Examples:  
    `GET https://localhost:22580/users`

Returns:

* Returns the users on success.
* Codes `200`/`40x`/`50x`
* `data` or `error` field

------------------------------------------------------------------------------------

### POST `/users/{objectID}/change-password`

Change the password of the given user.

>Authorization: `JWT` with any user role

Examples:  
   `POST https://localhost:22580/users/3/change-password`
   `BODY: { "old-password": "<your old password>", "new-password": "<your new-password>" }`

Returns:

* Returns nothing on success but a 200 status code.
* Codes `200`/`40x`/`50x`
* Nothing or `error` field
  
------------------------------------------------------------------------------------

### POST `/users/{objectID}/suspend`

Suspends the given user.

>Authorization: `JWT` with user role set to `ADMIN`

Examples:  
    `POST https://localhost:22580/users/3/suspend`

Returns:

* Returns nothing on success but a 200 status code.
* Codes `200`/`40x`/`50x`
* Nothing or `error` field

------------------------------------------------------------------------------------

### POST `/users/{objectID}/role/{resourceID}`

Sets the given user role for the given user. See [here](jwt/user.go) for possible values.

>Authorization: `JWT` with user role set to `ADMIN`

Examples:  
    `POST https://localhost:22580/users/3/role/42`

Returns:

* Returns nothing on success but a 200 status code.
* Codes `200`/`40x`/`50x`
* Nothing or `error` field

------------------------------------------------------------------------------------

### POST `/users/{objectID}/{festival|artist|location}/{resourceID}`

Associates the given user with the specified festival, artist or location.

>Authorization: `JWT` with user role set to `ADMIN` or `service key`

Examples:  
    `POST https://localhost:22580/users/3/artist/134`

Returns:

* Returns nothing on success but a 200 status code.
* Codes `200`/`40x`/`50x`
* Nothing or `error` field

------------------------------------------------------------------------------------

### DELETE `/users/{objectID}/{festival|artist|location}/{resourceID}`

Removes the association between the given user and the specified festival, artist or location.

>Authorization: `JWT` with user role set to `ADMIN` or `service key`

Examples:  
    `DELETE https://localhost:22580/users/3/festival/26`

Returns:

* Returns nothing on success but a 200 status code.
* Codes `200`/`40x`/`50x`
* Nothing or `error` field

------------------------------------------------------------------------------------

## Validation-Key

------------------------------------------------------------------------------------

### GET `/validation-key`

Returns the public key used to sign the jwt's issued by this identity service.

>Authorization: `JWT` with user role set to `ADMIN` or `service key`

Examples:  
    `GET https://localhost:22580/validation-keys`

Returns:

* Returns the validation key on success.
* Codes `200`/`40x`/`50x`
* raw file or `error` field

------------------------------------------------------------------------------------

## Service-Keys

------------------------------------------------------------------------------------

### GET `/service-keys`

Returns all registered service keys.

>Authorization: `JWT` with user role set to `ADMIN` or `service key`

Examples:  
    `GET https://localhost:22580/service-keys`

Returns:

* Returns the service keys on success.
* Codes `200`/`40x`/`50x`
* `data` or `error` field

------------------------------------------------------------------------------------

### POST `/service-keys`

Registers a new service key.

>Authorization: `JWT` with user role set to `ADMIN` or `service key`
  
Examples:  
    `POST https://localhost:22580/service-keys`
    `BODY: { "service_key": "<service key>", "service_key_comment": "<Comment for the service key>" }`

Returns:

* Returns nothing on success but a 201 status code.
* Codes `201`/`40x`/`50x`
* Nothing or `error` field

------------------------------------------------------------------------------------

### PATCH `/service-keys/{objectID}`

Updates the given service key.

>Authorization: `JWT` with user role set to `ADMIN`
  
Examples:  
    `POST https://localhost:22580/service-keys/23`
    `BODY: { "service_key": "<service key>", "service_key_comment": "<Comment for the service key>" }`

Returns:

* Returns nothing on success but a 200 status code.
* Codes `200`/`40x`/`50x`
* Nothing or `error` field

------------------------------------------------------------------------------------

### DELETE `/service-keys/{objectID}`

Deletes the given service key.

>Authorization: `JWT` with user role set to `ADMIN`

Examples:  
    `DELETE https://localhost:22580/service-keys/23`

Returns:

* Returns nothing on success but a 200 status code.
* Codes `200`/`40x`/`50x`
* Nothing or `error` field

------------------------------------------------------------------------------------

## API-Keys

------------------------------------------------------------------------------------

### GET `/api-keys`

Returns all registered API keys.

>Authorization: `JWT` with user role set to `ADMIN` or `service key`

Examples:  
    `GET https://localhost:22580/api-keys`

Returns:

* Returns the API keys on success.
* Codes `200`/`40x`/`50x`
* `data` or `error` field
 
------------------------------------------------------------------------------------

### POST `/api-keys`

Registers a new API key.

>Authorization: `JWT` with user role set to `ADMIN`
  
Examples:  
    `POST https://localhost:22580/api-keys`
    `BODY: { "api_key": "<api key>", "api_key_comment": "<Comment for the api key>" }`

Returns:

* Returns nothing on success but a 201 status code.
* Codes `201`/`40x`/`50x`
* Nothing or `error` field

------------------------------------------------------------------------------------

### PATCH `/api-keys/{objectID}`

Updates the given API key.

>Authorization: `JWT` with user role set to `ADMIN`
  
Examples:  
    `POST https://localhost:22580/api-keys/23`
    `BODY: { "api_key": "<api key>", "api_key_comment": "<Comment for the API key>" }`

Returns:

* Returns nothing on success but a 200 status code.
* Codes `200`/`40x`/`50x`
* Nothing or `error` field

------------------------------------------------------------------------------------

### DELETE `/api-keys/{objectID}`

Deletes the given api key.

>Authorization: `JWT` with user role set to `ADMIN`
Examples:  
    `DELETE https://localhost:22580/api-keys/23`

Returns:

* Returns nothing on success but a 200 status code.
* Codes `200`/`40x`/`50x`
* Nothing or `error` field

------------------------------------------------------------------------------------
