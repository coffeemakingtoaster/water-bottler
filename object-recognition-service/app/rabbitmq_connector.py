import pika
import json
import os


class RabbitMQConnector:
    def __init__(self):
        print("Connecting to RabbitMQ...")
        try:
            self.connection = pika.BlockingConnection(
                pika.ConnectionParameters(
                    os.environ.get("QUEUE_HOST", "localhost"),
                    credentials=pika.PlainCredentials(
                        os.environ.get("QUEUE_USER", "water"),
                        os.environ.get("QUEUE_PASS", "bottler"),
                    ),
                )
            )

            self.channel = self.connection.channel()
        except Exception as e:
            raise RuntimeError(f"Could not connect to RabbitMQ server: {str(e)}")

    def register_callback(self, queue, callback):
        # Make sure the queue exists
        self.channel.queue_declare(queue=queue)

        self.channel.basic_consume(
            queue=queue,
            on_message_callback=callback,
        )

    def publish_task_finish_event(self, image_id, email):
        # Make sure the queue exists
        queue = os.environ.get("QUEUE_OUTPUT_NAME", "image-workload")
        self.channel.queue_declare(queue=queue)

        # Publish the done event
        self.channel.basic_publish(
            exchange="",
            routing_key=queue,
            body=json.dumps({"image_id": image_id, "email": email}),
        )

    def start_listening(self):
        self.channel.start_consuming()
