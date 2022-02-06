import { ConnectionRepository } from '../connection/ConnectionRepository';
import { WebSocketEvent, WebSocketHandler } from './Handler';

export default class HelloHandler {
  constructor(
    private connectionRepository: ConnectionRepository
  ) {}

  get handler(): WebSocketHandler {
    return async (event: WebSocketEvent) => {
      const ws = this.connectionRepository.getConnection(event.metadata.connectionUuid);
      if (ws === null) {
        return;
      }

      ws.send('hello ' + event.metadata.connectionUuid);
    };
  }
}