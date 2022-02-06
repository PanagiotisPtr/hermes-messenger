import { table } from 'console';
import { Entity } from '../entity/Entity';

type EntityClass<E extends Entity> = new(...args: any) => E;

/** MemoryDB
 *
 * Maps a the table name for an entity to a list of entities.
 *
 * e.g.
 *
 * MemoryDB<User.constructor.name, <someUsersUuid, SomeUserEntity> >
 */
export class MemoryDB {
  private data: Map<string, Map<string, Entity>>;

  constructor() {
    this.data = new Map<string, Map<string, Entity>>();
  }

  getTableForEntity<E extends Entity>(entity: EntityClass<E>): Map<string, E> | null {
    const table = this.data.get(entity.constructor.name);
    if (table) {
      return table as Map<string, E>;
    }

    return null;
  }

  createTableForEntity<E extends Entity>(entity: EntityClass<E>): void {
    if (this.data.has(entity.constructor.name)) {
      return;
    }

    this.data.set(entity.constructor.name, new Map<string, E>());
  }

  getOrCreateTableForEntity<E extends Entity>(entity: EntityClass<E>): Map<string, E> {
    if (this.getTableForEntity(entity) === null) {
      this.createTableForEntity(entity);
    }

    return this.getTableForEntity(entity)!;
  }
}