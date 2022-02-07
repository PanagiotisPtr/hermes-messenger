import 'reflect-metadata';
import WebSocket, { RawData, WebSocketServer } from 'ws';
import { config as actionConfig } from './config/actions.config';
import { WebSocketHandler } from './handlers/Handler';
import { applicationContainer } from './inversify.config';
import { TYPES } from './types';
import { IConnectionRepository } from './connection/IConnectionRepository';

const wss = new WebSocketServer({ port: 3000 });

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
for (const [action, handlerClass] of Object.entries(actionConfig)) {
  handlers.set(action, applicationContainer.resolve(handlerClass).handler);
}

wss.on('connection', function connection(ws: WebSocket) {
  const connectionRepo = applicationContainer.get<IConnectionRepository>(TYPES.ConnectionRepository);

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
