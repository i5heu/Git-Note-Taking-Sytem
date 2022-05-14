import fs from "fs";
import puppeteer from "puppeteer";
import Config from "../config";

/**
 * npm i puppeteer puppeteer-extra-plugin-adblocker
 */

export interface SaveLinkAsPdfArchiveSettings {
  archivePath: string;
  enableJs: String[] | undefined;
  ignore: String[] | undefined;
}

export default class SaveLinkAsPdfArchive {
  settings: SaveLinkAsPdfArchiveSettings;
  config: Config;
  file: any;

  constructor(settings, config, file, chunk) {
    this.settings = settings;
    this.config = config;
    this.file = file;
  }

  async run() {
    this.createFolderIfNotExists(
      Config.basePath + "/" + this.settings.archivePath
    );
    const fileString = fs.readFileSync(this.file.path, "utf8");
    const links = this.findLinks(fileString);
    if (links)
      for (const link of links) {
        const domainRegex =
          /^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/?\n]+)/;
        const domainOfLink = link.match(domainRegex)[1];
        if (
          this.settings.ignore &&
          this.settings.ignore.indexOf(domainOfLink) !== -1
        )
          continue;

        const enableJs =
          this.settings.enableJs &&
          this.settings.enableJs.indexOf(domainOfLink) !== -1;

        try {
          const filePath = await this.saveLinkAsPdf(link, enableJs);
          this.addFooterIfNotExists(this.file.path);
          this.addLinkToFooterIfNotExists(this.file.path, link, filePath);
        } catch (error) {
          console.log("ERROR", error);
        }
      }
  }

  // find links in file
  findLinks(fileText): string[] {
    const regex =
      /(http|https):\/\/([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,!@?^=%&:\/~+#-]*[\w!@?^=%&\/~+#-])/g;
    return fileText.match(regex);
  }

  // save link as pdf
  async saveLinkAsPdf(link, enableJs = false) {
    // check if file exists
    const filename = this.generateFilename(link);
    const filePath =
      Config.basePath +
      "/" +
      this.settings.archivePath +
      "/" +
      filename +
      ".pdf";
    if (fs.existsSync(filePath)) return filePath;
    console.log("SAVING PAGE ... ", link);

    const browser = await puppeteer.launch();
    const page = await browser.newPage();
    await page.setJavaScriptEnabled(enableJs);
    await page.setViewport({
      width: 1024,
      height: 1024 * 3,
    });
    await page.setUserAgent(
      "Mozilla/5.0(iPad; U; CPU iPhone OS 3_2 like Mac OS X; en-us) AppleWebKit/531.21.10 (KHTML, like Gecko) Version/4.0.4 Mobile/7B314 Safari/531.21.10"
    );
    await page.emulateMediaFeatures([
      { name: "prefers-color-scheme", value: "dark" },
    ]);
    await page.goto(link, { waitUntil: ["networkidle0", "domcontentloaded"] });
    await page.waitForNetworkIdle();
    const height = await page.evaluate(
      () => document.documentElement.offsetHeight
    );
    await page.waitForTimeout(5 * 1000);
    await page.pdf({
      path: filePath,
      printBackground: true,
      width: 1024,
      height: height + 100,
    });
    await browser.close();

    console.log("SAVED LINK AS PDF", link);

    return filePath;
  }

  // generate filename from link
  generateFilename(link: string) {
    return link
      .replace(/[/\\?.%*#&$@€!:|`'"<>]/g, "-")
      .replace("--", "-")
      .replace("--", "-");
  }

  // create folder if not exists
  createFolderIfNotExists(path) {
    if (!fs.existsSync(path)) {
      fs.mkdirSync(path, { recursive: true });
    }
  }

  // add footer to markdown file if not exists
  async addFooterIfNotExists(filePath) {
    const fileString = fs.readFileSync(filePath, "utf8");
    const footer = "--------- Footer ---------";
    if (fileString.indexOf(footer) === -1) {
      fs.appendFileSync(filePath, "\n\n\n" + footer + "\n");
    }
  }

  // add link to footer if not exists
  async addLinkToFooterIfNotExists(filePath, link, archivePath) {
    const fileString = fs.readFileSync(filePath, "utf8");
    const afterFooter = fileString.split("--------- Footer ---------")[1];

    if (afterFooter.indexOf(archivePath.replace(Config.basePath, "")) === -1) {
      fs.appendFileSync(
        filePath,
        "\n" +
          "[Archive to " +
          this.generateFilename(link) +
          "]" +
          "(" +
          archivePath.replace(Config.basePath, "") +
          ")"
      );
    }
  }
}
