import { Container } from 'inversify';
import { IConnectionRepository } from './connection/domain/repository/IConnectionRepository';
import { ConnectionMemoryRepository, ConnectionStorage } from './connection/infra/repository/ConnectionMemoryRepository';
import { TYPES as ConnectionTYPES } from './connection/types';
import { IMessageRepository } from './message/domain/repository/IMessageRepository';
import { MemoryDB } from './lib/common/infra/MemoryDB';
import { TYPES as CommonTYPES } from './lib/common/types';
import { IUserRepository } from './user/domain/repository/IUserRepository';
import { TYPES as UserTYPES } from './user/types';
import { TYPES as MessageTYPES } from './message/types';
import { UserMemoryRepository } from './user/infra/repository/UserMemoryRepository';
import { MessageMemoryRepository } from './message/infra/repository/MessageMemoryRepository';

const applicationContainer = new Container();
applicationContainer.bind<MemoryDB>(CommonTYPES.MemoryDB).to(MemoryDB).inSingletonScope();
applicationContainer.bind<ConnectionStorage>(ConnectionTYPES.ConnectionStorage).toConstantValue(new Map());
applicationContainer.bind<IConnectionRepository>(ConnectionTYPES.IConnectionRepository).to(ConnectionMemoryRepository);
applicationContainer.bind<IUserRepository>(UserTYPES.IUserRepository).to(UserMemoryRepository);
applicationContainer.bind<IMessageRepository>(MessageTYPES.IMessageRepository).to(MessageMemoryRepository);

export { applicationContainer };
