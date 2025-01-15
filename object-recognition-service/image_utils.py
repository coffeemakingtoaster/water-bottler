from PIL import Image
import numpy as np

def overlay_water_bottle(background_image, box_coordinates):
    # Handle both file path and numpy array inputs
    if isinstance(background_image, str):
        background = Image.open(background_image).convert('RGBA')
    else:
        background = Image.fromarray(background_image).convert('RGBA')

    water_bottle_path = "example-images/water_bottle.png"
    water_bottle = Image.open(water_bottle_path).convert('RGBA')

    x1, y1, x2, y2 = map(int, box_coordinates)
    target_width = x2 - x1
    target_height = y2 - y1

    water_bottle = water_bottle.resize((target_width, target_height), Image.Resampling.LANCZOS)
    background.paste(water_bottle, (x1, y1), water_bottle)

    return np.array(background)
