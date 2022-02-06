import { Entity } from '../common/entity/Entity';

interface UserAttributes {
  username: string;
}

export class User extends Entity<UserAttributes> {
  get username(): string {
    return this.username;
  }
}