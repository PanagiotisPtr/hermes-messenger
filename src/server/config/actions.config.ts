import { HandlerClass } from '../common/handler/HandlerClass';
import HelloHandler from '../handlers/HelloHandler';

interface ActionsConfig {
  [action: string]: new(...args: any) => HandlerClass
}

export const config: ActionsConfig = {
  hello: HelloHandler
};