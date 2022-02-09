import 'reflect-metadata';
import WebSocket, { RawData, WebSocketServer } from 'ws';
import { config as actionConfig } from './config/actions.config';
import { WebSocketHandler } from './handlers/Handler';
import { applicationContainer } from './inversify.config';
import { TYPES } from './types';
import { IConnectionRepository } from './connection/IConnectionRepository';
import { v4 as uuidV4 } from 'uuid';
import { IUserRepository } from './user/IUserRepository';
import { User } from './user/User';

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

let userCount = 1;
wss.on('connection', function connection(ws: WebSocket) {
  const connectionRepo = applicationContainer.get<IConnectionRepository>(TYPES.IConnectionRepository);
  const connectionUuid = connectionRepo.addConnection(ws);

  // @TODO should have registered and logged in / gotten token before connecting to the websocket server
  const userRepo = applicationContainer.get<IUserRepository>(TYPES.IUserRepository);
  const connectedUser = new User({
    username: `User ${userCount++}`,
    activeConnections: new Set()
  });
  connectedUser.addConnection(connectionUuid);
  userRepo.addUser(connectedUser);

  console.log('User connected:', {
    username: connectedUser.username,
    connections: connectedUser.activeConnections
  });

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
          uuid: uuidV4(),
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
    connectedUser.removeConnection(connectionUuid);
  });

  ws.send('something');
});
