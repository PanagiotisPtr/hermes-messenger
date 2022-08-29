import WebSocket from "ws"

export default class WebSocketMessenger {
    private webSocketConnections: {
        [userUuid: string]: {
            [connectionUuid: string]: WebSocket.WebSocket
        }
    } = {}


    public addConnectionForUser(userUuid: string, connUuid: string, conn: WebSocket.WebSocket) {
        if (!this.webSocketConnections[userUuid]) {
            this.webSocketConnections[userUuid] = {}
        }
        this.webSocketConnections[userUuid][connUuid] = conn
    }

    public removeConnectionForUser(userUuid: string, connUuid: string) {
        delete this.webSocketConnections[userUuid][connUuid]
    }

    // SendMessageToUser sends a message to a user using websockets
    // sends to all open websocket connections for that user
    public async sendMessageToUser(userUuid: string, message: string): Promise<void> {
        if (!this.webSocketConnections[userUuid]) {
            return
        }

        for (const connectionUuid of Object.keys(this.webSocketConnections[userUuid])) {
            this.webSocketConnections[userUuid][connectionUuid].send(message)
        }
    }
}
