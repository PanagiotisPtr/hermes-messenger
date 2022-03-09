import { Entity } from '../../../lib/common/entity/Entity';

type UserUuid = string;

interface IMessageProps {
  from: UserUuid;
  to: UserUuid;
  content: string;
}

export class Message extends Entity<IMessageProps> {
  get from(): UserUuid {
    return this.props.from;
  }

  get to(): UserUuid {
    return this.props.to;
  }

  get content(): string {
    return this.props.content;
  }
}