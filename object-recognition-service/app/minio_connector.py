import io
import os
from minio import Minio
from PIL import Image


class MinioConnector:
    def __init__(self):
        print("Connecting to Minio...")
        # Connect to Minio server
        try:
            self.client = Minio(
                os.environ.get("MINIO_ENDPOINT", "localhost:9000"),
                access_key=os.environ.get("MINIO_USER", "minio"),
                secret_key=os.environ.get("MINIO_KEY", "minio123"),
                secure=False,
            )

            self.bucket_name = os.environ.get("MINIO_BUCKET")

            # Make sure the bucket exists
            if not self.client.bucket_exists(self.bucket_name):
                self.client.make_bucket(self.bucket_name)
        except:
            raise RuntimeError("Could not connect to Minio server")

    def get_image(self, image_id) -> Image:
        response = self.client.get_object(self.bucket_name, image_id, f"{image_id}")
        responseData = io.BytesIO(response.read())
        return Image.open(responseData)

    def set_image(self, image_id, image: Image):
        image_bytes = io.BytesIO()
        image.save(image_bytes, format="PNG")
        image_bytes.seek(0)

        self.client.put_object(
            self.bucket_name,
            image_id,
            image_bytes,
            length=len(image_bytes.getvalue()),
            content_type="image/png",
        )
