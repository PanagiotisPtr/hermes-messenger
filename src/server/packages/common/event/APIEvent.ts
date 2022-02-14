export interface APIEvent {
  uuid: string;
  metadata: {
    connectionUuid: string;
    [attribute: string]: string;
  };
}
