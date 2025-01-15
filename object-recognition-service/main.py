from ultralytics import YOLOWorld
from PIL import Image
import numpy as np
from image_utils import overlay_water_bottle

if __name__ == "__main__":
    print("Starting Object Recognition Service")
    model = YOLOWorld()

    classes = ["beer bottle", "beer can", "beer glass", "beer mug", "beer stein"]
    modifiers = ["partial visible", "blurry", "empty"]

    model.set_classes(
        classes
        + [" ".join([modifier, cls]) for cls in classes for modifier in modifiers]
    )

    results = model.predict(
        [
            "example-images/image.jpeg",
            "example-images/image2.jpeg",
            "example-images/image3.jpeg",
            "example-images/image4.jpeg",
            "example-images/image5.jpeg",
            "example-images/image6.jpeg",
            "example-images/image7.jpeg",
        ],
        iou=0.4,
        conf=0.5,
        agnostic_nms=True,
        save=False,
        project="./example-images",
        name="results",
        exist_ok=True,
    )

    for i, result in enumerate(results):
        img = result.path

        if len(result.boxes) > 0:
            boxes = result.boxes.xyxy.cpu().numpy()
            for box in boxes:
                x1, y1, x2, y2 = map(int, box)
                center_x = (x1 + x2) // 2
                center_y = (y1 + y2) // 2
                width = x2 - x1
                height = y2 - y1

                coords = np.array([x1, y1, x2, y2])
                img = overlay_water_bottle(img, coords)

        output_path = f"./example-images/results/image{i}_modified.jpg"
        final_image = Image.fromarray(img)
        if final_image.mode == 'RGBA':
            final_image = final_image.convert('RGB')
        final_image.save(output_path)
