import Config, { Plugin } from "../config";
import fs from "fs";
import FileHelper from "../helper/fileHelper";
const fsPromises = fs.promises;

export interface TodoItem {
  id: string | false;
  dependencies: string[] | undefined;
  repeat: string | undefined;
  priority: number;
  type: string;
  title: string;
  text: string;
  completed: boolean;
  createdAt: Date;
  beginsAt: Date | undefined;
  endsAt: Date | undefined;
  tags: string[] | undefined;
  durationInMinutes: number | undefined;
}

export default class ToDo {
  toDos: TodoItem[] = [];
  todoBase: string;

  constructor(settings: Plugin["settings"], config: Config) {
    this.todoBase = Config.basePath + "/" + settings.todoPath;
  }

  async run() {
    // console.log("Plugin ToDo: is running");
    // await this.createToDoFolderIfNotExists();
    // await this.fillTodoArray();
    // console.log("Plugin ToDo: is finished");
  }

  createToDoFolderIfNotExists() {
    if (!fs.existsSync(this.todoBase)) {
      fs.mkdirSync(this.todoBase);
    }
  }

  async fillTodoArray() {
    this.toDos = [];

    // read from file system
    FileHelper.getFileListInFolder(this.todoBase).then(async (files) => {
      files.forEach(async (file) => {
        const todoText = await FileHelper.getContentOfFile(
          this.todoBase + "/" + file
        );
        const todo = ParserTodo.parse(todoText);
        this.toDos.push(todo);
      });
    });
  }

  static async add(todo: TodoItem, todoBase: string): Promise<TodoItem["id"]> {
    const newItem = new AddToDo(todo, todoBase);
    return await newItem.run();
  }

  findById(id: string) {
    return this.toDos.find((todo) => todo.id === id);
  }

  setById(id: string, todo: TodoItem) {
    const index = this.toDos.findIndex((t) => t.id === id);
    this.toDos[index] = todo;
  }
}

class AddToDo {
  todo: TodoItem;
  todoBase: string;

  constructor(todo: TodoItem, todoBase: string) {
    this.todo = todo;
    this.todoBase = todoBase;
  }

  async run(): Promise<TodoItem["id"]> {
    const id = await this.findFreeId();

    let footer = " ";
    if (this.todo.dependencies && this.todo.dependencies.length > 0) {
      footer = "Dependencies: ";
      for (const dependency of this.todo.dependencies) {
        footer += `[${dependency}](./${dependency}.md) `;
      }
    }

    const content = `title: ${this.todo.title}
dependencies: ${this.todo.dependencies ? this.todo.dependencies.join(", ") : ""}
repeat: ${this.todo.repeat}
priority: ${this.todo.priority}
type: ${this.todo.type}
text: ${this.todo.text}
completed: ${this.todo.completed}
createdAt: ${this.todo.createdAt.toISOString()}
beginsAt: ${this.todo.beginsAt ? this.todo.beginsAt.toISOString() : ""}
endsAt: ${this.todo.endsAt ? this.todo.endsAt.toISOString() : ""}
tags: ${this.todo.tags ? this.todo.tags.join(", ") : ""}
durationInMinutes: ${this.todo.durationInMinutes}
--- body ---
${this.todo.text}
--- footer ---
${footer}
`;

    // create file
    await fsPromises.writeFile(this.todoBase + "/" + id + ".md", content);

    return id;
  }

  async findFreeId() {
    const filesInToDoFolder = FileHelper.getFileListInFolder(this.todoBase);
    const possibleChars = "abcdefghijklmnopqrstuvwxyz123456789".split("");
    let id = ["a"];

    while (fs.existsSync(this.todoBase + "/" + id.join("") + ".md")) {
      // if the next possible char dose not exist
      if (!possibleChars[possibleChars.indexOf(id[id.length - 1]) + 1]) {
        if (id[id.length - 2] == possibleChars[possibleChars.length - 1]) {
          // z to 11
          id[id.length - 1] = possibleChars[0];
          id.push(possibleChars[0]);
        } else {
          // 1z to 2a
          id[id.length - 1] = possibleChars[0];
          id[id.length - 2] =
            possibleChars[possibleChars.indexOf(id[id.length - 2]) + 1];
        }
      } else {
        id[id.length - 1] =
          possibleChars[possibleChars.indexOf(id[id.length - 1]) + 1];
      }
    }

    return id.join("");
  }
}

class ParserTodo {
  static parse(todoText: string): TodoItem {
    const sections = FileHelper.sectionContent(todoText);

    const todo: TodoItem = {
      id: this.getStringFromSection(sections.header, "id"),
      dependencies: this.getArrayFromSection(sections.header, "dependencies"),
      repeat: this.getStringFromSection(sections.header, "repeat"),
      priority: this.getIntFromSection(sections.header, "priority"),
      type: this.getStringFromSection(sections.header, "type"),
      title: this.getStringFromSection(sections.header, "title"),
      text: sections.body,
      completed: this.getBooleanFromSection(sections.header, "completed"),
      createdAt: this.getDateFromSection(sections.header, "createdAt"),
      beginsAt: this.getDateFromSection(sections.header, "beginsAt"),
      endsAt: this.getDateFromSection(sections.header, "endsAt"),
      tags: this.getArrayFromSection(sections.header, "tags"),
      durationInMinutes: this.getIntFromSection(
        sections.header,
        "durationInMinutes"
      ),
    };

    return todo;
  }

  static getStringFromSection(header: string, key: string): string | undefined {
    const keyLine = header.split("\n").find((line) => line.indexOf(key) > -1);
    if (!keyLine) return;

    const keyValue = keyLine.split(":")[1].trim();
    return keyValue;
  }
  static getIntFromSection(header: string, key: string): number | undefined {
    const keyLine = header.split("\n").find((line) => line.indexOf(key) > -1);
    if (!keyLine) return;

    const keyValue = keyLine.split(":")[1].trim();
    return parseInt(keyValue);
  }

  static getArrayFromSection(
    header: string,
    key: string
  ): string[] | undefined {
    const keyLine = header.split("\n").find((line) => line.indexOf(key) > -1);
    if (!keyLine) return;
    const keyValueArray = keyLine.split(":")[1].trim().split(",");
    return keyValueArray;
  }

  static getDateFromSection(header: string, key: string): Date | undefined {
    const keyLine = header.split("\n").find((line) => line.indexOf(key) > -1);
    if (!keyLine) return;

    const keyValue = keyLine.split(":")[1].trim();
    return new Date(keyValue);
  }

  static getBooleanFromSection(
    header: string,
    key: string
  ): boolean | undefined {
    const keyLine = header.split("\n").find((line) => line.indexOf(key) > -1);
    if (!keyLine) return;

    const keyValue = keyLine.split(":")[1].trim();
    return keyValue === "true";
  }
}
