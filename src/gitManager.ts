import simpleGit, { SimpleGit, SimpleGitOptions } from 'simple-git';
import os from 'os';
import fs from 'fs';
import Config from './helper/config';

export default class GitManager {
    git: SimpleGit;
    conf: Config;

    options: Partial<SimpleGitOptions> = {
        baseDir: Config.basePath,
        maxConcurrentProcesses: 6,
    };

    constructor(conf: Config) {
        this.conf = conf;
        this.git = simpleGit(this.options);
    }

    public async initialPullOrClone() {
        const isRepo =  await this.git.checkIsRepo();

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
}