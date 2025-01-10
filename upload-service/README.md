# Upload Service

Service for uploading images to the system.
Runs on port 8081.

To specify the domain of the auth service used for API key validation use the `AUTH_SERVICE_URL` env variable.
To specify the domain of the rabbitmq instance use the `QUEUE_URL` env variable.

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
curl -X POST -H "X-API-KEY: amVmZnMtd2F0ZXItYm90dGxlci1leGFtcGxlLWFwaS1rZXk=" -F image=./testimage.png localhost:8080/upload
```

     Note: The file testimage.png must be present for this.
