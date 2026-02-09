const logger = require("./logger");

/**
 * Delays execution for a specified amount of time.
 * @param {number} ms - The number of milliseconds to sleep.
 * @returns {Promise<void>} A promise that resolves after the specified time.
 */
const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

/**
 * @param {string} [type=""] - Log level type.
 *     When only one argument is provided, defaults to "debug".
 *     Valid levels: "debug", "info", "notice", "warning", "error",
 *     "crit", "alert", "emerg", or any custom string
 * @default ""
 */
const serverLog = (type = "", message) => {
    if (message === undefined) {
        message = type;
        type = "debug";
    }

    const now = new Date();
    const options = { timeZone: process.env.TZ || "UTC", hour12: false };
    const dateParts = now.toLocaleDateString("id-ID", options).split("/");
    const timeParts = now.toLocaleTimeString("id-ID", options).split(".");
    const formattedTimestamp = `${dateParts[2]}-${dateParts[1].padStart(
        2,
        "0"
    )}-${dateParts[0].padStart(2, "0")} ${timeParts[0].padStart(
        2,
        "0"
    )}:${timeParts[1].padStart(2, "0")}:${timeParts[2].padStart(2, "0")}`;

    switch (type) {
    case "debug": {
        logger.debug(message);
        break;
    }
    case "info": {
        logger.info(message);
        break;
    }
    case "notice": {
        logger.notice(message);
        break;
    }
    case "warning": {
        logger.warning(message);
        break;
    }
    case "error": {
        logger.error(message);
        break;
    }
    case "crit": {
        logger.crit(message);
        break;
    }
    case "alert": {
        logger.alert(message);
        break;
    }
    case "emerg": {
        logger.emerg(message);
        break;
    }
    default: {
        logger.debug(message);
        break;
    }
    }

    console.log(
        `${formattedTimestamp} - BACKEND_LOG: ${message}`.toUpperCase()
    );
};

/**
 * Formats a given phone number into a WhatsApp-compatible format.
 * @param {string} number - The phone number to format.
 * @returns {string} The formatted phone number.
 */
const phoneNumberFormatter = function (number) {
    let formatted = number.replaceAll(/\D/g, "");

    if (formatted.startsWith("0")) {
        formatted = "62" + formatted.slice(1);
    }

    if (!formatted.endsWith("@c.us")) {
        formatted += "@c.us";
    }

    return formatted;
};

/**
 * Removes duplicate contacts based on their number.
 * @param {Array} contacts - Array of contact objects.
 * @returns {Array} - Array of unique contact objects.
 */
const removeDuplicateContacts = (contacts) => {
    const uniqueContactsMap = new Map();

    for (const contact of contacts) {
        if (
            !uniqueContactsMap.has(contact.number) &&
            uniqueContactsMap.id?.server === "c.us"
        ) {
            uniqueContactsMap.set(contact.number, contact);
        }
    }

    return [...uniqueContactsMap.values()];
};

/**
 * Formats a given phone number into a WhatsApp-compatible international format.
 * @param {string} number - The phone number to format.
 * @returns {string} The formatted phone number.
 */
const formatInternationalPhoneNumber = (number) => {
    number = number.toString();

    if (!number.startsWith("62")) {
        return number;
    }

    let localNumber = number.slice(2);

    switch (localNumber.length) {
    case 10: {
        return `+62 ${localNumber.replace(
            /(\d{3})(\d{3})(\d{4})/,
            "$1-$2-$3"
        )}`;
    }
    case 11: {
        return `+62 ${localNumber.replace(
            /(\d{3})(\d{4})(\d{4})/,
            "$1-$2-$3"
        )}`;
    }
    case 12: {
        return `+62 ${localNumber.replace(
            /(\d{4})(\d{4})(\d{4})/,
            "$1-$2-$3"
        )}`;
    }
    case 13: {
        return `+62 ${localNumber.replace(
            /(\d{4})(\d{4})(\d{5})/,
            "$1-$2-$3"
        )}`;
    }
    case 14: {
        return `+62 ${localNumber.replace(
            /(\d{4})(\d{5})(\d{5})/,
            "$1-$2-$3"
        )}`;
    }
    case 15: {
        return `+62 ${localNumber.replace(
            /(\d{4})(\d{5})(\d{6})/,
            "$1-$2-$3"
        )}`;
    }
    default: {
        return "Invalid length for an Indonesian number";
    }
    }
};

