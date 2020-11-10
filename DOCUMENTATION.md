<!--suppress ALL -->

<h1 align="center">
    FestivalsIdentityAPI Documentation
</h1>

<p align="center">
  <a href="#overview">Overview</a> •
  <a href="#user">User</a> •
  <a href="#api-keys">API-Keys</a> •
  <a href="#status">Status</a>
</p>

### Used Languages

* Documentation: `Markdown`, `HTML`
* Server Application: `golang`
* Deployment: `bash`

### Authentication

To use the API you need to provide an API key with your requests HTTP header:
```
'API-KEY':'YOUR_API_KEY_GOES_HERE'
```

### Requests

The FestivalFilesAPI API supports the HTTP `GET`, `POST`, `PATCH` and `DELETE` methods.

#### Query Parameter

* `width`  
    The width parameter expects a simple string. `width=^[0-9]+`
* `height`  
    The height parameter expects numbers separated by a comma. `height=[0-9]+`

### Response

Get requests that are handled gracefully by the server will always return the requested file directly,
otherwise an `error` field is returned and will always contain a string with the error message.
```
{
    "error": "An error occured"
}
```

## Overview

[User](#user)
* POST              `/user/login`

[Status](#status)
* GET               `/status`

------------------------------------------------------------------------------------
## User

#### POST `/images`

Upload an image.

 * Examples:  
      `POST https://localhost:8080/images`  
        
 * Returns
      * Returns the image 
      * Codes `200`/`40x`/`50x`
      * image or `error` field

------------------------------------------------------------------------------------
#### GET `/images/{objectID}`

Get the image with the given `objectID`.

* Query Parameter:  
    `height`, `width`: Crops the image to the given size, respecting the aspect ratio of the image.

 * Examples:  
    `GET https://localhost:8080/images/test.jpg`
    `GET https://localhost:8080/images/test.jpg?width=500&height=500`
      
 * Returns 
     * Returns the image on success.
     * Codes `201`/`40x`/`50x`
     * image or `error` field

------------------------------------------------------------------------------------
#### DELETE `/images/{objectID}`

Delete the images with the given `objectID`.
 
 * Examples:  
    `DELETE https://localhost:8080/images/test.jpg`
    
 * Returns 
     * Returns no object.
     * Codes `204`/`40x`/`50x`
     * `data` or `error` field
 
------------------------------------------------------------------------------------
#### GET `/status`

Get the status of the fileserver.

 * Examples:  
    `GET https://localhost:8080/status`
      
 * Returns 
     * Returns the status on success.
     * Codes `201`/`40x`/`50x`
     * image or `error` field
