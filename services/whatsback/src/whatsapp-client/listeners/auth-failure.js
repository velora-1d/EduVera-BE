const { serverLog } = require("../../helper");

/**
 * Handles the "auth_failure" event for the WhatsApp client.
 * Logs the authentication failure and notifies all connected sockets.
 * @param {Array<Object>} connectedSockets - An array of connected sockets.
 */
module.exports = function authFailureHandler(connectedSockets) {
    serverLog("WhatsApp client failed to authenticate");
    for (const socket of connectedSockets) {
        socket.emit("logs", "Auth failure, restarting...");
    }
};
