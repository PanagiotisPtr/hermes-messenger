import { ConnectionRepository } from '../connection/ConnectionRepository';
import { WebSocketEvent, WebSocketHandler } from './Handler';

export const getHelloHandler = (
  repository: ConnectionRepository
): WebSocketHandler => {
  return async (event: WebSocketEvent) => {
    const ws = repository.getConnection(event.metadata.connectionUuid);
    if (ws === null) {
      return;
    }

    ws.send('hello ' + event.metadata.connectionUuid);
  };
};
