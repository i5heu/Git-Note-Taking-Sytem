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
      console.error("Error occurred while reading directory!", err);
    }
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