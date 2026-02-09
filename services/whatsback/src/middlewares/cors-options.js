const { parseOrigins } = require("../helper");
require("dotenv").config();

/**
 * Additional CORS middleware to validate the origin of incoming requests,
 * allowing only requests from hosts specified in the API_CORS_ORIGIN environment variable.
 * Requests from cronjobs are exempt from this validation.
 *
 * @param {import("express").Request} req - The request object.
 * @param {import("express").Response} res - The response object.
 * @param {import("express").NextFunction} next - The next middleware function.
 * @returns {void} - Responds with a 403 status code if the origin is not allowed, otherwise calls the next middleware.
 */
const additionalCors = (req, res, next) => {
    const allowedOrigins = parseOrigins(process.env.API_CORS_ORIGIN);
    const origin = req.headers.origin;
    const host = req.headers.host;
    const isCronjob = req.headers["whatsback-source"] === "cronjob";

    if (isCronjob) return next();

    if (!origin) {
        const expectedHost =
        process.env.NODE_ENV === "production"
            ? `https://${host}`
            : `http://${host}`;

        if (!allowedOrigins.includes(expectedHost)) {
            return res.status(403).end();
        }
    }

    if (origin && !allowedOrigins.includes(origin)) {
        return res.status(403).end();
    }

    next();
};

const apiCors = {
    origin: parseOrigins(process.env.API_CORS_ORIGIN),
    allowedHeaders: ["Content-Type"],
    methods: ["GET", "POST", "OPTIONS", "DELETE", "PUT"],
};

module.exports = {
    additionalCors,
    apiCors
};
