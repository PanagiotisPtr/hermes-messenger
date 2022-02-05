import { v4 as uuidV4 } from 'uuid';
import WebSocket from 'ws';
import { ConnectionStorage } from './ConnectionStorage';

export class ConnectionRepository {
  constructor(
    private connections: ConnectionStorage
  ) {}

  addConnection(ws: WebSocket): string {
    const connectionUuid = uuidV4();

    this.connections.set(connectionUuid, ws);

    return connectionUuid;
  }

  getConnection(uuid: string): WebSocket | null {
    const exisitngConnection = this.connections.get(uuid);
    if (exisitngConnection) {
      return exisitngConnection;
    }

    return null;
  }

  removeConnection(uuid: string): void {
    const existingConnection = this.connections.get(uuid);
    if (existingConnection) {
      this.connections.delete(uuid);
    }
  }
}