import Repository from '../common/repository/Repository';
import { Message } from './Message';
import { IMessageRepository } from './IMessageRepository';
import { MemoryDB } from '../common/infra/MemoryDB';

export class MessageMemoryRepository extends Repository implements IMessageRepository {
  private messages: Map<string, Message>;

  constructor(
    memoryDB: MemoryDB
  ) {
    super();
    this.messages = memoryDB.getOrCreateTableForEntity(Message);
  }

  addMessage(message: Message): string {
    this.messages.set(message.uuid, message);

    return message.uuid;
  }

  getMessage(uuid: string): Message | null {
    const existingUser = this.messages.get(uuid);
    if (existingUser) {
      return existingUser;
    }

    return null;
  }

  removeMessage(uuid: string): void {
    const existingUser = this.messages.get(uuid);
    if (existingUser) {
      this.messages.delete(uuid);
    }
  }
}