import Repository from '../../../common/repository/Repository';
import { User } from '../entity/User';

export interface IUserRepository extends Repository {
  addUser(user: User): string;
  getUser(uuid: string): User | null;
  removeUser(uuid: string): void;
  getUserFromConnectionUuid(connectionUuid: string): User | null;
  getUserFromUsername(username: string): User | null;
}
