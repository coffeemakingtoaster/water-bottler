from dataclasses import dataclass


@dataclass
class Config:
    SCALE_FACTOR: float = 1.5
    IOU_THRESHOLD: float = 0.4
    CONFIDENCE_THRESHOLD: float = 0.5
    WATER_BOTTLE_PATH: str = "water_bottle_cropped.png"
    INPUT_DIR: str = "../example-data/example-images"
    OUTPUT_DIR: str = "../example-data/example-results"


CONFIG = Config()
