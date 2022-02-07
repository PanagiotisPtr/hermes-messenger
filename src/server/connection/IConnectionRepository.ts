import { WebSocket } from 'ws';

export interface IConnectionRepository {
  addConnection(ws: WebSocket): string;
  getConnection(uuid: string): WebSocket | null;
  removeConnection(uuid: string): void;
}
