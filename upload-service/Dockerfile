# -- BUILD STAGE --
FROM golang:1.23.3-alpine3.19 AS build

WORKDIR /build

COPY . .

RUN go build -o upload-service

# -- RUNTIME STAGE --
FROM alpine:3.19 AS runtime

WORKDIR /app

ARG USER=docker
ARG UID=5432
ARG GID=5433

# Create user for execution

#User group has same name as user
RUN addgroup -g $GID $USER 


RUN adduser \
    --disabled-password \
    --gecos "" \
    --ingroup "$USER" \
    --no-create-home \
	--uid "$UID" \
    "$USER"

# Copy build with permissions
COPY --from=build --chown=$USER:$USER /build/upload-service /app/upload-service

# Ensure that binary can be run
RUN chmod +x /app/upload-service

USER $USER 

EXPOSE 8080

CMD ["/app/upload-service"]
