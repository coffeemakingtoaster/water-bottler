import os
from typing import Tuple

from ultralytics import YOLOWorld
from numpy import ndarray
from config import CONFIG


class BeerDetector:
    """A class for detecting beer containers in images using YOLO model."""

    model: YOLOWorld

    def __init__(self):
        try:
            model = YOLOWorld()

            classes = ["beer bottle", "beer can", "beer glass", "beer mug"]
            modifiers = ["", "partial visible", "blurry", "empty"]

            model.set_classes(
                [f"{modifier} {cls}" for cls in classes for modifier in modifiers]
            )

            self.model = model
        except Exception as e:
            raise RuntimeError(f"Failed to initialize YOLO model: {str(e)}")

    def predict(self, image_path: str) -> Tuple[ndarray, ndarray]:
        """
        Predicts the bounding boxes of potential beer containers in an image.

        Args:
            image_path (str): The path to the image file.

        Returns:
            Tuple[ndarray, ndarray]: A tuple containing:
                - The bounding box coordinates as numpy array of shape (n, 4) in [x1, y1, x2, y2] format
                - The confidence scores as numpy array of shape (n,)

        Raises:
            FileNotFoundError: If the image file does not exist
            RuntimeError: If prediction fails
        """
        if not os.path.exists(image_path):
            raise FileNotFoundError(f"No Image found at: {image_path}")

        try:
            result = self.model.predict(
                image_path,
                iou=CONFIG.IOU_THRESHOLD,
                conf=CONFIG.CONFIDENCE_THRESHOLD,
                agnostic_nms=True,
                save=False,
                verbose=False,
            )
            boxes = result[0].boxes
            return boxes.xyxy.cpu().numpy(), boxes.conf.cpu().numpy()
        except Exception as e:
            raise RuntimeError(f"Failed to process image: {str(e)}")
