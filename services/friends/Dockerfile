FROM golang:1.19-alpine

RUN apk add --no-cache git
RUN apk add --no-cache protoc
RUN apk add --no-cache make
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR $GOPATH/github.com/panagiotisptr

RUN git clone https://github.com/PanagiotisPtr/hermes-messenger.git

WORKDIR $GOPATH/github.com/panagiotisptr/hermes-messenger

RUN make protos-go

WORKDIR $GOPATH/github.com/panagiotisptr/hermes-messenger/services/friends

RUN go mod tidy

RUN go build . 

ENTRYPOINT [ "./friends" ]
