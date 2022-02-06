import { User } from './User';

export interface IUserRepository {
  addUser(user: User): string;
  getUser(uuid: string): User | null;
  removeUser(uuid: string): void;
}
