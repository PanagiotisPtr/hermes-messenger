import { inject, injectable } from 'inversify';
import { TYPES } from '../../../../types';
import { MemoryDB } from '../../../common/infra/MemoryDB';
import Repository from '../../../common/repository/Repository';
import { User } from '../../domain/entity/User';
import { IUserRepository } from '../../domain/repository/IUserRepository';

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

  getUserFromConnectionUuid(connectionUuid: string): User | null {
    for (const user of this.users.values()) {
      for (const userConnectionUuid of user.activeConnections) {
        console.log(user.username);
        console.log(user.activeConnections);
        if (userConnectionUuid === connectionUuid) {
          return user;
        }
      }
    }

    return null;
  }

  getUserFromUsername(username: string): User | null {
    for (const user of this.users.values()) {
      if (user.username === username) {
        return user;
      }
    }

    return null;
  }
}