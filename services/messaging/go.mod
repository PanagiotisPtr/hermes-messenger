module github.com/panagiotisptr/hermes-messenger/messaging

go 1.18

require (
	github.com/elastic/go-elasticsearch/v8 v8.4.0
	github.com/go-redis/redis/v9 v9.0.0-beta.2
	github.com/google/uuid v1.3.0
	github.com/panagiotisptr/hermes-messenger/libs/utils v0.0.0
	github.com/panagiotisptr/hermes-messenger/protos v0.0.0
	github.com/panagiotisptr/hermes-messenger/services/friends v0.0.0
	go.uber.org/fx v1.18.1
	go.uber.org/zap v1.22.0
	google.golang.org/grpc v1.46.2
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/elastic/elastic-transport-go/v8 v8.1.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	go.mongodb.org/mongo-driver v1.10.2 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/dig v1.15.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d // indirect
	golang.org/x/net v0.0.0-20220425223048-2871e0cb64e4 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20220422013727-9388b58f7150 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20200623002339-fbb79eadd5eb // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)

replace github.com/panagiotisptr/hermes-messenger/libs/service v0.0.0 => ../../libs/service

replace github.com/panagiotisptr/hermes-messenger/libs/utils v0.0.0 => ../../libs/utils

replace github.com/panagiotisptr/hermes-messenger/protos v0.0.0 => ../../protos

replace github.com/panagiotisptr/hermes-messenger/services/friends v0.0.0 => ../friends
