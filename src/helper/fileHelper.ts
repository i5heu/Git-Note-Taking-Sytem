import fs from "fs";
const fsPromises = fs.promises;

export default class FileHelper {

  /**
   * will return the names of all files and folders of given path
   */
  static async getFileListInFolder(path: string): Promise<string[]> {
    try {
      return fsPromises.readdir(path);
    } catch (err) {
      console.error("Error occured while reading directory!", err);
    }
  }

  // find a markdown table in a string with line breaks
  static async findMarkdownTable(text: string) {
    let lines = text.split("\n");

    let tables: string[] = [];

    lines.forEach((line, index) => {
      if (line[0] == "|") {
        if (tables[tables.length - 1].length != 0) tables[tables.length - 1] += "\n";
        tables[tables.length - 1] += line;
      } else {
        if (typeof tables[tables.length] != "string" && tables[tables.length - 1] != "") {
          tables.push("");
        }
      }
    });

    return tables;
  }

  static async getContentOfFile(path: string) {
    return fs.readFileSync(path, "utf8");
  }

  static removeExcessWhitespace(string) {
    while (string.indexOf("  ") > -1) {
      string = string.replace("  ", " ");
    }
    return string;
  }
}