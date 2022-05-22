import * as http from "http";
import axios from "axios";
import { v4 as uuidv4 } from "uuid";

// HTTP SERVER HELLO WORLD
function main() {
  setInterval(caller, 2000);

  const server = http.createServer((req, res) => {
    res.statusCode = 200;
    res.setHeader("Content-Type", "text/plain");
    res.end("Hello World\n");
  });

  server.listen(80, "0.0.0.0", () => {
    console.log("Server running at http://localhost:80/");
  });
}
main();

function caller() {
  const data = {
    id: uuidv4(),
    name: "Test Service",
  };

  axios
    .post("http://base/register", data)
    .then((res) => {
      console.log(`Status: ${res.status}`);
    })
    .catch((err) => {
      console.error(err);
    });
}
