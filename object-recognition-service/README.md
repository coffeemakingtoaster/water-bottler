# Object recognition service 

## Quickstart
To setup a new venv for development you can use `venv` command:
```sh
python3 -m venv  <path to whereever you store venvs>/water-bottler

# Activate it
source python3 -m venv  <path to whereever you store venvs>/water-bottler/bin/activate

pip install -r requirements.txt
```

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
