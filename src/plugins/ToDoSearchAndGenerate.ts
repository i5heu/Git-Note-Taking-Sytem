import Config, { Plugin } from "../config";
import fs from 'fs';
import ToDo from "./Todo";

export default class ToDoSearchAndGenerate {
    todoBase: string;
    filePath: string;


    constructor(settings: Plugin["settings"], config: Config) {
        this.todoBase = Config.basePath + "/" + settings.todoPath;

        if (typeof settings.path !== "string")
            throw new Error("No path for file defined in settings!");
        this.filePath = settings.path;
    }

    async run() {
        if (!fs.existsSync(this.todoBase)) {
            fs.mkdirSync(this.todoBase);
        }

        const content = await this.getContentOfFile(this.filePath);
        
        //check if file contains todo creation tag
        const todoCreationTag = "$TODO";
        if (content.indexOf(todoCreationTag) === -1) return;

        const contentLineByLine = content.split("\n");
        
        //iterate over all lines and find todo creation tag
        contentLineByLine.forEach(async (line, index) => {
            if (line.indexOf(todoCreationTag) == -1) return;
            
        });
    }

    async createNewTodo(){
        
        ToDo.add({
            id: false,
            dependencies: undefined,
            repeat: undefined,
            priority: 0,
            type: "",
            title: "",
            text: "",
            completed: false,
            createdAt: new Date(),
            beginsAt: undefined,
            endsAt: undefined,
            tags: undefined,
            durationInMinutes: undefined
        })
    }

    async findFreeId() {

    }

    async getContentOfFile(path: string) {
        return fs.readFileSync(path, "utf8");
    }
}
