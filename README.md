# Hermes Messenger
Hermes-messeneger is an instant messaging application that's built using microservices in Go.

In addition to using microservices, this project also uses Consul service mesh to do things like parameter sharing, service discory etc. The services communicate with each other through gRPC.

This is still work in progress

## Services
- Authentication
  - Responsible for registering users
  - Authenticating users
  - Issuing tokens
  - Storing public/private RSA key pairs (this will move to Consul)
- Friends
  - Manages connections between users (friends)
- User
  - Manages user information
- Messaging
  - Sends and stores messages
- Chat App
  - Serves the front-end code
  - acts as an API gateway and proxies calls to other (internal) services

## Running a service
After installing all the dependencies for a service run
```
LISTEN_PORT=SOME_PORT HEALTH_CHECK_PORT=SOME_OTHER_PORT go run main.go
```
from inside the service's folder. You can also build the services or make Docker images for them (will be done eventually for all services)