const { Client, LocalAuth } = require("whatsapp-web.js");
require("dotenv").config();

let puppeteerOptions = {
    headless: true,
    args: [
        "--disable-gpu",
        "--disable-dev-shm-usage",
        "--disable-setuid-sandbox",
        "--no-first-run",
        "--no-sandbox",
        "--no-zygote",
        "--deterministic-fetch",
        "--disable-features=IsolateOrigins",
        "--disable-site-isolation-trials"
    ],
};

if (process.env.NODE_ENV === "production") {
    puppeteerOptions = {
        executablePath: process.env.PUPPETEER_EXECUTABLE_PATH || "/usr/bin/google-chrome-stable",
        headless: true,
        args: [
            "--disable-gpu",
            "--disable-dev-shm-usage",
            "--disable-setuid-sandbox",
            "--no-first-run",
            "--no-sandbox",
            "--no-zygote",
            "--deterministic-fetch",
            "--disable-features=IsolateOrigins",
            "--disable-site-isolation-trials",
            "--single-process",
        ],
    };
}

/**
 * WhatsApp client instance configured with LocalAuth and custom Puppeteer settings.
 * @type {Client}
 */
const client = new Client({
    authStrategy: new LocalAuth(),
    restartOnAuthFail: true,
    takeoverOnConflict: true,
    puppeteer: puppeteerOptions,
    qrMaxRetries: 10,
});

/**
 * Array of connected sockets.
 * @type {Array<Object>}
 */
let connectedSockets = [];

// Import event listeners (each is in its own module)
const readyHandler = require("./listeners/ready");
const authenticatedHandler = require("./listeners/authenticated");
const authFailureHandler = require("./listeners/auth-failure");
const qrCodeHandler = require("./listeners/qr");
const disconnectedHandler = require("./listeners/disconnected");
const messageHandler = require("./listeners/message");

const state = require("./state");

// Register WhatsApp client event listeners, passing in the client, sockets, and state.
client.on("ready", () => readyHandler(client, connectedSockets, state));
client.on("authenticated", () => authenticatedHandler(client, connectedSockets, state));
client.on("auth_failure", () => authFailureHandler(connectedSockets));
client.on("qr", (qr) => qrCodeHandler(qr, connectedSockets, state));
client.on("disconnected", (reason) => disconnectedHandler(reason, client, connectedSockets));
client.on("message", (message) => messageHandler(message, client));

client.initialize();

/**
 * Updates the reference to the connected sockets.
 * @param {Array<Object>} sockets - The new array of connected sockets.
 * @returns {void}
 */
const setSocketManager = (sockets) => {
    connectedSockets = sockets;
};

/**
 * Exported module containing the WhatsApp client and a setter for the socket manager.
 * @module WhatsAppClient
 * @property {Client} client - The WhatsApp client instance.
 * @property {Function} setSocketManager - Function to update the connected sockets.
 */
module.exports = { client, setSocketManager };
