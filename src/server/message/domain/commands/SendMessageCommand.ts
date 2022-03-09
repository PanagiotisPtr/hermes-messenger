import { User } from '../../../user/domain/entity/User';
import { Message } from '../entity/Message';

export interface SendMessageRequest {
  fromUser: User;
  toUser: User;
  message: Message;
}

export interface SendMessageResponse {
  success: boolean;
  error: Error;
}

export class SendMessageCommand {
  constructor(
    protected request: SendMessageRequest
  ) {}
}
