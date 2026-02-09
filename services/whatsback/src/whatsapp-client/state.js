/**
 * @module state
 * @property {string|undefined} lastQR - The last QR code data, or `undefined` if no QR code is available.
 * @property {boolean} isAuthenticated - Indicates whether the user is authenticated.
 * @property {boolean} isReady - Indicates whether the application is ready.
 */
module.exports = {
    lastQR: undefined,
    isAuthenticated: false,
    isReady: false,
};
