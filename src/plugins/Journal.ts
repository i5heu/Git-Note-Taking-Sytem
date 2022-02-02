import { Plugin } from "../helper/config";

export default class Journal {
    constructor(settings: Plugin["settings"]) {
        console.log("Journal constructor", settings);
    }
}