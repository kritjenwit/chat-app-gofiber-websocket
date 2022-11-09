class Ws {
  dataSend = {
    eventName: "connected",
    data: {},
  };

  dataRecv = {};

  async connect() {
    this.ws = new WebSocket(`ws://${window.location.host}/ws`);

    return new Promise((resolve, reject) => {
      this.ws.onopen = (event) => {
        resolve(this);
      };
    });
  }

  set(eventName, data) {
    this.dataSend.eventName = eventName;
    this.dataSend.data = data;
    return this;
  }

  async send() {
    this.ws.send(JSON.stringify(this.dataSend));
    return new Promise((resolve, reject) => {
      this.ws.onmessage = (event) => {
        resolve(event);
      };
    });
  }

  onMessage() {
    this.ws.onmessage = (event) => {
      console.log(event)
    }
  }

  onCallback() {}
}