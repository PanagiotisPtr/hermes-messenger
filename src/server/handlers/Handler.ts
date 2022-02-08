export interface WebSocketEvent {
  uuid: string;
  payload: string;
  metadata: {
    connectionUuid: string;
    [attribute: string]: string;
  };
}

export type WebSocketHandler = (event: WebSocketEvent) => Promise<any>;
