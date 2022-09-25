.PHONY: protos-go protos-ts

# Generate Go clients
protos-go:
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	protos/*.proto

# Generate TypeScript clients (for NodeJS)
protos-ts:
	./services/chat-client/node_modules/grpc-tools/bin/protoc \
	--plugin=services/chat-client/node_modules/.bin/protoc-gen-ts_proto \
	--ts_proto_opt=outputServices=grpc-js,env=node,useOptionals=messages,exportCommonSymbols=false,esModuleInterop=true \
	--ts_proto_out=services/chat-client/grpc-clients \
	--proto_path protos protos/*.proto
