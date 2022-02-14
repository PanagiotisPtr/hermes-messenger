import { Entity } from '../../../common/entity/Entity';

type ConnectionUuid = string;

interface UserProps {
  username: string;
  activeConnections: Set<string>;
}

export class User extends Entity<UserProps> {
  get username(): string {
    return this.props.username;
  }

  get activeConnections(): ConnectionUuid[] {
    return Array.from(this.props.activeConnections.values());
  }

  addConnection(connectionUuid: ConnectionUuid): void {
    this.props.activeConnections.add(connectionUuid);
  }

  removeConnection(connectionUuid: ConnectionUuid): void {
    this.props.activeConnections.delete(connectionUuid);
  }
}