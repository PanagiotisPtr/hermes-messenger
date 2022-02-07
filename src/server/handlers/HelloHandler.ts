import { inject, injectable } from 'inversify';
import { IConnectionRepository } from '../connection/IConnectionRepository';
import { TYPES } from '../types';
import { WebSocketEvent, WebSocketHandler } from './Handler';

@injectable()
class HelloHandler {
  constructor(
    @inject(TYPES.ConnectionRepository)
    private connectionRepository: IConnectionRepository
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

export default HelloHandler;
