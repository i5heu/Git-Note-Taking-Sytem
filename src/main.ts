import express from "express";
import { renderSend } from "./helper/renderHelper";
import InitializeGitRepo from "./gitManager";
import PreRequest from "./preRequest";
import Pug from "./pug";
import GitManager from "./gitManager";
import { Tree } from "./helper/fileTreeHelper";
const { resolve } = require('path');
const { readdir } = require('fs').promises;

/**
 * Config
 * TODO: Add config file and infrastructure.
 */
const port = 8080;
const repoSSH = "git@github.com:i5heu/Git-Note-Taking-Test.git"

const git = new GitManager(repoSSH);
git.initialPullOrClone();

// Compile the Pug templates
const pug = new Pug();
// init express
const app = express();

let loginAttempt: string[number] | undefined[] = [];

app.get("/", async (req, res) => {
    if (!PreRequest.userSpace(req, res)) return;


    const tree = new Tree(git.options.baseDir);
    const filesClean = await tree.getFileTree();;

    renderSend(res, pug.home, {
        name: 'Timothy',
        files: JSON.stringify(filesClean, null, 2)
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