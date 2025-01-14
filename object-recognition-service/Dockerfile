FROM python:3.12.8-alpine3.21

# Install opencv dependencies via apk if needed

WORKDIR /app

COPY ./requirements.txt ./requirements.txt

RUN pip install -r ./requirements.txt

COPY . .

# -u so it prints directly
CMD ["python","-u","main.py"]
