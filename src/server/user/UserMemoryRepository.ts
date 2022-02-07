import Repository from '../common/repository/Repository';
import { User } from './User';
import { IUserRepository } from './IUserRepository';
import { MemoryDB } from '../common/infra/MemoryDB';
import { inject, injectable } from 'inversify';
import { TYPES } from '../types';

@injectable()
export class UserMemoryRepository extends Repository implements IUserRepository {
  private users: Map<string, User>;

  constructor(
    @inject(TYPES.MemoryDB) memoryDB: MemoryDB
  ) {
    super();
    this.users = memoryDB.getOrCreateTableForEntity(User);
  }

  addUser(user: User): string {
    this.users.set(user.uuid, user);

    return user.uuid;
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