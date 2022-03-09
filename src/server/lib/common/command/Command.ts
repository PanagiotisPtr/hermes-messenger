import { Request } from './Request';
import { Response } from './Response';

export abstract class Command<
  Req extends Request,
  Res extends Response
> {
  constructor(
    protected request: Request
  ) {}

  abstract execute(): Promise<Response>;
}
