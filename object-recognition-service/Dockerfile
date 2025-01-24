FROM python:3.11-slim

WORKDIR /app

COPY ./requirements.txt ./requirements.txt

RUN apt-get update && apt-get install -y ffmpeg libsm6 libxext6 git wget

# Force cpu install of torch toolkit before installing ultralytics
# The ultralytics install will then omit gpu dependencies -> reduce image size by around 7 GB
RUN pip install --no-cache-dir --upgrade pip && \
    pip install --no-cache-dir torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cpu && \
    pip install --no-cache-dir -r ./requirements.txt

# Install clip using git and then remove it -> we dont need it during runtime
RUN pip install --no-cache-dir "git+https://github.com/ultralytics/CLIP.git" && apt-get remove git -y

COPY . .

# Download the base model to reduce startup time
RUN wget https://github.com/ultralytics/assets/releases/download/v8.3.0/yolov8s-world.pt -O ./app/yolov8s-world.pt

WORKDIR /app/app
# -u so it prints directly
CMD ["python","-u","main.py"]