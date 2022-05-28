FROM golang:1.18-alpine

RUN apk add --no-cache git

WORKDIR $GOPATH/service/

COPY . .

RUN go mod tidy

RUN go build .

ENTRYPOINT [ "./authentication" ]
