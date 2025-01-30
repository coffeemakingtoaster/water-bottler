# Upload Service

Service for uploading images to the system.
Runs on port 8081.

The following evnironment variables should be set
- AUTH_SERVICE_URL 
- QUEUE_URL
- MINIO_ACCESS_KEY
- MINIO_SECRET_KEY
- MINIO_ENDPOINT
- MINIO_BUCKET_NAME

<details>

<summary> Start a rabbitmq instance locally </summary>

To start a rabbitmq instance locally:
```sh
docker run --rm -e RABBITMQ_DEFAULT_USER=water -e RABBITMQ_DEFAULT_PASS=bottler -p 3111:15672 -p 5672:5672 rabbitmq:3-management-alpine
```

You can now see the management dashboard at localhost:3111

</details>

## Example request
```sh
curl -X POST -F image=@docker-compose.yaml -H "X-API-KEY: amVmZnMtd2F0ZXItYm90dGxlci1leGFtcGxlLWFwaS1rZXk=" localhost:8081/upload -s -o /dev/null -w "%{http_code}"
```
