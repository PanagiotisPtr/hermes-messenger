export interface WebSocketEvent {
  payload: string;
  metadata: {
    connectionUuid: string;
    [attribute: string]: string;
  };
}

export type WebSocketHandler = (event: any) => Promise<any>;