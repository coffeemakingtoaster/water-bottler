# Notification Service

Service for sending emails containing the finished image URL to the users.
Runs on port 8080.

Run a sample smtp server for testing 

```sh
docker run -d -e "ServerOptions__Urls=http://*:80"  -p 80:80 -p 2525:25 rnwood/smtp4dev
```
This will start the smtp server at port `2525` and expose the web dashboard at [localhost](localhost:8080).
