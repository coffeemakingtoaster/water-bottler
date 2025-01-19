from ultralytics import YOLOWorld

if __name__ != "__main__":
    exit(1)
model = YOLOWorld()
classes = ["beer bottle", "beer can", "beer glass", "beer mug", "beer stein"]
modifiers = ["partial visible", "blurry", "empty"]
all_classes = classes + [
    f"{modifier} {cls}"
    for cls in classes
    for modifier in modifiers
]
model.set_classes(all_classes)
