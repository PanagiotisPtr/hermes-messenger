import { Entity } from '../common/entity/Entity';

export class User extends Entity {
  constructor(
    public uuid: string,
    public username: string
  ) {
    super();
  }
}