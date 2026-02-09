const database = require("../database");
const { serverLog } = require("../helper");
const contact = require("./contact");
const group = require("./group");

const table = "message_histories";

module.exports = {
    /**
   * Records a message history.
   * @param {string} messageTarget - A string in the format of a phone number
   *   (e.g. "+1234567890") or a group ID (e.g. "1234567890-123456@g.us").
   * @param {string} messageContent - The message that was sent.
   * @param {string} messageType - The type of message (e.g. "text", "image", etc.).
   */
    recordMessageHistory: (messageTarget, messageContent, messageType) => {
        try {
            if (messageType !== "DIRECT_MESSAGE" && messageType !== "GROUP_MESSAGE") {
                return;
            }

            if (!messageTarget || !messageContent) {
                return;
            }

            if (messageType === "GROUP_MESSAGE") {
                const dataGroup = group.find(messageTarget);
                messageTarget = `${dataGroup.groupName}|${dataGroup.groupId}`;
            } else {
                let searchTarget = messageTarget.replace("@c.us", "");
                const dataContact = contact.find(searchTarget);
                messageTarget = `${dataContact.name}|${dataContact.number}`;
            }

            const stmt = database.prepare(
                `INSERT INTO ${table} (message_target, message_content, message_type) VALUES (?, ?, ?)`
            );
            stmt.run(messageTarget, messageContent, messageType);
            serverLog(`Message history recorded for ${messageTarget}`);
        } catch (error) {
            serverLog("Error recording message history:", error);
        }
    },

    /**
   * Returns the total number of direct messages recorded in the message history
   * for the current month.
   * @returns {number} - Total count of direct messages for the current month.
   */
    countDirectMessage: () => {
        const stmt = database.prepare(
            `SELECT COUNT(*) AS count FROM ${table} WHERE message_type = 'DIRECT_MESSAGE' AND strftime('%m', created_at) = strftime('%m', 'now')`
        );
        return stmt.get().count;
    },

    /**
   * Returns the total number of group messages recorded in the message history
   * for the current month.
   * @returns {number} - Total count of group messages for the current month.
   */
    countGroupMessage: () => {
        const stmt = database.prepare(
            `SELECT COUNT(*) AS count FROM ${table} WHERE message_type = 'GROUP_MESSAGE' AND strftime('%m', created_at) = strftime('%m', 'now')`
        );
        return stmt.get().count;
    },
};
