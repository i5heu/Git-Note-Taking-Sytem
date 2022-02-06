import Config from "./config";
import { Tree } from "./helper/fileTreeHelper";
import Lock from "./lock";
const CronJob = require('cron').CronJob;

export default class PluginManager {
    config: Config;
    lock: Lock;
    dirty: boolean = false;

    constructor(config, lock: Lock) {
        this.config = config;
        this.lock = lock;
        this.dirty = false;

        this.runPluginsOverIfDirty();
    }

    async runPluginsOverFiles() {
        const tree = new Tree(Config.basePath);
        await tree.prepare();

        // iterate over plugins and filter for fileExtensions trigger
        for (const plugin of this.config.plugins) {
            if (!plugin.runOnAllWithType) continue;

            await this.lock.waitForFreeLockAndLock(plugin.name, plugin.timeout, async () => {
                console.log("RUNNING PLUGIN", plugin.name);

                // use the iterateOverTree to visit all folders
                await tree.iterateOverTree(async chunk => {
                    const run = async module => {
                        const a = await import(module);

                        // go over induvidual files
                        for await (const file of chunk.subItems) {
                            // filter files
                            if (plugin.runOnAllWithType.includes(file.fileExtension) || plugin.runOnAllWithType.includes("*")) {
                                const instanceRun = new a.default(file, chunk, plugin.settings, this.config);
                                if (instanceRun.run)
                                    await instanceRun.run();
                            }
                        }
                    }

                    await run("./plugins/" + plugin.name);
                });

                console.log("FINISHED RUNNING PLUGIN", plugin.name);
            });
        }
    }

    schedulePluginRuns() {
        for (const plugin of this.config.plugins) {
            if (!plugin.cron) continue;

            const job = new CronJob(plugin.cron, async () => {
                const run = async module => {
                    const a = await import(module);
                    new a.default(plugin.settings, this.config);
                }

                await this.lock.waitForFreeLockAndLock(plugin.name, plugin.timeout, async () => {
                    await run("./plugins/" + plugin.name);
                });
            });

            job.start();
        }
    }

    // run all plugins over all files if dirty
    private async runPluginsOverIfDirty() {
        setInterval(async () => {
            if (this.dirty) {
                this.dirty = false;
                await this.runPluginsOverFiles();
            }
        }, 3000);
    }

    setDirty() {
        this.dirty = true;
    }
}