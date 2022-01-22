import os from 'os';

export default class Config {

    static get basePath() {
        return os.homedir() + '/.git-note-taking-system';
    }

    
}