import { type SetStateAction } from "react";
import { Client as StompClient, StompConfig } from "@stomp/stompjs";

import { PUBLIC_WEBSOCKET_URL } from "#environment/derived";

type ISetConnFunc = (value: SetStateAction<WSClient | null>) => void;

export class WSConfig extends StompConfig {
  constructor() {
    super();

    // brokerURL: "ws://localhost:61614/ws",
    this.webSocketFactory = webSocketFactory;
    this.connectHeaders = {
      login: "guest", // TODO
      passcode: "guest",
    };
    this.reconnectDelay = 5000;
    this.heartbeatIncoming = 4000;
    this.heartbeatOutgoing = 4000;
    this.connectionTimeout = 5000;
  }
}

export class WSClient extends StompClient {
  /**
   * @TODO fix this monstrosity of assigning methods in constructor
   */
  constructor(config: WSConfig, setConnectionFunc: ISetConnFunc) {
    super(config);
    this.onConnect = () => {
      setConnectionFunc(this);
    };
    this.onDisconnect = () => {
      setConnectionFunc(null);
    };
    this.onWebSocketError = async (error) => {
      console.log(`onWebSocketError ${JSON.stringify(error)}`, "WS");
    };

    this.onStompError = (frame) => {
      // Will be invoked in case of error encountered at Broker
      // Bad login/passcode typically will cause an error
      // Complaint brokers will set `message` header with a brief message. Body may contain details.
      // Compliant brokers will terminate the connection after any error
      console.log("Broker reported error: " + frame.headers["message"]);
      console.log("Additional details: " + frame.body);
    };
  }
}

function webSocketFactory() {
  return new WebSocket(PUBLIC_WEBSOCKET_URL);
}
