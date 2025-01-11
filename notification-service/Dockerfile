# -- BUILD STAGE --
FROM golang:1.23.3-alpine3.19 AS build

WORKDIR /build

COPY . .

RUN go build -o notification-service

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
COPY --from=build --chown=$USER:$USER /build/notification-service /app/notification-service
COPY ./mail.tmpl ./mail.tmpl

# Ensure that binary can be run
RUN chmod +x /app/notification-service

USER $USER 

EXPOSE 8080

CMD ["/app/notification-service"]
