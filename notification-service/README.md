# Notification Service

Service for sending emails containing the finished image URL to the users.
Runs on port 8080.

The service takes the following env variables:

- SMTP_SERVER_URL: Required! The url where the smtp server is reachable. May contain port specification
- SMTP_SERVER_USERNAME: Optional. Only needed when the smtp server requires basicauth
- SMTP_SERVER_PASSWORD: Optional. Only needed when the smtp server requires basicauth

## Local setup

Run a sample smtp server for testing 

```sh
docker run -d -e "ServerOptions__Urls=http://*:80" -p 80:80 -p 2525:25 rnwood/smtp4dev
```
This will start the smtp server at port `2525` and expose the web dashboard at [localhost](localhost:8080).

<details>

<summary>With basic auth</summary>

```sh
docker run -d -e "ServerOptions__Urls=http://*:80" -e "RelayOptions__Login=water" -e "RelayOptions__Password=bottler" -e "ServerOptions__HostName=water-bottler-mail" -p 80:80 -p 2525:25 rnwood/smtp4dev
```

Run the notification service with envs:

- SMTP_SERVER_URL=localhost:2525 
- SMTP_SERVER_USERNAME=water 
- SMTP_SERVER_PASSWORD=bottler

</details>
