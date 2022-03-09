import { inject } from 'inversify';
import { IMessageRepository } from '../../domain/repository/IMessageRepository';
import { TYPES } from '../../types';

export class MessageService {
  constructor(
    @inject(TYPES.IMessageRepository)
    private messageRepository: IMessageRepository
  ) {}

  sendMessage(
    fromUserUuid: string,
    toUserUuid: string,
    message: string
  ): [boolean, Error|null] {


    return [true, null];
  }
}