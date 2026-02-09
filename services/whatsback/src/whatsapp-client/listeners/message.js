const database = require("../../database");
const { phoneNumberFormatter } = require("../../helper");

/**
 * Handles incoming messages for the WhatsApp client.
 * Responds to specific commands or messages with predefined replies.
 * @param {Message} msg - The incoming message object.
 * @param {Client} client - The WhatsApp client instance.
 * @returns {void}
 */
module.exports = async function messageHandler(message, client) {
    try {
        if (message.body.startsWith("!")) {
            const stmt = database.prepare(
                "SELECT command, response FROM commands WHERE command = ?"
            );
            const data = stmt.get(message.body);

            switch (message.body) {
            case data.command: {
                message.reply(data.response);

                break;
            }
            case "!whois": {
                message.reply(`I am ${client.info?.pushname || "unknown"}`);

                break;
            }
            case "!whoami": {
                const contact = await message.getContact();
                const contactName = contact.pushname || contact.name || "Unknown";
                const contactNumber = phoneNumberFormatter(contact.number);
                client.sendMessage(
                    contactNumber,
                    `Hi! Your name is *${contactName}* and your number is ${contact.number}`
                );

                break;
            }
            default: {
                console.warn("Command not found");
                message.reply("Command not found!");
            }
            }
        }
    } catch {
        console.error("Error handling message");
    }
};
