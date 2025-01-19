FROM python:3.11-slim
# Install opencv dependencies via apk if needed

WORKDIR /app

COPY ./requirements.txt ./requirements.txt

# Force cpu install of torch toolkit before installing ultralytics
# The ultralytics install will then omit gpu dependencies -> reduce image size by around 7 GB
RUN pip install --no-cache-dir --upgrade pip && \
    pip install --no-cache-dir torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cpu && \
    pip install --no-cache-dir -r ./requirements.txt

COPY . .

# -u so it prints directly
CMD ["python","-u","main.py"]
