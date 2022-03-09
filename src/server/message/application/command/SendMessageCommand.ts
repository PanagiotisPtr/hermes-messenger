import { v4 as uuidV4 } from 'uuid';
import { Command } from '../../../lib/common/command/Command';
import { Request } from '../../../lib/common/command/Request';
import { Response } from '../../../lib/common/command/Response';
import { Message } from '../../domain/entity/Message';

export interface SendMessageRequest extends Request {
  fromUserUuid: string;
  toUserUuid: string;
  message: Message;
}

export interface SendMessageResponse extends Response {
  success: boolean;
  error: Error|null;
}

export class SendMessageCommand extends Command<
  SendMessageRequest,
  SendMessageResponse
> {
  constructor(
    protected request: SendMessageRequest
  ) {
    super(request);
  }

  async execute(): Promise<SendMessageResponse> {
    return {
      uuid: uuidV4(),
      success: true,
      error: null
    };
  }
}
