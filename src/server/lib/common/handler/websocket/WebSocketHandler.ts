import { APIEvent } from '../../event/APIEvent';
import { WebSocketEvent } from '../../event/websocket/WebSocketEvent';
import { APIHandler } from '../APIHandler';

export type WebSocketHandler = APIHandler<WebSocketEvent>;
