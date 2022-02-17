import 'reflect-metadata';
import WebSocket, { RawData, WebSocketServer } from 'ws';
import { applicationContainer } from './inversify.config';
import { v4 as uuidV4 } from 'uuid';
import { getAPIConfig } from './application/config/api.config';
import { IConnectionRepository } from './connection/domain/repository/IConnectionRepository';
import { TYPES as ConnectionTYPES} from './connection/types';
import { TYPES as UserTYPES } from './user/types';
import { Connection } from './connection/domain/entity/Connection';
import { IUserRepository } from './user/domain/repository/IUserRepository';
import { User } from './user/domain/entity/User';

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

const APIConfig = getAPIConfig(applicationContainer);

let userCount = 1;
wss.on('connection', function connection(ws: WebSocket) {
  const connectionRepo = applicationContainer.get<IConnectionRepository>(ConnectionTYPES.IConnectionRepository);
  const connectionUuid = connectionRepo.addConnection(new Connection({ websocket: ws }));

  // @TODO should have registered and logged in / gotten token before connecting to the websocket server
  const userRepo = applicationContainer.get<IUserRepository>(UserTYPES.IUserRepository);
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

      const handler = APIConfig.websocket[messageObject.action];
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
