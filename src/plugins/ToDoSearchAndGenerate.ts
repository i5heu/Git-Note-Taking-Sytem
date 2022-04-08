import Config, { Plugin } from "../config";
import fs from 'fs';
import ToDo, { TodoItem } from "./Todo";

export default class ToDoSearchAndGenerate {
    todoBase: string;
    filePath: string;


    constructor(settings: Plugin["settings"], config: Config, file) {
        if (typeof settings.todoPath === "string" && typeof settings.autoToDosPath === "string")
            this.todoBase = Config.basePath + "/" + settings.todoPath + "/" + settings.autoToDosPath;
        this.todoBase = Config.basePath + "/" + "ToDo/auto";

        if (typeof file.path !== "string")
            throw new Error("No path for file defined in settings!");
        this.filePath = file.path;
    }

    async run() {
        if (!fs.existsSync(this.todoBase)) {
            fs.mkdirSync(this.todoBase, { recursive: true });
        }

        const content = await this.getContentOfFile(this.filePath);

        //check if file contains todo creation tag
        const todoCreationTag = "$TODO";
        if (content.indexOf(todoCreationTag) === -1) return;

        const contentLineByLine = content.split("\n");

        let parentLevels: { lvl: number, toDoId: string }[] = [];
        let currentLevel = 0;

        //iterate over all lines and find todo creation tag
        for await (const line of contentLineByLine) {
            if (line.indexOf("-") !== -1) {
                // count spaces before dash to determine level
                currentLevel = line.split("-")[0].length;

                //remove all higher levels or same level
                parentLevels = parentLevels.filter(t => t.lvl < currentLevel);
            } else {
                parentLevels = [];
                currentLevel = 0;
            }

            if (line.indexOf(todoCreationTag) == -1) continue;

            //reverse sort by lvl
            parentLevels = parentLevels.sort((a, b) => b.lvl - a.lvl);

            const dependencies = parentLevels[0] ? [parentLevels[0].toDoId] : [];
            const toDoId = await this.createNewTodo(line, dependencies);
            parentLevels.push({ lvl: currentLevel, toDoId: toDoId ? toDoId : "" });

            console.log("parnetLevels", parentLevels);
        };
    }

    async createNewTodo(line: string, dependencies: TodoItem["dependencies"]): Promise<TodoItem["id"]> {
        console.log("Creating new todo", line, dependencies);

        const priorityMatch = line.match(/p:(\d*)/)
        const priority = priorityMatch ? parseInt(priorityMatch[1]) : 0;

        const endMatch = line.match(/e:([0-9T:-]*)/);
        const end = endMatch ? new Date(endMatch[1]) : undefined;

        const durationMatch = line.match(/d:([0-9]*)/);
        const duration = durationMatch ? parseInt(durationMatch[1]) : undefined;

        const titleMatch = line.match(/\$TODO ([^$]*)/);
        let title = titleMatch ? titleMatch[1] : "";
        if (!title) title = line.replace("- $TODO", "");

        const id = await ToDo.add({
            id: false,
            dependencies: dependencies,
            repeat: undefined, //TODO check if this todo has a repeat tag
            priority: priority,
            type: "",
            title: title,
            text: "",
            completed: false,
            createdAt: new Date(),
            beginsAt: undefined,
            endsAt: end,
            tags: undefined,
            durationInMinutes: duration
        },
            this.todoBase
        );

        return id;
    }

    async getContentOfFile(path: string) {
        return fs.readFileSync(path, "utf8");
    }
}
