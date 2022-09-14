// Create WebSocket connection.
const socket = new WebSocket("ws://localhost:8080/ws");
socket.addEventListener("open", (event) => {
    console.log("Sending message to server");
    document.getElementById("title").style.color = "#00ff00";
    document.getElementById("loginSubmit").disabled = false;
});
socket.addEventListener("close", (event) => {
    console.log("Connection closed", event);
    document.getElementById("title").style.color = "#ff0000";
});

let t0;
let t1;
let UUID = "";

async function login() {
    var username = document.getElementById("username").value;
    var apiKey = document.getElementById("apikey").value;
    await socket.send(
        JSON.stringify({
            type: "REGISTER",
            data: {
                name: username,
                apiKey: apiKey,
            },
        })
    );

    setInterval(() => {
        t0 = performance.now();
        socket.send(
            JSON.stringify({
                type: "PING",
            })
        );
    }, 1000);

    document.getElementById("login").remove();
}

socket.addEventListener("message", (event) => {
    if (event.data == "PONG") {
        t1 = performance.now();
        document.getElementById("ping").innerHTML = `${(t1 - t0).toFixed(0)}ms`;
        return;
    }

    const data = JSON.parse(event.data);

    switch (data.type) {
        case "REGISTER":
            UUID = data.data;
            console.log("Registered", data);
            break;
        default:
            console.log("Unknown message", data);
    }
});