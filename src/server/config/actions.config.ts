import { HandlerClass } from '../common/handler/HandlerClass';
import HelloHandler from '../handlers/HelloHandler';
import MessageHandler from '../handlers/MessageHandler';

interface ActionsConfig {
  [action: string]: new(...args: any) => HandlerClass
}

export const config: ActionsConfig = {
  hello: HelloHandler,
  message: MessageHandler
};