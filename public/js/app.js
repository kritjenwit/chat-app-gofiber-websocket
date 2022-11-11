class Ws {
  webSocketURL = `ws://localhost:3000/ws/`;

  connect(method) {
    return new WebSocket(this.webSocketURL + method);
  }
}
