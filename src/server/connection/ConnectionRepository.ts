import { inject, injectable } from 'inversify';
import { v4 as uuidV4 } from 'uuid';
import WebSocket from 'ws';
import Repository from '../common/repository/Repository';
import { TYPES } from '../types';
import { ConnectionStorage } from './ConnectionStorage';
import { IConnectionRepository } from './IConnectionRepository';

@injectable()
export class ConnectionRepository extends Repository implements IConnectionRepository {
  constructor(
    @inject(TYPES.ConnectionStorage) private connections: ConnectionStorage
  ) {
    super();
  }

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