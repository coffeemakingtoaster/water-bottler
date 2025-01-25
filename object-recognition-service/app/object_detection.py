import os
from ultralytics import YOLOWorld
from numpy import ndarray
from PIL import Image


class BeerDetector:
    """A class for detecting beer containers in images using YOLO model."""

    model: YOLOWorld

    def __init__(self):
        try:
            print("Initializing YOLO model...")
            model = YOLOWorld()
            print("Ich vermute hier sein")

            classes = ["beer bottle", "beer can", "beer glass", "beer mug"]
            modifiers = ["", "partial visible", "blurry", "empty"]

            model.set_classes(
                [f"{modifier} {cls}" for cls in classes for modifier in modifiers]
            )

            self.model = model
            print("Initialized YOLO model!")
        except Exception as e:
            raise RuntimeError(f"Failed to initialize YOLO model: {str(e)}")

    def predict(self, image: Image) -> ndarray:
        """
        Predicts the bounding boxes of potential beer containers in an image.

        Args:
            Image: The image to process.

        Returns:
            ndarray: The bounding box coordinates as numpy array of shape (n, 4) in [x1, y1, x2, y2] format

        Raises:
            RuntimeError: If prediction fails
        """

        try:
            result = self.model.predict(
                image,
                iou=0.4,
                conf=0.5,
                agnostic_nms=True,
                save=False,
                verbose=False,
            )
            boxes = result[0].boxes
            return boxes.xyxy.cpu().numpy()
        except Exception as e:
            raise RuntimeError(f"Failed to process image: {str(e)}")
