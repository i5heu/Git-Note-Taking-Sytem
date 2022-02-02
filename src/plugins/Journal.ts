import Config, { Plugin } from "../helper/config";

export default class Journal {
    constructor(settings: Plugin["settings"], config: Config) {
        console.log("Journal constructor", settings);
    }
}