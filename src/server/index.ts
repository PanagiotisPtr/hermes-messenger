import WebSocket, { RawData, WebSocketServer } from 'ws';
import { ConnectionRepository } from './connection/ConnectionRepository';
import { ConnectionStorage } from './connection/ConnectionStorage';
import { WebSocketHandler } from './handlers/Handler';
import { getHelloHandler } from './handlers/HelloHandler';

const wss = new WebSocketServer({ port: 3000 });

const connectionStorage: ConnectionStorage = new Map();

interface WebSocketAction {
  action: string;
  payload: string;
}

const isWebSocketAction = (something: any): something is WebSocketAction => {
  return typeof something === 'object' &&
    something.action &&
    something.payload &&
    typeof something.action === 'string' &&
    typeof something.payload === 'string';
};

const handlers = new Map<string, WebSocketHandler>();

handlers.set('hello', getHelloHandler(
  new ConnectionRepository(connectionStorage)
));

wss.on('connection', function connection(ws: WebSocket) {
  const connectionRepo = new ConnectionRepository(connectionStorage);

  const connectionUuid = connectionRepo.addConnection(ws);

  ws.on('message', (data: RawData) => {
    ws.send(connectionUuid);

    try {
      const messageObject = JSON.parse(data.toString());

      if (! isWebSocketAction(messageObject)) {
        throw Error('Not an action');
      }

      const handler = handlers.get(messageObject.action);
      if (handler !== undefined) {
        handler({
          payload: messageObject.payload,
          metadata: {
            connectionUuid
          }
        });
      }
    } catch (err) {
      ws.send('invalid message:' + err);
    }
  });

  ws.on('close', (_code: number, _reason: Buffer) => {
    connectionRepo.removeConnection(connectionUuid);
  });

  ws.send('something');
});
