const qrcode = require("qrcode");
const { serverLog } = require("../../helper");

/**
 * Handles the "qr" event for the WhatsApp client.
 * Logs the QR code reception, updates the state, and notifies all connected sockets with the QR code data URL.
 * @param {string} qr - The QR code string.
 * @param {Array<Object>} connectedSockets - An array of connected sockets.
 * @param {Object} state - The application state object.
 * @returns {void}
 */
module.exports = function qrCodeHandler(qr, connectedSockets, state) {
    serverLog("QR Code is received");
    state.lastQR = qr;
    for (const socket of connectedSockets) {
        qrcode.toDataURL(qr, (error, url) => {
            if (!error) {
                socket.emit("qr", url);
                socket.emit("logs", "QR Code received, scan please!");
            }
        });
    }
};
