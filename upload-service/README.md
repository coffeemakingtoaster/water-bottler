# Upload Service

Service for uploading images to the system.
Runs on port 8080.

## Example request

```sh
curl -X POST -F image=./testimage.png localhost:8080/upload
```

     Note: The file testimage.png must be present for this.
