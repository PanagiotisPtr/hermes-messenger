import { inject, injectable } from 'inversify';
import { IConnectionRepository } from '../connection/IConnectionRepository';
import { IMessageRepository } from '../message/IMessageRepository';
import { TYPES } from '../types';
import { IUserRepository } from '../user/IUserRepository';
import { WebSocketEvent, WebSocketHandler } from './Handler';

@injectable()
class MessageHandler {
  constructor(
    @inject(TYPES.IConnectionRepository)
    private connectionRepository: IConnectionRepository,
    @inject(TYPES.IUserRepository)
    private userRepository: IUserRepository,
    @inject(TYPES.IMessageRepository)
    private messageRepository: IMessageRepository
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

export default MessageHandler;
