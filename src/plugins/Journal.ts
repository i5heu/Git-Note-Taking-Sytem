import Config, { Plugin } from "../config";
import fs from 'fs';
import path from 'node:path';

export default class Journal {
    journalBase: string;

    constructor(settings: Plugin["settings"], config: Config) {
        const date = new Date();
        console.log("Plugin Journal: is running");
        if (typeof settings.journalPath === "string")
            this.journalBase = Config.basePath + "/" + settings.journalPath;

        this.createJournalFolderIfNotExists();
        this.createTemplateFolderIfNotExists();
        this.createTemplateFileIfNotExists();
        this.createYearFolderIfNotExists();
        this.createNewEntry();
        this.createNewEntry(true);
        this.createLinkForCurrentDay();

        console.log("Plugin Journal: is finished");
    }

    createJournalFolderIfNotExists() {
        if (!fs.existsSync(this.journalBase)) {
            fs.mkdirSync(this.journalBase);
        }
    }

    // create folder of current year if not exists
    createYearFolderIfNotExists() {
        const year = new Date().getFullYear();

        const yearFolder = this.journalBase + "/" + year;
        if (!fs.existsSync(yearFolder)) {
            fs.mkdirSync(yearFolder);
        }
    }

    createLinkForCurrentDay() {
        const date = new Date();
        const todayFilename = `${date.getFullYear()}-${date.getMonth()}-${date.getDate()}.md`;
        const todayFilepath = this.journalBase + "/" + date.getFullYear();
        const relativePath = path.relative(Config.basePath, todayFilepath);

        fs.writeFileSync(Config.basePath + "/today.md", `[${todayFilename}](${relativePath + "/" + todayFilename})`);
    }

    // create template folder if not exists
    createTemplateFolderIfNotExists() {
        const templateFolder = this.journalBase + "/template";
        if (!fs.existsSync(templateFolder)) {
            fs.mkdirSync(templateFolder);
        }
    }

    // create template file if not exists
    createTemplateFileIfNotExists() {
        let templateFile: string;

        // check if the current weekday file exists if not, create a empty file
        if (!fs.existsSync(this.journalBase + `/template/${this.weekday}.md`)) {
            templateFile = "";
            // write template file
            fs.writeFileSync(this.journalBase + `/template/${this.weekday}.md`, templateFile);
        }
    }

    // copy current weekday template file , with the current date as filename, into the current year folder
    createNewEntry(nextDay = false) {
        const date = new Date();

        let filename: string;
        if (nextDay) {
            filename = `${date.getFullYear()}-${date.getMonth()}-${date.getDate()}.md`;
        } else {
            const dateTomorrow = new Date();
            dateTomorrow.setDate(date.getDate() + 1);
            filename = `${dateTomorrow.getFullYear()}-${dateTomorrow.getMonth()}-${dateTomorrow.getDate()}.md`;
        }

        const filepath = this.journalBase + "/" + date.getFullYear() + "/" + filename;
        if (!fs.existsSync(filepath))
            fs.writeFileSync(filepath, this.templateFileContentOfCurrentWeekday);
    }

    // get template file content for the current weekday
    get templateFileContentOfCurrentWeekday() {
        const templateFile = fs.readFileSync(this.journalBase + `/template/${this.weekday}.md`, 'utf8');
        return templateFile;
    }

    get weekday() {
        const date = new Date();
        return date.toLocaleDateString("en-gb", { weekday: "long" });
    }
}