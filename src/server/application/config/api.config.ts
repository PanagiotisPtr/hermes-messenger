import { Container } from 'inversify';
import { Controller } from '../../lib/common/controller/Controller';
import { HttpHandler } from '../../lib/common/handler/http/HttpHandler';
import { WebSocketHandler } from '../../lib/common/handler/websocket/WebSocketHandler';
import { MessageController } from '../../message/application/controller/websocket/MessageController';

type Constructor<T> = new (...args: any) => T;

type ControllerConstructor<T extends Controller = Controller> = Constructor<T>;

interface WebSocketConfig {
  [action: string]: WebSocketHandler
}

interface HttpConfig {
  [path: string]: {
    [method: string]: HttpHandler
  }
}

interface APIConfig {
  websocket: WebSocketConfig;
  http: HttpConfig;
}

export const getAPIConfig = (
  applicationContainer: Container
): APIConfig => {
  const init = <T>(controller: ControllerConstructor<T>) => applicationContainer.resolve(controller);

  return {
    websocket: {
      message: init(MessageController).messageHandler
    },
    http: {

    }
  }
}
