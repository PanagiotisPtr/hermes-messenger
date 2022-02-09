import { inject, injectable } from 'inversify';
import { IConnectionRepository } from '../connection/IConnectionRepository';
import { IMessageRepository } from '../message/IMessageRepository';
import { Message } from '../message/Message';
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
      const fromConnectionUuid = event.metadata.connectionUuid;
      const fromConnection = this.connectionRepository.getConnection(fromConnectionUuid);
      console.log(fromConnectionUuid);
      const fromUser = this.userRepository.getUserFromConnectionUuid(fromConnectionUuid);

      if (!fromUser) {
        fromConnection?.send('no from user');
      }

      const payload = JSON.parse(event.payload);
      if (!payload.to || !payload.message) {
        fromConnection?.send('invalid payload');
      }

      const toUser = this.userRepository.getUserFromUsername(payload.to);
      if (toUser === null) {
        fromConnection?.send('invalid receiver');
      } else {
        const message = new Message({ from: fromUser!, to: toUser, content: payload.message });
        this.messageRepository.addMessage(message);
        const toConnectionUuids = toUser.activeConnections;
        for (const toConnectionUuid of toConnectionUuids) {
          const toConnection = this.connectionRepository.getConnection(toConnectionUuid);
          toConnection?.send(payload.message);
        }
      }
    };
  }
}

export default MessageHandler;
