import { RedisClientType, RedisFunctions, RedisModules, RedisScripts } from "redis";
import { Message } from "../../grpc-clients/messaging";
import WebSocketMessenger from "./WebSocketMessenger";

export async function startRedisListener<T>(
    messenger: WebSocketMessenger,
    client: RedisClientType<T & RedisModules, RedisFunctions, RedisScripts>,
) {
    client.subscribe('messages', rawMessage => {
        const message = JSON.parse(rawMessage) as Message
        const wsMessage = { type: "message", data: rawMessage }
        messenger.sendMessageToUser(message.From, JSON.stringify(wsMessage))
        messenger.sendMessageToUser(message.To, JSON.stringify(wsMessage))
    })
}