import express from "express";
import { renderSend } from "./helper";
import PreRequest from "./preRequest";
import Pug from "./pug";

/**
 * Config
 * TODO: Add config file and infrastructure.
 */
const port = 8080;

// Compile the Pug templates
const pug = new Pug();
// init express
const app = express();

let loginAttempt: string[number] | undefined[] = [];

app.get("/", (req, res) => {
    if (!PreRequest.userSpace(req, res)) return;

    renderSend(res, pug.home, {
        name: 'Timothy'
    });
});

app.get("/login", (req, res) => {
    if (!PreRequest.loginAttempts(req, res, loginAttempt)) return;

    renderSend(res, pug.login, {});
});


// start the Express server
app.listen(port, () => {
    console.log(`server started at http://localhost:${port}`);
});

// empty the loginAttempt array every 3 minutes
// TODO: maybe find a better place to put this
setInterval(() => {
    console.log("Clearing loginAttempts");
    loginAttempt = [];
}, 30000);