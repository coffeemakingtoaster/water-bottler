import os
import numpy as np

from PIL import Image
from typing import List

class WaterBottleImageProcessor:
    def __init__(self, water_bottle_path: str, water_glass_path: str):
        if not os.path.exists(water_bottle_path):
            raise FileNotFoundError(
                f"Did not find water bottle image at {water_bottle_path}"
            )
        if not os.path.exists(water_glass_path):
            raise FileNotFoundError(
                f"Did not find water glass image at {water_glass_path}"
            )

        self.water_bottle = Image.open(water_bottle_path).convert("RGBA")
        self.water_glass = Image.open(water_glass_path).convert("RGBA")

    def process(
        self,
        image: Image,
        box_coordinates: List[np.ndarray],
    ) -> Image.Image:
        """
        Overlays the water bottle image on top of given image at the provided box_coordinates.

        Args:
            image (Image): The image to overlay the water bottle on.
            box_coordinates (List[np.ndarray]): A list of numpy arrays containing the bounding box coordinates of the detected beer containers.
        """
        # Read in the image as PIL Image
        image = image.convert("RGBA")

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

            # Calculate the width and height of the bounding box
            width = x2 - x1
            height = y2 - y1

            # Scaling the water bottle image to fit the bounding box
            # but keeping the aspect ratio
            aspect_ratio_bottle = self.water_bottle.height / self.water_bottle.width
            aspect_ratio_glass = self.water_glass.height / self.water_glass.width
            aspect_ratio_target = height / width

            #choose image with ratio that is closer to the target
            fitting_replacement = self.water_glass if abs(aspect_ratio_glass - aspect_ratio_target) < abs(aspect_ratio_bottle - aspect_ratio_target) else self.water_bottle
            aspect_ratio_fitting_replacement = aspect_ratio_glass if fitting_replacement == self.water_glass else aspect_ratio_bottle

            if(aspect_ratio_fitting_replacement > aspect_ratio_target):
                scale_factor = width / fitting_replacement.width
                target_width = width
                target_height = int(fitting_replacement.height * scale_factor)
            else:
                scale_factor = height / fitting_replacement.height
                target_height = height
                target_width = int(fitting_replacement.width * scale_factor)

            # Resize bottle to fit the bounding box
            resized_replacement = fitting_replacement.resize(
                (target_width, target_height),
                Image.Resampling.LANCZOS,
            )

            paste_x = x1 + (x2 - x1 - target_width) // 2
            paste_y = y1 + (y2 - y1 - target_height) // 2

            # Paste the bottle on the image using the alpha channel as mask
            image.paste(resized_replacement, (paste_x, paste_y), resized_replacement)

        return image.convert("RGB")
