import express from "express";
import { renderSend } from "./helper/renderHelper";
import Pug from "./pug";

/**
 * A collection of functions that are run before the request is handled.
 */
export default class PreRequest {

    /**
     * Will check if the user has attempted to login too many times.
     * If so, it will show a 429 error.
     */
    public static loginAttempts(req: express.Request, res: express.Response, loginAttempt: string[number] | undefined[]) {
        const ip = req.socket.remoteAddress || req.headers['x-forwarded-for'][0];

        if (loginAttempt[ip] === undefined) loginAttempt[ip] = 0;
        loginAttempt[ip]++;

        if (loginAttempt[ip] > 5) {
            res.status(429).send("Too many login attempts.  Please try again later.");
            return false;
        }

        return true;
    }

    /**
     *  Will check if the user is logged in.  
     *  If not, it will redirect to the login page.
     */
    public static userSpace(req: express.Request, res: express.Response) {
        return true;
        res.redirect("/login");
        return false;
    }
}