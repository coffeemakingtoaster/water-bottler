FROM python:3.11-slim
# Install opencv dependencies via apk if needed

WORKDIR /app

COPY ./requirements.txt ./requirements.txt

RUN apt-get update && apt-get install ffmpeg libsm6 libxext6 git -y

# Force cpu install of torch toolkit before installing ultralytics
# The ultralytics install will then omit gpu dependencies -> reduce image size by around 7 GB
RUN pip install --no-cache-dir --upgrade pip && \
    pip install --no-cache-dir torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cpu && \
    pip install --no-cache-dir -r ./requirements.txt

# Install clip using git and then remove it -> we dont need it during runtime
RUN pip install --no-cache-dir "git+https://github.com/ultralytics/CLIP.git" && apt-get remove git -y

COPY . .

# Load model into dockerfile to reduce startup time
RUN python -u ./setup-model.py 

# -u so it prints directly
CMD ["python","-u","main.py"]
