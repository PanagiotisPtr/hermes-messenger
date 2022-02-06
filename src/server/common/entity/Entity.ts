import { v4 as uuidV4 } from 'uuid';

export class Entity<Props> {
  public readonly uuid: string;
  protected props: Props;

  constructor(props: Props, uuid?: string) {
    this.uuid = uuid ? uuid : uuidV4();
    this.props = props;
  }
};
