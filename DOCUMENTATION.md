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

To use the API you need to provide an API or service key via a custom header or a JWT
with your requests authorization header, for login you need to use basic authentication:

```ini
Api-Key: <api-key>
Service-Key: <service-key>
Authorization: Bearer <jwt>
Authorization: Basic <encodedcredentials>
```

### Requests

The FestivalsIdentityAPI supports the HTTP `GET`, `POST`, `PATCH` and `DELETE` methods.

### Response

Get requests that are handled gracefully by the server will always return the requested ressource directly,
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

* GET              `/validation-key`

[Service-Keys](#service-keys)

* GET              `/service-keys`
* POST             `/service-keys`
* DELETE           `/service-keys`

[API-Keys](#api-keys)

* GET              `/api-keys`
* POST             `/api-keys`
* DELETE           `/api-keys`

------------------------------------------------------------------------------------

## Server Status

Determine the state of the server.

Info object

```json
{
    "BuildTime":      string,
    "GitRef":         string,
    "Version":        string
}
```

------------------------------------------------------------------------------------

### GET `/info`

 * Authorization: JWT
 
 * Returns
      * Returns the info object 
      * Codes `200`/`40x`/`50x`
      * `data` or `error` field

------------------------------------------------------------------------------------
#### GET `/version`

 * Authorization: JWT
 
 * Returns
      * The version of the server application.
      * Codes `200`/`40x`/`50x`
      * server version as a string `text/plain`

------------------------------------------------------------------------------------
#### POST `/update`

Updates to the newest release on github and restarts the service.

 * Authorization: JWT
 * Returns
      * The version of the server application.
      * Codes `202`/`40x`/`50x`
      * server version as a string `text/plain`

------------------------------------------------------------------------------------
#### GET `/health`

 * Authorization: JWT
 
 * Returns
      * Always returns HTTP status code 200
      * Code `200`
      * empty `text/plain`

------------------------------------------------------------------------------------
#### GET `/log`

Returns the service log.

 * Authorization: JWT
 
 * Returns
      * Returns a string
      * Codes `200`/`40x`/`50x`
      * empty or `text/plain`

------------------------------------------------------------------------------------
#### GET `/log/trace`

Returns the service trace log.

 * Authorization: JWT
 
 * Returns
      * Returns a string
      * Codes `200`/`40x`/`50x`
      * empty or `text/plain`

------------------------------------------------------------------------------------
## Users

#### POST `/users/signup`

Signup to the festivalsapp backend. 

 * Authorization: API Token
 
 * Example:  
      `POST https://localhost:22580/users/signup`  
      `BODY: { "email": "your email", "password": "<your password>" }`
 * Returns:
      * Returns nothing but a 201 status code. 
      * Codes `201`/`40x`/`50x`
      * Nothing or `error` field

------------------------------------------------------------------------------------
#### GET `/users/login`

Login to the festivalsapp backend.
 
 * Authorization: API Token & Basic Auth
 
 * Examples:  
    `GET https://localhost:22580/users/login`
 * Returns:
     * Returns the JWT on success.
     * Codes `200`/`40x`/`50x`
     * The raw JWT or `error` field

------------------------------------------------------------------------------------
#### GET `/users/refresh`

Refresh the JWT to the festivalsapp backend. This will only refresh the users claims but not the expiration date. 
 
 * Authorization: JWT
 
 * Examples:  
    `GET https://localhost:22580/users/refresh`
 * Returns:
     * Returns the JWT on success.
     * Codes `200`/`40x`/`50x`
     * The raw JWT or `error` field

------------------------------------------------------------------------------------
#### GET `/users`

Retruns all registered users.
 
 * Authorization: JWT
 
 * Examples:  
    `GET https://localhost:22580/users`
 * Returns:
     * Returns the users on success.
     * Codes `200`/`40x`/`50x`
     * `data` or `error` field

------------------------------------------------------------------------------------
#### POST `/users/{objectID}/change-password`

Change the password of the given user.

 * Authorization: JWT
 
 * Examples:  
    `POST https://localhost:22580/users/3/change-password`
    `BODY: { "old-password": "<your old password>", "new-password": "<your new-password>" }`
 * Returns:
     * Returns nothing on success but a 200 status code.
     * Codes `200`/`40x`/`50x`
     * Nothing or `error` field
  
------------------------------------------------------------------------------------
#### POST `/users/{objectID}/suspend`

Suspends the given user.

 * Authorization: JWT
 
 * Examples:  
    `POST https://localhost:22580/users/3/suspend`
 * Returns:
     * Returns nothing on success but a 200 status code.
     * Codes `200`/`40x`/`50x`
     * Nothing or `error` field

------------------------------------------------------------------------------------
#### POST `/users/{objectID}/role/{resourceID}`

Sets the given user role for the given user. See [here](jwt/user.go) for possible values.

 * Authorization: JWT
 
 * Examples:  
    `POST https://localhost:22580/users/3/role/42`
 * Returns:
     * Returns nothing on success but a 200 status code.
     * Codes `200`/`40x`/`50x`
     * Nothing or `error` field

------------------------------------------------------------------------------------
#### POST `/users/{objectID}/{festival|artist|location}/{resourceID}`

Associates the given user with the specified festival, artist or location.

 * Authorization: JWT or service key
 
 * Examples:  
    `POST https://localhost:22580/users/3/artist/134`
 * Returns:
     * Returns nothing on success but a 200 status code.
     * Codes `200`/`40x`/`50x`
     * Nothing or `error` field

------------------------------------------------------------------------------------
#### DELETE `/users/{objectID}/{festival|artist|location}/{resourceID}`

Removes the association between the given user and the specified festival, artist or location.

 * Authorization: JWT or service key
 
 * Examples:  
    `DELETE https://localhost:22580/users/3/festival/26`
 * Returns:
     * Returns nothing on success but a 200 status code.
     * Codes `200`/`40x`/`50x`
     * Nothing or `error` field
 
------------------------------------------------------------------------------------
## Validation-Key

#### GET `/validation-key`

Returns the public key used to sign the jwt's issued by this identity service.

 * Authorization: JWT or service key
 
 * Examples:  
    `GET https://localhost:22580/validation-keys`
 * Returns:
     * Returns the validation key on success.
     * Codes `200`/`40x`/`50x`
     * raw file or `error` field

------------------------------------------------------------------------------------
## Service-Keys

#### GET `/service-keys`

Returns all registered service keys.

 * Authorization: JWT or service key
 
 * Examples:  
    `GET https://localhost:22580/service-keys`
 * Returns:
     * Returns the service keys on success.
     * Codes `200`/`40x`/`50x`
     * `data` or `error` field
 
------------------------------------------------------------------------------------
#### POST `/service-keys`

Registers a new service key.

 * Authorization: JWT
  
 * Examples:  
    `POST https://localhost:22580/service-keys`
    `BODY: { "service_key": "<service key>", "service_key_comment": "<Comment for the service key>" }`
 * Returns:
     * Returns nothing on success but a 201 status code.
     * Codes `201`/`40x`/`50x`
     * Nothing or `error` field
 
------------------------------------------------------------------------------------
#### PATCH `/service-keys/{objectID}`

Updates the given service key.

 * Authorization: JWT
  
 * Examples:  
    `POST https://localhost:22580/service-keys/23`
    `BODY: { "service_key": "<service key>", "service_key_comment": "<Comment for the service key>" }`
 * Returns:
     * Returns nothing on success but a 200 status code.
     * Codes `200`/`40x`/`50x`
     * Nothing or `error` field

------------------------------------------------------------------------------------
#### DELETE `/service-keys/{objectID}`

Deletes the given service key.

 * Authorization: JWT
 
 * Examples:  
    `DELETE https://localhost:22580/service-keys/23`
 * Returns:
     * Returns nothing on success but a 200 status code.
     * Codes `200`/`40x`/`50x`
     * Nothing or `error` field

------------------------------------------------------------------------------------
## API-Keys

#### GET `/api-keys`

Returns all registered API keys.

 * Authorization: JWT or service key
 
 * Examples:  
    `GET https://localhost:22580/api-keys`
 * Returns:
     * Returns the API keys on success.
     * Codes `200`/`40x`/`50x`
     * `data` or `error` field
 
------------------------------------------------------------------------------------
#### POST `/api-keys`

Registers a new API key.

 * Authorization: JWT
  
 * Examples:  
    `POST https://localhost:22580/api-keys`
    `BODY: { "api_key": "<api key>", "api_key_comment": "<Comment for the api key>" }`
 * Returns:
     * Returns nothing on success but a 201 status code.
     * Codes `201`/`40x`/`50x`
     * Nothing or `error` field
 
------------------------------------------------------------------------------------
#### PATCH `/api-keys/{objectID}`

Updates the given API key.

 * Authorization: JWT
  
 * Examples:  
    `POST https://localhost:22580/api-keys/23`
    `BODY: { "api_key": "<api key>", "api_key_comment": "<Comment for the API key>" }`
 * Returns:
     * Returns nothing on success but a 200 status code.
     * Codes `200`/`40x`/`50x`
     * Nothing or `error` field

------------------------------------------------------------------------------------
#### DELETE `/api-keys/{objectID}`

Deletes the given api key.

 * Authorization: JWT
 
 * Examples:  
    `DELETE https://localhost:22580/api-keys/23`
 * Returns:
     * Returns nothing on success but a 200 status code.
     * Codes `200`/`40x`/`50x`
     * Nothing or `error` field

------------------------------------------------------------------------------------