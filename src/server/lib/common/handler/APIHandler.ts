import { APIEvent } from '../event/APIEvent';

export type APIHandler<E extends APIEvent> = (event: E) => Promise<any>;
