import fs from "fs";

export default class SaveLinkAsPdfArchive {
    constructor(file, chunk) {
        
        const fileFromHardDrive = fs.readFileSync(file.path , 'utf8');

        console.log("SaveLinkAsPdfArchive constructor", file, fileFromHardDrive.slice(0, 100));
    }
}