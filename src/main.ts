import express from "express";
import { renderSend } from "./helper";
import InitializeGitRepo from "./gitManager";
import PreRequest from "./preRequest";
import Pug from "./pug";
import GitManager from "./gitManager";
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

    async function* getFiles(dir): AsyncIterableIterator<string> {
        const dirents = await readdir(dir, { withFileTypes: true });
        for (const dirent of dirents) {
            const res = resolve(dir, dirent.name);
            if (dirent.isDirectory()) {
                yield* getFiles(res);
            } else {
                yield res;
            }
        }
    }
    const files = await getFiles(git.options.baseDir);

    const filesClean = [];

    for await (const f of getFiles(git.options.baseDir)) {
        filesClean.push(f.replace(git.options.baseDir, ""));
    }

    renderSend(res, pug.home, {
        name: 'Timothy',
        files: filesClean
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