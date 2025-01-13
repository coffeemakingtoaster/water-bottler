from ultralytics import YOLOWorld

if __name__ == "__main__":
    print("Starting Object Recognition Service")
    model = YOLOWorld()

    classes = ["beer bottle", "beer can", "beer glass", "beer mug", "beer stein"]
    modifiers = ["partial visible", "blurry", "empty"]

    model.set_classes(
        classes
        + [" ".join([modifier, cls]) for cls in classes for modifier in modifiers]
    )

    results = model.predict(
        [
            "example-images/image.jpeg",
            "example-images/image2.jpeg",
            "example-images/image3.jpeg",
            "example-images/image4.jpeg",
            "example-images/image5.jpeg",
            "example-images/image6.jpeg",
            "example-images/image7.jpeg",
        ],
        iou=0.4,
        conf=0.5,
        agnostic_nms=True,
        save=True,
        project="./example-images",
        name="results",
        exist_ok=True,
    )
