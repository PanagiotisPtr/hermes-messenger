import { Entity } from '../common/entity/Entity';

interface UserProps {
  username: string;
}

export class User extends Entity<UserProps> {
  get username(): string {
    return this.props.username;
  }
}