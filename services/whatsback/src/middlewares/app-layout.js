require("dotenv").config();
const state = require("../whatsapp-client/state");
const currentVersion = require("../../package.json").version;

/**
 * Sets app-wide layout variables.
 *
 * @param {Object} req - Express.js request object.
 * @param {Object} res - Express.js response object.
 * @param {Function} next - Express.js next middleware function.
 */
module.exports = (req, res, next) => {
    res.locals.APP_PORT = process.env.APP_PORT || 5001;
    res.locals.AUTHENTICATED = state.isAuthenticated;
    res.locals.NODE_ENV = process.env.NODE_ENV || "development";
    res.locals.APP_VERSION = currentVersion;

    res.locals.layout = state.isAuthenticated
        ? "./layouts/dashboard"
        : "./layouts/auth";

    next();
};
