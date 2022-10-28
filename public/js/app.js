window.addEventListener("DOMContentLoaded", (_) => {
  let websocket = new WebSocket(`ws://${window.location.host}/ws`);
  let dataSend = {
    eventName: "connected",
    data: {},
  };

  websocket.onopen = () => {
    websocket.send(JSON.stringify(dataSend));
  };

  websocket.onmessage = (e) => {
    console.log(e);
    let data = JSON.parse(e.data);
    if (data.eventName == "connected") {
      console.log(data.chatHistorys);
    }
  };

  let form = document.getElementById("input-form");
  form.addEventListener("submit", function (event) {
    event.preventDefault();
    let username = document.getElementById("input-username");
    let text = document.getElementById("input-text");
    dataSend.eventName = "chat";
    dataSend.data = {
      username: username.value,
      text: text.value,
    };
    websocket.send(JSON.stringify(dataSend));
    text.value = "";
  });
});
