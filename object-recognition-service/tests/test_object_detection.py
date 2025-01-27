import unittest
from unittest.mock import Mock, patch
import numpy as np
from PIL import Image
from app.object_detection import BeerDetector

class TestBeerDetector(unittest.TestCase):
    @patch('app.object_detection.YOLOWorld')
    def test_initialization(self, mock_yolo):
        detector = BeerDetector()
        mock_yolo.assert_called_once()

        expected_classes = [
            f"{modifier} {cls}"
            for cls in ["beer bottle", "beer can", "beer glass", "beer mug"]
            for modifier in ["", "partial visible", "blurry", "empty"]
        ]
        mock_yolo.return_value.set_classes.assert_called_once_with(expected_classes)

    @patch('app.object_detection.YOLOWorld')
    def test_initialization_failure(self, mock_yolo):
        mock_yolo.side_effect = Exception("Model loading failed")

        with self.assertRaises(RuntimeError) as context:
            BeerDetector()
        self.assertIn("Failed to initialize YOLO model", str(context.exception))

    @patch('app.object_detection.YOLOWorld')
    def test_predict(self, mock_yolo):
        # Configure the complete mock chain
        mock_result = Mock()
        mock_boxes = Mock()
        mock_boxes.xyxy = Mock()
        mock_boxes.xyxy.cpu = Mock()
        mock_boxes.xyxy.cpu.return_value = Mock()
        mock_boxes.xyxy.cpu.return_value.numpy.return_value = np.array([
            [100, 100, 200, 200],
            [300, 300, 400, 400]
        ])
        mock_result.boxes = mock_boxes
        mock_yolo.return_value.predict.return_value = [mock_result]

        detector = BeerDetector()
        test_image = Image.new('RGB', (640, 480))
        result = detector.predict(test_image)
        self.assertEqual(result.shape, (2, 4))

    @patch('app.object_detection.YOLOWorld')
    def test_predict_failure(self, mock_yolo):
        mock_yolo.return_value.predict.side_effect = Exception("Prediction failed")
        detector = BeerDetector()
        test_image = Image.new('RGB', (640, 480))

        with self.assertRaises(RuntimeError):
            detector.predict(test_image)

if __name__ == '__main__':
    unittest.main()
