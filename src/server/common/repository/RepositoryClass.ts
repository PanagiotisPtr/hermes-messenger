import { WebSocketHandler } from '../../handlers/Handler';
import Repository from './Repository';

export interface RepositoryClass {
  new(...args: any): Repository;
}