import json
from object_detection import BeerDetector
from image_processing import WaterBottleImageProcessor
from rabbitmq_connector import RabbitMQConnector
from minio_connector import MinioConnector


def onImageEventReceived(ch, method, properties, body):
    # Read the payload
    payload = json.loads(body)
    image_id = payload["image_id"]
    email = payload["user_mail"]

    # Get the image from Minio
    image = minio.get_image(image_id)

    # Predict the bounding boxes of potential beer containers in the image
    boxes = beer_detection_model.predict(image)

    # Process the image by overlaying a water bottle on top of the detected beer containers
    edited_image = water_bottle_processor.process(image, boxes)

    # Save the edited image back to Minio
    minio.set_image(image_id, edited_image)

    # Acknowledge the message and publish a task finish event
    queue_connector.publish_task_finish_event(image_id, email)
    ch.basic_ack(delivery_tag=method.delivery_tag)


if __name__ == "__main__":
    print("Starting Object Recognition Service")

    # Setup Image processing classes
    beer_detection_model = BeerDetector()
    water_bottle_processor = WaterBottleImageProcessor("water_bottle.png")

    # Setup Minio connection
    minio = MinioConnector()

    # Setup RabbitMQ connections
    queue_connector = RabbitMQConnector()
    queue_connector.register_callback(
        queue="image-workload",
        callback=onImageEventReceived,
    )

    print("Waiting for messages... To exit press CTRL+C")
    queue_connector.start_listening()
