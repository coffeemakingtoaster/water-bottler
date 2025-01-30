# Object recognition service

## Environment Variables
```
MINIO_USER: <MinIO Username>
MINIO_KEY: <MinIO Password>
MINIO_ENDPOINT: <Minio Host (with Port)>
MINIO_BUCKET: <Minio Bucket to read/storage images>
QUEUE_HOST: <RabbitMQ Queue>
QUEUE_USER: <RabbitMQ User>
QUEUE_PASS: <RabbitMQ Password>
QUEUE_INPUT_NAME: <RabbitMQ Queue name to get image task events>
QUEUE_OUTPUT_NAME: <RabbitMQ Queue name to put image done events>
SLOW_MODE_DELAY: <Number of seconds to sleep between image detection -> image procesing>
```
