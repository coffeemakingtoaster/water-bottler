import os
import json
from object_detection import BeerDetector
from image_processing import WaterBottleImageProcessor
from rabbitmq_connector import RabbitMQConnector
from minio_connector import MinioConnector

SLOW_MODE_DELAY = os.getenv("SLOW_MODE_DELAY", 0)


def onImageEventReceived(ch, method, properties, body):
    # Read the payload
    try:
        payload = json.loads(body)
        image_id = payload["image_id"]
        email = payload["user_mail"]
    except:
        print("ERROR: Invalid message payload")
        # TODO: How to handle messages with invalid payload?
        ch.basic_ack(delivery_tag=method.delivery_tag)
        return

    # Get the image from Minio
    try:
        print(f"Getting image from Minio with id {image_id}")
        image = minio.get_image(image_id)
    except:
        print("ERROR: Could not get image from Minio")
        return

    # Predict the bounding boxes of potential beer containers in the image
    boxes = beer_detection_model.predict(image)
    print("Detected beer containers:", boxes, "in image", image_id)

    # Slow down the processing to simulate a slow service for presentation
    # purposes
    if SLOW_MODE_DELAY > 0:
        import time

        time.sleep(SLOW_MODE_DELAY)

    # Process the image by overlaying a water bottle on top of the detected beer containers
    edited_image = water_bottle_processor.process(image, boxes)

    # Save the edited image back to Minio
    try:
        minio.set_image(image_id, edited_image)
    except:
        print("ERROR: Could not save edited image to Minio")
        return

    # Acknowledge the message and publish a task finish event
    try:
        queue_connector.publish_task_finish_event(image_id, email)
    except:
        print("ERROR: Could not publish task finish event")
        return

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
