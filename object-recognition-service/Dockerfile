FROM python:3.11-slim
# Install opencv dependencies via apk if needed

WORKDIR /app

COPY ./requirements.txt ./requirements.txt

RUN pip install -r ./requirements.txt

COPY . .

# -u so it prints directly
CMD ["python","-u","main.py"]
