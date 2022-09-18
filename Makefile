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

deploy-service-user:
	kubectl apply -f services/user/deployment/elasticsearch.yml -n hermes-messenger; \
	helm install user-redis bitnami/redis -n hermes-messenger; \
	kubectl apply -f services/user/deployment/service.yml -n hermes-messenger

deploy-service-friends:
	helm install friends-redis bitnami/redis -n hermes-messenger; \
	kubectl apply -f services/friends/deployment/service.yml -n hermes-messenger

deploy-service-messaging:
	kubectl apply -f services/messaging/deployment/elasticsearch.yml -n hermes-messenger; \
	helm install messaging-redis bitnami/redis -n hermes-messenger; \
	kubectl apply -f services/messaging/deployment/service.yml -n hermes-messenger


deploy-service-authentication:
	helm install authentication-redis bitnami/redis -n hermes-messenger; \
	kubectl apply -f services/authentication/deployment/service.yml -n hermes-messenger


deploy-all: deploy-service-user \
	deploy-service-friends \
	deploy-service-messaging \
	deploy-service-authentication
