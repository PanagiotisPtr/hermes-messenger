import Repository from '../../../common/repository/Repository';
import { Message } from '../entity/Message';

export interface IMessageRepository extends Repository {
  addMessage(Message: Message): string;
  getMessage(uuid: string): Message | null;
  removeMessage(uuid: string): void;
}
