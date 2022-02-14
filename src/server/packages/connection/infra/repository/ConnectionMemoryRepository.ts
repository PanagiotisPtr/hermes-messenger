import { inject, injectable } from 'inversify';
import { v4 as uuidV4 } from 'uuid';
import WebSocket from 'ws';
import Repository from '../../../common/repository/Repository';
import { Connection } from '../../domain/entity/Connection';
import { IConnectionRepository } from '../../domain/repository/IConnectionRepository';
import { TYPES } from '../../types';

@injectable()
export class ConnectionMemoryRepository extends Repository implements IConnectionRepository {
  constructor(
    @inject(TYPES.ConnectionStorage) private connections: Map<string, Connection>
  ) {
    super();
  }

  addConnection(connection: Connection): string {
    const connectionUuid = uuidV4();

    this.connections.set(connectionUuid, connection);

    return connectionUuid;
  }

  getConnection(uuid: string): Connection | null {
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