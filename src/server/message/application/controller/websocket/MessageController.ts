import { inject, injectable } from 'inversify';
import { IConnectionRepository } from '../../../../connection/domain/repository/IConnectionRepository';
import { TYPES as ConnectionTYPES } from '../../../../connection/types';
import { TYPES as UserTYPES } from '../../../../user/types';
import { IUserRepository } from '../../../../user/domain/repository/IUserRepository';
import { IMessageRepository } from '../../../domain/repository/IMessageRepository';
import { TYPES } from '../../../types';
import { Message } from '../../../domain/entity/Message';
import { WebSocketHandler } from '../../../../lib/common/handler/websocket/WebSocketHandler';
import { WebSocketEvent } from '../../../../lib/common/event/websocket/WebSocketEvent';
import { Controller } from '../../../../lib/common/controller/Controller';

@injectable()
export class MessageController extends Controller {
  constructor(
    @inject(ConnectionTYPES.IConnectionRepository)
    private connectionRepository: IConnectionRepository,
    @inject(UserTYPES.IUserRepository)
    private userRepository: IUserRepository,
    @inject(TYPES.IMessageRepository)
    private messageRepository: IMessageRepository
  ) {
    super();
  }

  get messageHandler(): WebSocketHandler {
    return async (event: WebSocketEvent) => {
      const fromConnectionUuid = event.metadata.connectionUuid;
      const fromConnection = this.connectionRepository.getConnection(fromConnectionUuid);
      const fromUser = this.userRepository.getUserFromConnectionUuid(fromConnectionUuid);

      if (!fromUser) {
        fromConnection?.websocket.send('no from user');
      }

      const payload = JSON.parse(event.payload);
      if (!payload.to || !payload.message) {
        fromConnection?.websocket.send('invalid payload');
      }

      const toUser = this.userRepository.getUserFromUsername(payload.to);
      if (toUser === null) {
        fromConnection?.websocket.send('invalid receiver');
      } else {
        const message = new Message({ from: fromUser?.uuid!, to: toUser?.uuid, content: payload.message });
        this.messageRepository.addMessage(message);
        const toConnectionUuids = toUser.activeConnections;
        for (const toConnectionUuid of toConnectionUuids) {
          const toConnection = this.connectionRepository.getConnection(toConnectionUuid);
          toConnection?.websocket.send(payload.message);
        }
      }
    };
  }
}
