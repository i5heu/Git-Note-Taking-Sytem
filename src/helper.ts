import pug from 'pug';
import express from "express";

/**
 * Helper function to render and send the Pug template.
 */
export function renderSend(res: express.Response, compiledTemplate: pug.compileTemplate, data: any) {
    return res.send(compiledTemplate(data));
}