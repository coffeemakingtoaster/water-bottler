from ultralytics import YOLOWorld

if __name__ == "__main__":
    print("Starting Object Recognition Service")
    model = YOLOWorld("beer_detection_model.pt")

    images = [
        "example-images/image.jpeg",
        "example-images/image2.jpeg",
        "example-images/image3.jpeg",
        "example-images/image4.jpeg",
        "example-images/image5.jpeg",
        "example-images/image6.jpeg",
        "example-images/image7.jpeg",
        "example-images/image8.jpeg",
    ]

    for image in images:
        print("Processing image: ", image)
        results = model.predict(
            image,
            iou=0.3,
            conf=0.3,
            agnostic_nms=True,
            save=False,
        )
