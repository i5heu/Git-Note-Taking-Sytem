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

  static async appendToFile(path: string, content: string) {
    try {
      await fsPromises.appendFile(path, content);
    } catch (err) {
      console.error("Error occurred while appending to file!", err);
    }
  }

  static async writeFile(path: string, content: string) {
    try {
      await fsPromises.writeFile(path, content, 'utf8');
    } catch (err) {
      console.error("Error occurred while appending to file!", err);
    }
  }

  // will return the header body and optionally the footer
  static sectionContent(fileText: string) {
    const lineByLine = fileText.split("\n");
    const bodyTagLineNumber = lineByLine.findIndex(line => line.indexOf("--- body ---") > -1);
    const footerTagLineNumber = lineByLine.findIndex(line => line.indexOf("--- footer ---") > -1);

    let header: string[] | undefined;
    let body: string[] | undefined;
    let footer: string[] | undefined;

    if (bodyTagLineNumber > -1) {
      header = lineByLine.slice(0, bodyTagLineNumber);
    }

    if (bodyTagLineNumber > -1) {
      if (footerTagLineNumber > -1) {
        body = lineByLine.slice(bodyTagLineNumber + 1, footerTagLineNumber);
      } else {
        body = lineByLine.slice(bodyTagLineNumber + 1, lineByLine.length - 1);
      }
    } else if (footerTagLineNumber > -1) {
      body = lineByLine.slice(0, lineByLine.length - 1);
    }

    if (footerTagLineNumber > -1) {
      footer = lineByLine.slice(footerTagLineNumber + 1, lineByLine.length - 1);
    }

    return {
      header: header ? header.join("\n") : undefined,
      body: body ? body.join("\n") : undefined,
      footer: footer ? footer.join("\n") : undefined
    }
  }
}