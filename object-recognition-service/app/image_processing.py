import os
import numpy as np

from PIL import Image
from typing import List

IMAGE_SCALE_FACTOR = 1.5


class WaterBottleImageProcessor:
    def __init__(self, water_bottle_path: str):
        if not os.path.exists(water_bottle_path):
            raise FileNotFoundError(
                f"Did not find water bottle image at {water_bottle_path}"
            )

        self.water_bottle = Image.open(water_bottle_path).convert("RGBA")

    def process(
        self,
        image_path: str,
        box_coordinates: List[np.ndarray],
    ) -> Image.Image:
        """
        Overlays the water bottle image on top of given image at the provided box_coordinates.

        Args:
            image (str): The path of the image to overlay the water bottle on.
            box_coordinates (List[np.ndarray]): A list of numpy arrays containing the bounding box coordinates of the detected beer containers.
        """
        # Check if image exists
        if not os.path.exists(image_path):
            raise FileNotFoundError(f"Image not found at {image_path}")

        # Read in the image as PIL Image
        image = Image.open(image_path).convert("RGBA")

        # Sort the bounding boxes by size
        # This is done to ensure that smaller bounding boxes are processed first
        # since smaller boxes are more likely to be "behind" larger boxes
        box_coordinates = sorted(
            box_coordinates,
            key=lambda x: (x[2] - x[0]) * (x[3] - x[1]),
        )

        # Iterate over all bounding boxes and overlay the water bottle image
        for coords in box_coordinates:
            # Convert to integers
            x1, y1, x2, y2 = map(int, coords)

            # Calculate the  width and height of the bounding box
            width = x2 - x1
            height = y2 - y1

            # Scaling the water bottle image to fit the bounding box
            # but keeping the aspect ratio
            target_width = int(width * IMAGE_SCALE_FACTOR)
            target_height = int(height * IMAGE_SCALE_FACTOR)

            aspect_ratio = self.water_bottle.width / self.water_bottle.height
            if target_width / target_height > aspect_ratio:
                target_width = int(target_height * aspect_ratio)
            else:
                target_height = int(target_width / aspect_ratio)

            # Resize bottle to fit the bounding box
            resized_water_bottle = self.water_bottle.resize(
                (target_width, target_height),
                Image.Resampling.LANCZOS,
            )

            paste_x = x1 + (x2 - x1 - target_width) // 2
            paste_y = y1 + (y2 - y1 - target_height) // 2

            # Paste the bottle on the image using the alpha channel as mask
            image.paste(resized_water_bottle, (paste_x, paste_y), resized_water_bottle)

        return image.convert("RGB")
