import { WebSocketHandler } from '../../handlers/Handler';

export interface HandlerClass {
  get handler(): WebSocketHandler;
}