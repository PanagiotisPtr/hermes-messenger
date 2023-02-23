module github.com/panagiotisptr/hermes-messenger/services/authentication

go 1.18

require (
	github.com/go-redis/redis/v9 v9.0.0-rc.1
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/uuid v1.3.0
	github.com/panagiotisptr/hermes-messenger/protos v0.0.0
	go.uber.org/fx v1.18.2
	go.uber.org/zap v1.23.0
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa
	google.golang.org/grpc v1.51.0-dev
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/dig v1.15.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/net v0.0.0-20221014081412-f15817d10f9b // indirect
	golang.org/x/sys v0.0.0-20220908164124-27713097b956 // indirect
	golang.org/x/text v0.3.8 // indirect
	google.golang.org/genproto v0.0.0-20221024183307-1bc688fe9f3e // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

replace github.com/panagiotisptr/hermes-messenger/libs/service v0.0.0 => ../../libs/service

replace github.com/panagiotisptr/hermes-messenger/libs/utils v0.0.0 => ../../libs/utils

replace github.com/panagiotisptr/hermes-messenger/protos v0.0.0 => ../../protos
