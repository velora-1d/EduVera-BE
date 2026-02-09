const path = require("node:path");
const winston = require("winston");
require("winston-daily-rotate-file");

const rotate = new winston.transports.DailyRotateFile({
    filename: path.join(__dirname, "..", "logs", "app-%DATE%.log"),
    datePattern: "YYYY-MM-DD",
    zippedArchive: true,
    maxSize: "20m",
    maxFiles: "14d",
    level: "debug",
    format: winston.format.combine(
        winston.format.timestamp({ format: "YYYY-MM-DD HH:mm:ss" }),
        winston.format.json()
    ),
});

const logger = winston.createLogger({
    level: "debug",
    transports: [rotate],
});

rotate.on("error", (error) => {
    logger.error("Error in log file rotation:", error);
});

rotate.on("rotate", (oldFilename, newFilename) => {
    logger.info(`Rotated log file: ${oldFilename} -> ${newFilename}`);
});

module.exports = logger;
