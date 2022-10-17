#!/bin/bash

echo -e "service:
  port: ${LISTEN_PORT:-80}
  grpcReflection: ${GRPC_REFLECTION:-false}

mongo:
  uri: \"${MONGO_URI:-mongodb://localhost:27017}\"
  db: \"${MONGO_DB:-test}\"

redis:
  addr: \"${REDIS_ADDRESS:-localhost:6379}\"
  password: \"${REDIS_PASSWORD}\"
  db: ${REDIS_DB:-0}
"
