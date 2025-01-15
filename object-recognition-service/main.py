from ultralytics import YOLOWorld
from PIL import Image
import numpy as np
import os
from config import CONFIG
from image_utils import overlay_water_bottle
from image_utils import process_image_paths

def setup_model():
    model = YOLOWorld()
    classes = ["beer bottle", "beer can", "beer glass", "beer mug", "beer stein"]
    modifiers = ["partial visible", "blurry", "empty"]
    all_classes = classes + [
        f"{modifier} {cls}"
        for cls in classes
        for modifier in modifiers
    ]
    model.set_classes(all_classes)
    return model

def process_images(image_paths):
    try:
        os.makedirs(CONFIG['OUTPUT_DIR'], exist_ok=True)
        model = setup_model()

        results = model.predict(
            image_paths,
            iou=CONFIG['IOU_THRESHOLD'],
            conf=CONFIG['CONFIDENCE_THRESHOLD'],
            agnostic_nms=True,
            save=False,
            project=CONFIG['OUTPUT_DIR'],
            name="results",
            exist_ok=True,
        )

        for i, result in enumerate(results):
            if len(result.boxes) > 0:
                img = result.path
                boxes = result.boxes.xyxy.cpu().numpy()

                for box in boxes:
                    img = overlay_water_bottle(img, box)

                output_path = os.path.join(
                    CONFIG['OUTPUT_DIR'],
                    f"image{i}_modified.jpg"
                )
                final_image = Image.fromarray(img)
                if final_image.mode == 'RGBA':
                    final_image = final_image.convert('RGB')
                final_image.save(output_path)
                print(f"Saved processed image to {output_path}")

    except Exception as e:
        print(f"Error in image processing pipeline: {str(e)}")
        raise

if __name__ == "__main__":
    print("Starting Object Recognition Service")

    image_paths = [
        "example-images/image.jpeg",
        "example-images/image2.jpeg",
        "example-images/image3.jpeg",
        "example-images/image4.jpeg",
        "example-images/image5.jpeg",
        "example-images/image6.jpeg",
        "example-images/image7.jpeg",
    ]

    process_images(process_image_paths(image_paths))
