import simpleGit, { SimpleGit, SimpleGitOptions } from 'simple-git';
import os from 'os';
import fs from 'fs';

export default class GitManager {
    git: SimpleGit;
    repoSSH: string;

    options: Partial<SimpleGitOptions> = {
        baseDir: os.homedir() + '/.git-note-taking-system',
        maxConcurrentProcesses: 6,
    };

    constructor(repoSSH) {
        //create the git repo folder if it doesn't exist
        if (!fs.existsSync(this.options.baseDir)) {
            fs.mkdirSync(this.options.baseDir);
        }

        this.git = simpleGit(this.options);
        this.repoSSH = repoSSH;
    }

    public async initialPullOrClone() {
        const isRepo =  await this.git.checkIsRepo();

        if (isRepo) {
            console.log('Repo exists, pulling...');
            await this.git.pull();
            console.log('Pulled');
        } else {
            console.log('Git repo not found, cloning...');
            await this.git.clone(this.repoSSH, this.options.baseDir);
            console.log('Cloned');
        }
    }
}