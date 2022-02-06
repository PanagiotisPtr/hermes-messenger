import { WebSocketHandler } from '../../handlers/Handler';

export interface HandlerClass {
  new(...args: any): {
    get handler(): WebSocketHandler;
  };
}