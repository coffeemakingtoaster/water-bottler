from PIL import Image
import numpy as np
import os
from typing import Union
from config import CONFIG

def process_image_paths(image_paths: list) -> list:
    valid_paths = []
    for path in image_paths:
        if not os.path.exists(path):
            print(f"Warning: Image not found: {path}")
            continue
        valid_paths.append(path)
    return valid_paths

def load_water_bottle(water_bottle_path: str) -> Image:
    if not os.path.exists(water_bottle_path):
        raise FileNotFoundError(f"Water bottle image not found at {water_bottle_path}")
    return Image.open(water_bottle_path).convert('RGBA')

def overlay_water_bottle(background_image: Union[str, np.ndarray], water_bottle_image: Image,
                        box_coordinates: np.ndarray) -> np.ndarray:
    if isinstance(background_image, str):
        background = Image.open(background_image).convert('RGBA')
    else:
        background = Image.fromarray(background_image).convert('RGBA')

    x1, y1, x2, y2 = map(int, box_coordinates)
    base_width = x2 - x1
    base_height = y2 - y1

    target_width = int(base_width * CONFIG.SCALING_FACTOR)
    target_height = int(base_height * CONFIG.SCALING_FACTOR)

    aspect_ratio = water_bottle_image.width / water_bottle_image.height
    if target_width / target_height > aspect_ratio:
        target_width = int(target_height * aspect_ratio)
    else:
        target_height = int(target_width / aspect_ratio)

    water_bottle_resized = water_bottle_image.resize(
        (target_width, target_height),
        Image.Resampling.LANCZOS
    )

    paste_x = x1 + (x2 - x1 - target_width) // 2
    paste_y = y1 + (y2 - y1 - target_height) // 2
    background.paste(water_bottle_resized, (paste_x, paste_y), water_bottle_resized)

    return np.array(background)