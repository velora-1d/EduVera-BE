const {
    phoneNumberFormatter,
    calculateTypingDuration,
    serverLog,
} = require("../../helper");
const group = require("../../models/group");
const { recordMessageHistory } = require("../../models/message_history");
const { client } = require("../../whatsapp-client");

/**
 * Simulates typing indication in a chat for a specified duration.
 *
 * @param {string} phoneNumber - The phone number of the chat where the typing state is to be set.
 * @param {string} message - The message content used to calculate the typing duration.
 */

async function typingMessage(phoneNumber, message) {
    const typingDuration = calculateTypingDuration(message);

    const chat = await client.getChatById(phoneNumber);
    await chat.sendStateTyping();

    await new Promise((resolve) => setTimeout(resolve, typingDuration));

    await chat.clearState();
}

/**
 * Sends a direct message to a specified WhatsApp user.
 *
 * @param {object} req - The HTTP request object.
 * @param {object} res - The HTTP response object.
 *
 * @prop {string} req.body.number - The recipient's phone number.
 * @prop {string} req.body.message - The text message to send.
 *
 * @returns {Promise<void>} - A promise that resolves if the message is sent successfully.
 */
const sendMessageToUser = async (req, res) => {
    try {
        const { number, message } = req.body;

        const phoneNumber = phoneNumberFormatter(number);

        await typingMessage(phoneNumber, message);

        await client.sendMessage(phoneNumber, message);
        recordMessageHistory(phoneNumber, message, "DIRECT_MESSAGE");

        res.status(200).json({
            status: true,
            message: "Message sent",
        });
    } catch (error) {
        res.status(500).json({
            status: false,
            message: `API Error - ${error}`,
        });
    }
};

/**
 * Sends a message to a specified WhatsApp group.
 *
 * @param {object} req - The HTTP request object.
 * @param {object} res - The HTTP response object.
 *
 * @prop {string} req.body.groupId - The WhatsApp group ID (must end with '@g.us').
 * @prop {string} req.body.message - The text message to send.
 *
 * @returns {Promise<void>} - A promise that resolves if the message is sent successfully.
 */
const sendMessageToGroup = async (req, res) => {
    try {
        let { groupId, message } = req.body;

        if (!groupId.endsWith("@g.us")) {
            let detailGroup = group.findByName(groupId);

            if (!detailGroup) {
                serverLog("send_grou_message: Group not found");
                res.status(404).json({
                    status: false,
                    message: "Group not found",
                });
                return;
            }

            groupId = detailGroup.groupId;
        }

        await typingMessage(groupId, message);

        await client.sendMessage(groupId, message);
        recordMessageHistory(groupId, message, "GROUP_MESSAGE");

        serverLog("send_group_message: Message sent to group");
        res.status(200).json({
            status: true,
            message: "Message sent to group",
        });
    } catch (error) {
        serverLog("send_group_message: Error sending message to group");
        res.status(500).json({
            status: false,
            message: `API Error - ${error.message}`,
        });
    }
};

module.exports = { sendMessageToUser, sendMessageToGroup };
