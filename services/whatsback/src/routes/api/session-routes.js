const express = require('express');
const router = express.Router();
const state = require('../../whatsapp-client/state');
const { client } = require('../../whatsapp-client');
const qrcode = require('qrcode');

// Get Session Status and QR Code
// Used by EduVera Backend Adapter to mimic Evolution API behavior
router.get('/status', async (req, res) => {
    try {
        let status = 'disconnected';
        if (state.isAuthenticated) {
            status = 'connected';
        } else if (state.isReady || (state.lastQR && !state.isAuthenticated)) {
            // If QR is available but not authenticated, it's connecting (waiting specifically for scan)
            status = 'connecting';
        }

        let qrBase64 = null;
        if (state.lastQR && !state.isAuthenticated) {
            // Convert QR string to Base64 Image
            qrBase64 = await qrcode.toDataURL(state.lastQR);
        }

        // Get phone number from client info if available
        let phoneNumber = null;
        if (state.isAuthenticated && client.info && client.info.wid) {
            phoneNumber = client.info.wid.user;
        }

        res.json({
            status: true,
            data: {
                status: status,
                qrcode: qrBase64,
                phone_number: phoneNumber,
                instance_name: "default"
            }
        });
    } catch (error) {
        console.error("Error in session status:", error);
        res.status(500).json({ status: false, message: error.message });
    }
});

// Logout Session
router.post('/logout', async (req, res) => {
    try {
        if (state.isAuthenticated) {
            await client.logout();
            // Reset state
            state.isAuthenticated = false;
            state.lastQR = undefined;
            res.json({ status: true, message: "Logged out successfully" });
        } else {
            res.status(400).json({ status: false, message: "Not authenticated" });
        }
    } catch (error) {
        console.error("Error logging out:", error);
        res.status(500).json({ status: false, message: error.message });
    }
});

module.exports = router;