/**
 * Takes a string of comma-separated origins and returns an array of valid origins.
 * If the input is empty, not a string, or contains only whitespace, returns "*".
 * @param {string} value - The string of comma-separated origins
 * @returns {string[]|string} - The array of valid origins or "*".
 */
const parseOrigins = (value) => {
    if (!value || typeof value !== "string") return "*";

    const origins = value
        .split(",")
        .map((o) => o.trim())
        .filter((o) => o.length > 0);

    return origins.length > 0 ? origins : "*";
};

/**
 * Converts a given string to snake_case.
 *
 * This function replaces spaces and uppercase letters with underscores, removes
 * non-word characters, and converts the entire string to lowercase.
 *
 * @param {string} text - The input string to be converted.
 * @returns {string} - The converted snake_case string.
 */

const toSnakeCase = (text) => {
    return text
        .replace(/([a-z])([A-Z])/g, "$1_$2")
        .replace(/\s+/g, "_")
        .replace(/[^\w]/g, "")
        .toLowerCase();
};

/**
 * Calculates the duration of a typing simulation based on the length of the input message.
 *
 * The calculation takes into account the average typing speed, the minimum and maximum delay,
 * and introduces a random variation to simulate human behavior. For longer messages, the
 * function adds a pause duration to simulate the natural pauses between words.
 *
 * @param {string} message - The input message to calculate the typing duration for.
 * @returns {number} - The calculated typing duration in milliseconds.
 */
const calculateTypingDuration = (message) => {
    const config = {
        avgSpeed: 225,
        minDelay: 800,
        maxDelay: 12000,
        speedVariation: 0.3,
        humanPauseThreshold: 75,
        pauseDuration: 1200,
    };

    let baseTime = (message.length / config.avgSpeed) * 60 * 1000;

    const speedVariation = 1 + (Math.random() * 2 - 1) * config.speedVariation;
    baseTime *= speedVariation;

    if (message.length > config.humanPauseThreshold) {
        const pauseCount = Math.floor(
            message.length / config.humanPauseThreshold
        );
        baseTime += pauseCount * config.pauseDuration;
    }

    return Math.min(config.maxDelay, Math.max(config.minDelay, baseTime));
};

/**
 * Prints a banner to the console with information about the running
 * Whatsback server instance.
 *
 * The banner includes the title, port number, REST API URL, Whatsback UI URL,
 * and links to support the project via GitHub Sponsors and Saweria.
 *
 * @param {number} PORT - The port number that the server is listening on.
 */
const bannerCLI = async (PORT) => {
    const { default: boxen } = await import("boxen");

    const title = "Whatsback Server is Running";
    const maxLength = Math.max(
        title.length,
        38 + PORT.toString().length,
        "[Support via GitHub Sponsors]".length + 23,
        "[Support via Saweria]".length + 15
    );
    const separator = "─".repeat(maxLength + 2);

    const UI_PORT =
        process.env.NODE_ENV === "production"
            ? process.env.UI_PORT || 8169
            : PORT;

    const text = [
        `\x1b[0m${title}\x1b[0m`,
        `\x1b[0m${separator}\x1b[0m`,
        `\x1b[0m📡 Running on: \x1b[1;35mhttp://localhost:${PORT}\x1b[0m`,
        `\x1b[0m🔌 REST API:   \x1b[1;32mhttp://localhost:${PORT}\x1b[0m`,
        `\x1b[0m💻 Whatsback UI: \x1b[1;32mhttp://localhost:${UI_PORT}\x1b[0m`,
        `\x1b[0m${separator}\x1b[0m`,
        "\x1b[0m💖 Support: \x1b[4;36mhttps://github.com/sponsors/darkterminal\x1b[0m",
        "\x1b[0m💖 Support: \x1b[4;36mhttps://saweria.co/darkterminal\x1b[0m",
    ].join("\n");

    const boxedBanner = boxen(text, {
        padding: 1,
        margin: 1,
        borderStyle: "round",
        borderColor: "cyan",
        align: "center",
    });

    console.log(boxedBanner);
};

module.exports = {
    sleep,
    serverLog,
    phoneNumberFormatter,
    removeDuplicateContacts,
    formatInternationalPhoneNumber,
    parseOrigins,
    toSnakeCase,
    calculateTypingDuration,
    bannerCLI,
};
