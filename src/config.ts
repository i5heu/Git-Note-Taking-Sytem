import os from 'os';
import fs from 'fs';

export interface Plugin {
    name: string;
    cron: string;
    timeout: number; // seconds: timeout is to control if a task should still run
    runOnAllWithType: string[];
    settings: { [key: string]: unknown };
}

export default class Config {

    conf: any;

    constructor() {
        //create the git repo folder if it doesn't exist
        if (!fs.existsSync(Config.basePath)) {
            fs.mkdirSync(Config.basePath);
        }

        const defaultConfFile = fs.readFileSync('defaultConfigs/config.json', 'utf8');
        const defaultConf = JSON.parse(defaultConfFile);

        let conf;

        if (fs.existsSync(Config.basePath + '/.Tyche.config.json')) {
            const confFile = fs.readFileSync(Config.basePath + '/.Tyche.config.json', 'utf8');
            conf = { ...defaultConf, ...JSON.parse(confFile) };
        } else {
            conf = defaultConf;
        }

        this.conf = conf;

        setInterval(() => {
            console.log('Reload config');
            let conf;

            if (fs.existsSync(Config.basePath + '/.Tyche.config.json')) {
                const confFile = fs.readFileSync(Config.basePath + '/.Tyche.config.json', 'utf8');
                conf = { ...defaultConf, ...JSON.parse(confFile) };
            } else {
                conf = defaultConf;
            }

            this.conf = conf;
        }, 15000);
    }

    createDefaultConfigIfNotExists() {
        if (fs.existsSync(Config.basePath + '/.Tyche.config.json')) return;
        const defaultConfFile = fs.readFileSync('defaultConfigs/config.json', 'utf8');
        fs.writeFileSync(Config.basePath + '/.Tyche.config.json', defaultConfFile);
    }

    static get basePath() {
        return os.homedir() + '/.Tyche';
    }

    get repoPath(): string {
        return this.conf.repo;
    }

    get port(): number {
        return this.conf.port;
    }

    get loginAttemptResetSeconds(): number {
        return this.conf.loginAttemptResetSeconds;
    }

    get plugins(): Plugin[] {
        return this.conf.plugins;
    }

    get repo(): string {
        return this.conf.repo;
    }

    get pullInterval(): number {
        return this.conf.pullInterval;
    }
}