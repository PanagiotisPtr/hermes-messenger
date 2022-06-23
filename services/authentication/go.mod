module github.com/panagiotisptr/hermes-messenger/services/authentication

go 1.18

require (
	github.com/go-redis/redis/v9 v9.0.0-beta.1
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/uuid v1.1.2
	github.com/panagiotisptr/hermes-messenger/libs/utils v0.0.0
	github.com/panagiotisptr/hermes-messenger/protos v0.0.0
	go.uber.org/fx v1.17.1
	go.uber.org/zap v1.21.0
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
	google.golang.org/grpc v1.47.0
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.7 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/dig v1.14.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f // indirect
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20200623002339-fbb79eadd5eb // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)

replace github.com/panagiotisptr/hermes-messenger/libs/service v0.0.0 => ../../libs/service

replace github.com/panagiotisptr/hermes-messenger/libs/utils v0.0.0 => ../../libs/utils

replace github.com/panagiotisptr/hermes-messenger/protos v0.0.0 => ../../protos
