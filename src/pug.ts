import pug from "pug";

export default class Pug {
  public home: pug.compileTemplate;
  public login: pug.compileTemplate;
  public compileTemplate: pug.compileTemplate;

  constructor() {
    this.home = pug.compileFile("templates/home.pug");
    this.login = pug.compileFile("templates/login.pug");
    console.log("Pug compiled");
  }
}
