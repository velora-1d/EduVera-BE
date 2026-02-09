const { serverLog } = require("../helper");

/**
 * Get current session profile info
 * @param {Client} client - The WhatsApp client instance.
 * @returns {Object}
 */
module.exports = async function getProfile(client) {
    try {
        const myId = client.info?.wid?._serialized;
        const phoneNumber = await client.getFormattedNumber(myId);
        const myContact = await client.getContactById(myId);

        const profilePicUrl =
      (await myContact.getProfilePicUrl()) || "https://robohash.org/unknown";
        const userName = client.info?.pushname || "unknown";

        return {
            name: userName,
            picture: profilePicUrl,
            phone: phoneNumber,
        };
    } catch {
        serverLog("Client is not ready yet!");
        return {
            name: "unknown",
            picture: "https://robohash.org/unknown",
        };
    }
};
