import 'reflect-metadata';
import WebSocket, { RawData, WebSocketServer } from 'ws';
import { config as actionConfig } from './config/actions.config';
import { WebSocketHandler } from './handlers/Handler';
import { User } from './user/User';
import { applicationContainer } from './inversify.config';
import { TYPES } from './types';
import { IConnectionRepository } from './connection/IConnectionRepository';
import { IUserRepository } from './user/IUserRepository';

const wss = new WebSocketServer({ port: 3000 });

const userRepo = applicationContainer.get<IUserRepository>(TYPES.UserRepository);

const Bob = new User({ username: 'Bob' });
userRepo.addUser(Bob);
console.log(userRepo.getUser(Bob.uuid));

const anotherUserRepo = applicationContainer.get<IUserRepository>(TYPES.UserRepository);
console.log(anotherUserRepo.getUser(Bob.uuid));

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
