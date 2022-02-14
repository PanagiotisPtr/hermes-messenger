import { HttpEvent } from '../../event/http/HttpEvent';
import { APIHandler } from '../APIHandler';

export type HttpHandler = APIHandler<HttpEvent>;
