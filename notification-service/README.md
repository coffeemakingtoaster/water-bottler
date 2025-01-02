# Notification Service

Service for sending emails containing the finished image URL to the users.
Runs on port 8080.

Run a sample smtp server for testing 

```sh
docker run -d -p 1025:1025 -p 1080:1080 haravich/fake-smtp-server
```
