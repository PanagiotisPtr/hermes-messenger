import { Connection } from '../entity/Connection';

export interface IConnectionRepository {
  addConnection(connection: Connection): string;
  getConnection(uuid: string): Connection | null;
  removeConnection(uuid: string): void;
}
