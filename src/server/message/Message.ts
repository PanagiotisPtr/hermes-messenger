import { Entity } from '../common/entity/Entity';
import { User } from '../user/User';

interface IMessageProps {
  from: User;
  to: User;
  content: string;
}

export class Message extends Entity<IMessageProps> {
  get from(): User {
    return this.props.from;
  }

  get to(): User {
    return this.props.to;
  }

  get content(): string {
    return this.props.content;
  }
}