# Upload Service

Service for uploading images to the system.
Runs on port 8081.

To specify the domain of the auth service used for api key validation use the `AUTH_SERVICE_URL` env variable.

## Example request

```sh
curl -X POST -F image=./testimage.png localhost:8080/upload
```

     Note: The file testimage.png must be present for this.
