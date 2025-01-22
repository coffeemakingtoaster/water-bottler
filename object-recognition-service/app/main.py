import os
from config import CONFIG
from object_detection import BeerDetector
from image_processing import WaterBottleImageProcessor

if __name__ == "__main__":
    print("Starting Object Recognition Service")

    # Initialize needed classes
    beer_detection_model = BeerDetector()
    water_bottle_processor = WaterBottleImageProcessor(CONFIG.WATER_BOTTLE_PATH)

    # Load images in the example folder...
    # This will later be replaced with a
    # listener to the rabbitmq queue to retrieve images
    # from the image storage solution
    print(f"Loading images from {CONFIG.INPUT_DIR}")

    image_paths = [
        os.path.join(CONFIG.INPUT_DIR, image_name)
        for image_name in os.listdir(CONFIG.INPUT_DIR)
        if image_name.endswith((".jpeg", ".jpg", ".png"))
    ]

    print(f"Found {len(image_paths)} images")

    # Make sure the output path exists...
    os.makedirs(CONFIG.OUTPUT_DIR, exist_ok=True)

    for image_path in image_paths:
        print(f"Processing Image {image_path}")
        boxes, conf = beer_detection_model.predict(image_path)
        high_conf_boxes = boxes[conf > CONFIG.CONFIDENCE_THRESHOLD]
        print(f"Found {len(high_conf_boxes)} high confidence beer containers")
        edited_image = water_bottle_processor.process(image_path, high_conf_boxes)
        edited_image.save(CONFIG.OUTPUT_DIR + "/" + os.path.basename(image_path))

    print(f"Saved processed images to {CONFIG.OUTPUT_DIR}")
