# A simple authentication service written in Go
This services registers and authenticates users using JWTs

## Building and running the code
First, generate the GRPC code
```
make protos
```

Then build the go code
```
go build
```

Lastly, run the server
```
./authentication
```

## Availble commands and examples
The available server commands are
```
Authentication.Authenticate
Authentication.GetPublicKey
Authentication.Refresh
Authentication.Register
```

See the `protos/authentication.proto` for more details.

You can call these RPC methods either by using grpcurl or by generating a client for your language of choice using the `protos/authentication.proto` file.