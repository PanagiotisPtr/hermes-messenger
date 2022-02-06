import { HandlerClass } from '../common/handler/HandlerClass';
import { RepositoryClass } from '../common/repository/RepositoryClass';
import { ConnectionRepository } from '../connection/ConnectionRepository';
import HelloHandler from '../handlers/HelloHandler';

interface ActionsConfig {
  [action: string]: {
    handlerClass: new(...args: any) => HandlerClass,
    repositories: RepositoryClass[]
  }
}

export const config: ActionsConfig = {
  hello: {
    handlerClass: HelloHandler,
    repositories: [
      ConnectionRepository
    ]
  }
};