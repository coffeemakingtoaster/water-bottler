from PIL import Image
import numpy as np

def overlay_water_bottle(background_image, box_coordinates):
    if isinstance(background_image, str):
        background = Image.open(background_image).convert('RGBA')
    else:
        background = Image.fromarray(background_image).convert('RGBA')

    water_bottle_path = "example-images/water_bottle.png"
    water_bottle = Image.open(water_bottle_path).convert('RGBA')

    x1, y1, x2, y2 = map(int, box_coordinates)
    base_width = (x2 - x1)
    base_height = (y2 - y1)

    scaling_factor = 4
    target_width = int(base_width * scaling_factor)
    target_height = int(base_height * scaling_factor)

    aspect_ratio = water_bottle.width / water_bottle.height
    if target_width / target_height > aspect_ratio:
        target_width = int(target_height * aspect_ratio)
    else:
        target_height = int(target_width / aspect_ratio)

    water_bottle_resized = water_bottle.resize((target_width, target_height), Image.Resampling.LANCZOS)

    paste_x = x1 + (x2 - x1 - target_width) // 2
    paste_y = y1 + (y2 - y1 - target_height) // 2

    background.paste(water_bottle_resized, (paste_x, paste_y), water_bottle_resized)

    return np.array(background)