import fs from "fs";
import FileHelper from "./fileHelper";
const fsPromises = fs.promises;

export class Tree {
  rootPath;
  itemTree;
  indexes = [];

  constructor(directoryPath) {
    this.rootPath = directoryPath;
  }

  async prepare() {
    this.itemTree = [];
    await this.iterateOverSubFolder(this.rootPath, this.itemTree);
  }

  async getFileTree() {
    this.itemTree = [];
    await this.iterateOverSubFolder(this.rootPath, this.itemTree);
    return this.itemTree;
  }

  async iterateOverSubFolder(subFolderPath, treeChunk) {
    const subFolderList = await FileHelper.getFileListInFolder(subFolderPath);

    for (const subItemName of subFolderList) {
      const subItemPath = subFolderPath + "/" + subItemName;
      const isFolder = fs.lstatSync(subItemPath).isDirectory();

      switch (subItemName) {
        case ".git":
        case "AA_INDEX.auto.md":
          continue;
          break;
      }

      const folder = {
        name: subItemName,
        fileExtension: isFolder ? undefined : subItemName.split(".").pop(),
        path: subItemPath,
        directory: isFolder,
        subItems: isFolder ? [] : null,
      };

      treeChunk.push(folder);

      if (isFolder) {
        await this.iterateOverSubFolder(subItemPath, folder.subItems);
      }
    }
  }

  async iterateOverTree(itemFunction: (chunk: any) => Promise<any>) {
    await this.subFolderIterator(
      {
        root: true,
        directory: true,
        subItems: this.itemTree,
      },
      itemFunction
    );
  }

  async subFolderIterator(chunk, itemFunction: (chunk: any) => Promise<any>) {
    if (chunk.directory) {
      for (const chunky of chunk.subItems) {
        await this.subFolderIterator(chunky, itemFunction);
      }

      await itemFunction(chunk);
    }
  }
}
