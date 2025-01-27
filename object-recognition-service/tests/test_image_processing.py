import unittest
from unittest.mock import patch, MagicMock
from PIL import Image
import numpy as np
from app.image_processing import WaterBottleImageProcessor

class TestWaterBottleImageProcessor(unittest.TestCase):
    def setUp(self):
        self.mock_path_bottle = "water_bottle.png"
        self.mock_path_glass = "water_glass.png"
        self.mock_water_bottle = Image.new("RGBA", (100, 200), (255, 0, 0, 255))
        self.mock_water_glass = Image.new("RGBA", (100, 100), (0, 0, 255, 255))

        self.patcher_exists = patch("os.path.exists", return_value=True)

        def mock_open_side_effect(path):
            if path == self.mock_path_bottle:
                return self.mock_water_bottle
            elif path == self.mock_path_glass:
                return self.mock_water_glass
            return None

        self.patcher_open = patch("PIL.Image.open", side_effect=mock_open_side_effect)

        self.mock_exists = self.patcher_exists.start()
        self.mock_open = self.patcher_open.start()

        self.processor = WaterBottleImageProcessor(self.mock_path_bottle, self.mock_path_glass)

    def tearDown(self):
        self.patcher_exists.stop()
        self.patcher_open.stop()

    def test_initialization_valid_path(self):
        self.assertEqual(self.processor.water_bottle.mode, "RGBA")
        self.assertEqual(self.processor.water_bottle.size, (100, 200))
        self.assertEqual(self.processor.water_glass.mode, "RGBA")
        self.assertEqual(self.processor.water_glass.size, (100, 100))

    def test_initialization_invalid_bottle_path(self):
        with patch("os.path.exists", return_value=False):
            with self.assertRaises(FileNotFoundError) as context:
                WaterBottleImageProcessor("invalid_bottle_path.png", self.mock_path_glass)
            self.assertIn("Did not find water bottle image at invalid_bottle_path.png", str(context.exception))

    def test_initialization_invalid_glass_path(self):
        with patch("os.path.exists") as mock_exists:
            mock_exists.side_effect = lambda path: path == self.mock_path_bottle
            with self.assertRaises(FileNotFoundError) as context:
                WaterBottleImageProcessor(self.mock_path_bottle, "invalid_glass_path.png")
            self.assertIn("Did not find water glass image at invalid_glass_path.png", str(context.exception))

    def test_process_empty_box_list(self):
        test_image = Image.new("RGBA", (500, 500), (0, 255, 0, 255))
        original_pixels = test_image.load()
        result = self.processor.process(test_image, [])

        # Verify image remains unchanged
        self.assertEqual(result.size, (500, 500))
        self.assertEqual(result.mode, "RGB")
        result_pixels = result.load()
        self.assertEqual(result_pixels[250, 250][:3], (0, 255, 0))

    def test_process_single_box(self):
        test_image = Image.new("RGBA", (500, 500), (0, 255, 0, 255))
        box = np.array([100, 100, 200, 300])
        result = self.processor.process(test_image, [box])
        center_x, center_y = 150, 200

        # Verify water bottle presence at center
        self.assertNotEqual(result.getpixel((center_x, center_y)), (0, 255, 0))
        outside_points = [(50, 50), (350, 350), (50, 350), (350, 50)]
        for point in outside_points:
            self.assertEqual(result.getpixel(point), (0, 255, 0))

    def test_process_multiple_boxes_with_aspect_ratio(self):
        test_image = Image.new("RGBA", (500, 500), (0, 255, 0, 255))
        boxes = [
            np.array([50, 50, 150, 250]),    # Tall box
            np.array([200, 200, 400, 300])   # Wide box
        ]
        result = self.processor.process(test_image, boxes)
        tall_center = (100, 200)
        wide_center = (300, 250)

        # Both centers should have overlay
        self.assertNotEqual(result.getpixel(tall_center), (0, 255, 0))
        self.assertNotEqual(result.getpixel(wide_center), (0, 255, 0))
        outside_points = [(25, 25), (450, 450)]
        for point in outside_points:
            self.assertEqual(result.getpixel(point), (0, 255, 0))

    def test_box_sorting_and_overlap(self):
        test_image = Image.new("RGBA", (500, 500), (0, 255, 0, 255))
        boxes = [
            np.array([0, 0, 400, 400]),     # Large box
            np.array([50, 50, 150, 150]),   # Small box overlapping
            np.array([300, 300, 350, 350])  # Small isolated box
        ]
        result = self.processor.process(test_image, boxes)
        centers = [(200, 200), (100, 100), (325, 325)]

        for center in centers:
            self.assertNotEqual(result.getpixel(center), (0, 255, 0))
        # Verify processing order (smaller boxes should be processed first)
        overlap_point = (75, 75)
        self.assertNotEqual(result.getpixel(overlap_point), (0, 255, 0))

    def test_correct_image_used_bottle(self):
        test_image = Image.new("RGBA", (500, 500), (0, 255, 0, 255))
        box = np.array([100, 100, 200, 300])  # Tall box, should use water bottle
        result = self.processor.process(test_image, [box])
        center_x, center_y = 150, 200

        # Verify that the pixel at the center matches the water bottle color
        self.assertEqual(result.getpixel((center_x, center_y)), (255, 0, 0))

    def test_correct_image_used_glass(self):
        test_image = Image.new("RGBA", (500, 500), (0, 255, 0, 255))
        box = np.array([100, 100, 200, 150])  # Wide box, should use water glass
        result = self.processor.process(test_image, [box])
        center_x, center_y = 150, 125

        # Verify that the pixel at the center matches the water glass color
        self.assertEqual(result.getpixel((center_x, center_y)), (0, 0, 255))

if __name__ == "__main__":
    unittest.main()
