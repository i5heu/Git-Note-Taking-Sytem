import Config from "./helper/config";
import { Tree } from "./helper/fileTreeHelper";

export default class PluginManager {
    config: Config;

    constructor(config) {
        this.config = config;
    }

    async runPluginsOverFiles() {
        const tree = new Tree(Config.basePath);
        await tree.prepare();

        // iterate over plugins and filter for fileExtensions trigger
        for (const plugin of this.config.plugins) {
            if (!plugin.runOnAllWithType) continue;

            console.log("RUNNING PLUGIN", plugin.name);

            // use the iterateOverTree to visit all folders
            await tree.iterateOverTree(async chunk => {
                async function run(module) {
                    const a = await import(module);

                    // go over induvidual files
                    for await (const file of chunk.subItems) {
                        // filter files
                        if (plugin.runOnAllWithType.includes(file.fileExtension))
                            new a.default(file, chunk);
                    }
                }

                await run("./plugins/" + plugin.name);
            });

            console.log("FINISHED RUNNING PLUGIN", plugin.name);
        }
    }

    schedulePluginRuns() {

    }
}