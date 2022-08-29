import { createServer, IncomingMessage, ServerResponse } from "http"
import next from "next"
import { parse } from "url"
import { parse as cookieParse } from "cookie"
import { WebSocketServer } from "ws"
import { tokenIsValid } from "../auth-utils/token"
import { v4 as uuidv4 } from "uuid";
import WebSocketMessenger from "./websockets/WebSocketMessenger"
import { createClient } from "redis"
import { startRedisListener } from "./websockets/redisListener"

(async () => {
    const dev = process.env.NODE_ENV !== "production"
    const hostname = "localhost"
    const port = 3000
    const app = next({ dev: true, hostname, port })
    const handle = app.getRequestHandler()
    const webSocketMessenger = new WebSocketMessenger()

    const client = createClient();

    client.on('error', (err) => console.warn('Redis Client Error', err));

    await client.connect();

    app.prepare().then(() => {
        const server = createServer(async (req: IncomingMessage, res: ServerResponse) => {
            try {
                const parsedUrl = parse(req.url ?? "", true)

                await handle(req, res, parsedUrl)
            } catch (err) {
                console.error("Error occurred handling", req.url, err)
                res.statusCode = 500
                res.end("internal server error")
            }
        }).listen(port)

        const wss = new WebSocketServer({
            noServer: true,
            path: "/api/ws",
        })

        server.on("upgrade", (request, socket, head) => {
            wss.handleUpgrade(request, socket, head, (websocket) => {
                wss.emit("connection", websocket, request)
            })
        })

        wss.on(
            "connection",
            async function connection(websocketConnection, connectionRequest) {
                const cookies = cookieParse(connectionRequest.headers.cookie ?? '')
                const { accessToken } = cookies

                const [dat, err] = await tokenIsValid(accessToken)
                if (err != null) {
                    websocketConnection.close(401)
                    return
                }
                const connectionUuid = uuidv4()
                webSocketMessenger.addConnectionForUser(dat, connectionUuid, websocketConnection)

                websocketConnection.onclose = () => {
                    webSocketMessenger.removeConnectionForUser(dat, connectionUuid)
                }

                // In case we want to listen to messages from client
                // we can use websocketConnection.on("message", ...) here
            }
        );
    })

    startRedisListener(webSocketMessenger, client)

})()