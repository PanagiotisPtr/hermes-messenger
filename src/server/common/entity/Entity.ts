import { v4 as uuidV4 } from 'uuid';

export class Entity<Attributes> {
  public readonly uuid: string;
  protected attributes: Attributes;

  constructor(attributes: Attributes, uuid?: string) {
    this.uuid = uuid ? uuid : uuidV4();
    this.attributes = attributes;
  }
};
