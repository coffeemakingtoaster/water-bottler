import pytest
from unittest.mock import Mock, MagicMock, patch
import os
from PIL import Image
import numpy as np
from main import process_images

@pytest.fixture
def mock_config():
    config = MagicMock()
    config.WATER_BOTTLE_PATH = "example-images/water_bottle.png"
    config.OUTPUT_DIR = "example-images/results"
    config.SCALING_FACTOR = 4.0
    config.IOU_THRESHOLD = 0.4
    config.CONFIDENCE_THRESHOLD = 0.5
    return config

@pytest.fixture
def mock_model():
    model = MagicMock()
    return model

@pytest.fixture
def mock_results():
    result = MagicMock()
    result.boxes = MagicMock()
    result.boxes.xyxy.cpu.return_value = np.array([[10, 10, 100, 100]])
    result.path = np.zeros((224, 224, 3), dtype=np.uint8)
    return [result]

@pytest.fixture
def mock_water_bottle():
    return np.zeros((50, 50, 4), dtype=np.uint8)

def test_process_images_success(mock_config, mock_model, mock_results, mock_water_bottle):
    with patch('config.CONFIG', mock_config), \
         patch('main.setup_model', return_value=mock_model), \
         patch('image_utils.load_water_bottle', return_value=mock_water_bottle), \
         patch('image_utils.overlay_water_bottle', return_value=np.zeros((224, 224, 3))):

        mock_model.predict.return_value = mock_results
        image_paths = ['image1.jpg', 'image2.jpg']

        process_images(image_paths)

        mock_model.predict.assert_called_once_with(
            image_paths,
            iou=mock_config.IOU_THRESHOLD,
            conf=mock_config.CONFIDENCE_THRESHOLD,
            agnostic_nms=True,
            save=False,
            project=mock_config.OUTPUT_DIR,
            name="results",
            exist_ok=True
        )

        assert os.path.exists(mock_config.OUTPUT_DIR)

def test_process_images_empty_input():
    with pytest.raises(ValueError):
        process_images([])

def test_process_images_invalid_path(mock_config, mock_model):
    with patch('config.CONFIG', mock_config), \
         patch('main.setup_model', return_value=mock_model):

        mock_model.predict.side_effect = Exception("Invalid path")

        with pytest.raises(Exception):
            process_images(['invalid_path.jpg'])

def test_process_images_rgba_conversion(mock_config, mock_model, mock_water_bottle):
    rgba_image = np.zeros((224, 224, 4), dtype=np.uint8)
    result = MagicMock()
    result.boxes.xyxy.cpu.return_value = np.array([[10, 10, 100, 100]])
    result.path = rgba_image

    with patch('config.CONFIG', mock_config), \
         patch('main.setup_model', return_value=mock_model), \
         patch('image_utils.load_water_bottle', return_value=mock_water_bottle), \
         patch('image_utils.overlay_water_bottle', return_value=rgba_image):

        mock_model.predict.return_value = [result]
        process_images(['image.jpg'])

        output_path = os.path.join(mock_config.OUTPUT_DIR, 'image0_modified.jpg')
        saved_image = Image.open(output_path)
        assert saved_image.mode == 'RGB'
