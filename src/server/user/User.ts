import { Entity } from '../common/entity/Entity';

interface UserProps {
  username: string;
  activeConnections: Set<string>;
}

export class User extends Entity<UserProps> {
  get username(): string {
    return this.props.username;
  }

  get activeConnections(): string[] {
    return Array.from(this.props.activeConnections.values());
  }

  addConnection(connectionUuid: string): void {
    this.props.activeConnections.add(connectionUuid);
  }

  removeConnection(connectionUuid: string): void {
    this.props.activeConnections.delete(connectionUuid);
  }
}