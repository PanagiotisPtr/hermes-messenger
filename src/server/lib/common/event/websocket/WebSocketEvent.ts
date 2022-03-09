import { APIEvent } from '../APIEvent';

export interface WebSocketEvent extends APIEvent {
  payload: string;
}