import { Container } from 'inversify';
import { MemoryDB } from './common/infra/MemoryDB';
import { ConnectionRepository } from './connection/ConnectionRepository';
import { ConnectionStorage } from './connection/ConnectionStorage';
import { IConnectionRepository } from './connection/IConnectionRepository';
import { IMessageRepository } from './message/IMessageRepository';
import { MessageMemoryRepository } from './message/MessageMemoryRepository';
import { TYPES } from './types';
import { IUserRepository } from './user/IUserRepository';
import { UserMemoryRepository } from './user/UserMemoryRepository';

const applicationContainer = new Container();
applicationContainer.bind<MemoryDB>(TYPES.MemoryDB).to(MemoryDB).inSingletonScope();
applicationContainer.bind<ConnectionStorage>(TYPES.ConnectionStorage).toConstantValue(new Map());
applicationContainer.bind<IConnectionRepository>(TYPES.ConnectionRepository).to(ConnectionRepository);
applicationContainer.bind<IUserRepository>(TYPES.UserRepository).to(UserMemoryRepository);
applicationContainer.bind<IMessageRepository>(TYPES.MessageRepository).to(MessageMemoryRepository);

export { applicationContainer };