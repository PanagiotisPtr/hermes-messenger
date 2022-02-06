import Repository from '../common/repository/Repository';
import { User } from './User';
import { v4 as uuidV4 } from 'uuid';
import { IUserRepository } from './IUserRepository';
import { MemoryDB } from '../common/infra/MemoryDB';

export class UserMemoryRepository extends Repository implements IUserRepository {
  private users: Map<string, User>;

  constructor(
    memoryDB: MemoryDB
  ) {
    super();
    this.users = memoryDB.getOrCreateTableForEntity(User);
  }

  addUser(user: User): string {
    const userUuid = uuidV4();

    this.users.set(userUuid, user);

    return userUuid;
  }

  getUser(uuid: string): User | null {
    const existingUser = this.users.get(uuid);
    if (existingUser) {
      return existingUser;
    }

    return null;
  }

  removeUser(uuid: string): void {
    const existingUser = this.users.get(uuid);
    if (existingUser) {
      this.users.delete(uuid);
    }
  }
}