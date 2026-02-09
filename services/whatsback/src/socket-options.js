const { serverLog, sleep } = require("./helper");
const qrcode = require("qrcode");
const state = require("./whatsapp-client/state");
const userInfo = require("./whatsapp-client/get-profile");

module.exports = async (socket, connectedSockets, client) => {
    connectedSockets.push(socket);
    socket.emit("connected", "WhatsApp Client is Connected!");
    socket.emit("logs", "WhatsApp Client is Connected!");
    serverLog("WhatsApp Client is Connected!");
    if (state.isAuthenticated) {
        socket.emit("is_authenticated", true);
    }

    if (state.lastQR) {
        qrcode.toDataURL(state.lastQR, (error, url) => {
            if (!error) {
                socket.emit("qr", url);
                socket.emit("logs", "QR Code received, scan please!");
                serverLog("QR Code sent to new socket");
            }
        });
    }
    if (state.isAuthenticated) {
        const info = await userInfo(client);
        socket.emit("authenticated", {
            log: "WhatsApp is authenticated!",
            user_info: info,
        });
        serverLog("Authenticated state sent to new socket");
    }
    if (state.isReady) {
        const info = await userInfo(client);
        socket.emit("ready", {
            log: "WhatsApp client is ready!",
            user_info: info,
        });
        socket.emit("logs", "WhatsApp client is ready!");
        serverLog("Ready state sent to new socket");
    }

    socket.on("logout", async () => {
        try {
            await client.logout();
            serverLog("Logged out the current client");
            socket.emit("disconnected", "Logged out the current client");
            await client.destroy();
            socket.emit("client_logout", "You're now logged out!");
            serverLog("Client destroyed after logout");
      
            state.isAuthenticated = false;
            state.isReady = false;
            state.lastQR = undefined;
      
            await sleep(3000);
            await client.initialize();
            serverLog("Reinitialized the client after logout");
            socket.emit("logs", "Reinitialized the client");
        } catch (error) {
            console.error("Logout error:", error);
        }
    });

    socket.on("disconnect", () => {
        const index = connectedSockets.indexOf(socket);
        if (index !== -1) {
            connectedSockets.splice(index, 1);
        }
        serverLog("Socket disconnected");
    });
};
