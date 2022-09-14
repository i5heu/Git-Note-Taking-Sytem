import { v4 as uuidv4 } from "https://cdn.jsdelivr.net/npm/uuid@9/dist/esm-browser/index.min.js";

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


function send(data) {
    createMessageLog(data, true);

    return socket.send(
        JSON.stringify(data)
    );
}

async function login() {
    var username = document.getElementById("username").value;
    var apiKey = document.getElementById("apikey").value;
    await send(
        {
            thread_id: uuidv4(),
            type: "REGISTER",
            data: {
                name: username,
                apiKey: apiKey,
            },
        }
    );

    setInterval(() => {
        t0 = performance.now();
        socket.send(
            JSON.stringify({
                thread_id: uuidv4(),
                type: "PING",
            }
            ));
    }, 1000);

    document.getElementById("login").remove();
    document.getElementById("controls").style.display = "block";
}
document.getElementById("loginSubmit").addEventListener("click", login);


function getFile() {
    send({
        thread_id: uuidv4(),
        type: "READ",
        data: {
            path: document.getElementById("path").value,
        },
    })
}
document.getElementById("getFile").addEventListener("click", getFile);











socket.addEventListener("message", (event) => {
    const data = JSON.parse(event.data);

    if (data.type != "PONG") {
        createMessageLog(data);
    }

    switch (data.type) {
        case "PONG":
            t1 = performance.now();
            document.getElementById("ping").innerHTML = `${(t1 - t0).toFixed(0)}ms`;
            break;
        case "REGISTER":
            UUID = data.data;
            console.log("Registered", data);
            break;
        default:
            console.log("Unknown message", data);
    }
});

function createThreadLog(data) {
    const messages = document.getElementById("messages");
    const message = document.createElement("details");
    message.className = "message";
    const summary = document.createElement("summary");
    summary.innerHTML = `${data.type} - ${data.thread_id}`;
    message.appendChild(summary);
    message.id = data.thread_id;

    //insert as first child
    messages.insertBefore(message, messages.firstChild);
}

function createMessageLog(data, send = false) {
    let message = document.getElementById(data.thread_id);
    if (!message) {
        createThreadLog(data);
        message = document.getElementById(data.thread_id);
    }

    const pre = document.createElement("pre");
    pre.className = send ? "send" : "receive";
    pre.innerHTML = JSON.stringify(data, null, 4);
    message.appendChild(pre);
}