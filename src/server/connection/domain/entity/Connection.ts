import { WebSocket } from 'ws';
import { Entity } from '../../../packages/common/entity/Entity';

interface ConnectionProps {
  websocket: WebSocket;
}

export class Connection extends Entity<ConnectionProps> {
  get websocket(): WebSocket {
    return this.props.websocket;
  }
}