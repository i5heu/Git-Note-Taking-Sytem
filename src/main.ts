import express from "express";
import { renderSend } from "./helper/renderHelper";
import InitializeGitRepo from "./gitManager";
import PreRequest from "./preRequest";
import Pug from "./pug";
import GitManager from "./gitManager";
import { Tree } from "./helper/fileTreeHelper";
import Config from "./config";
import PluginManager from "./pluginManager";
import Lock from "./lock";
const { resolve } = require('path');
const { readdir } = require('fs').promises;

async function init() {
    const conf = new Config();
    const lock = new Lock();

    const pluginManager = new PluginManager(conf, lock);
    const git = new GitManager(conf, lock, pluginManager);
    await git.initialPullOrClone();
    conf.createDefaultConfigIfNotExists();

    pluginManager.runPluginsOverFiles();
    pluginManager.schedulePluginRuns();

    git.pullInterval();
    lock.emptyQueueFunction = git.commitAndPush.bind(git);
    server(conf, git);
}

function server(conf, git: GitManager) {
    // Compile the Pug templates
    const pug = new Pug();
    // init express
    const app = express();

    let loginAttempt: string[number] | undefined[] = [];

    // app.get("/", async (req, res) => {
    //     if (!PreRequest.userSpace(req, res)) return;


    //     const tree = new Tree(git.options.baseDir);
    //     const filesClean = await tree.getFileTree();;

    //     renderSend(res, pug.home, {
    //         name: 'Timothy',
    //         files: JSON.stringify(filesClean, null, 2)
    //     });
    // });

    // app.get("/login", (req, res) => {
    //     if (!PreRequest.loginAttempts(req, res, loginAttempt)) return;

    //     renderSend(res, pug.login, {});
    // });

    app.get("/health", (req, res) => {
        res.send('OK');
    });

    app.get("/pull", (req, res) => {
        git.pull();
        res.send('OK');
    });


    // start the Express server
    app.listen(conf.port, () => {
        console.log(`server started at http://localhost:${conf.port}`);
    });

    // empty the loginAttempt array every 3 minutes
    // TODO: maybe find a better place to put this
    setInterval(() => {
        console.log("Clearing loginAttempts");
        loginAttempt = [];
    }, 30000);
}

init();