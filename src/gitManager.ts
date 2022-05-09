import simpleGit, { PullResult, SimpleGit, SimpleGitOptions } from 'simple-git';
import os from 'os';
import fs from 'fs';
import Config from './config';
import Lock from './lock';
import PluginManager from './pluginManager';

export default class GitManager {
    git: SimpleGit;
    conf: Config;
    lock: Lock;
    pluginManager: PluginManager;

    options: Partial<SimpleGitOptions> = {
        baseDir: Config.basePath,
        maxConcurrentProcesses: 6,
    };

    constructor(conf: Config, lock: Lock, pluginManager: PluginManager) {
        this.conf = conf;
        this.lock = lock;
        this.pluginManager = pluginManager;
        this.git = simpleGit(this.options);
    }

    public async initialPullOrClone() {
        const isRepo = await this.git.checkIsRepo();

        if (isRepo) {
            console.log('Repo exists, pulling...');
            await this.git.pull();
            console.log('Pulled');
        } else {
            console.log('Git repo not found, cloning...');
            await this.git.clone(this.conf.repo, Config.basePath);
            console.log('Cloned');
        }
    }

    // pull from remote
    public async pull() {
        await this.lock.waitForFreeLockAndLock('pull', 30, async () => {
            console.log('pulling...');
            const result = await this.git.pull();
            console.log('run plugins...');
            this.runPluginsIfChangesFromPull(result);
        });
    }

    // commit changes
    public async commit(message: string) {
        await this.lock.waitForFreeLockAndLock('commit', 30, async () => {
            await this.git.add('./*');
            await this.git.commit(message);
        });
    }

    // push to remote
    public async push() {
        await this.lock.waitForFreeLockAndLock('push', 30, async () => {
            await this.git.push();
        });
    }

    // commit and push
    public async commitAndPush(message: string, noUnlockCommit: boolean = false) {
        await this.lock.waitForFreeLockAndLock('commit and push', 30, async () => {
            const status = await this.git.status();
            if (status.ahead == 0 && status.files.length == 0) return;
            console.log('committing and pushing...');
            
            await this.git.add('./*');
            await this.git.commit(message);
            await this.git.push();
            console.log('committing and pushing... done');
        }, noUnlockCommit);
    }

    // run a pull every x seconds
    public async pullInterval() {
        setInterval(async () => {
            await this.lock.waitForFreeLockAndLock('pull interval', 30, async () => {
                const result = await this.git.pull();
                this.runPluginsIfChangesFromPull(result);
            });
        }, this.conf.pullInterval * 1000);
    }

    private runPluginsIfChangesFromPull(result: PullResult) {
        if (result.summary.changes) {
            this.pluginManager.setDirty();
        }
    }
}