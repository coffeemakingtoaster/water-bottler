from dataclasses import dataclass

@dataclass
class Config:
    SCALING_FACTOR: float = 4.0
    IOU_THRESHOLD: float = 0.4
    CONFIDENCE_THRESHOLD: float = 0.5
    WATER_BOTTLE_PATH: str = 'example-images/water_bottle.png'
    OUTPUT_DIR: str = 'example-images/results'

CONFIG = Config()