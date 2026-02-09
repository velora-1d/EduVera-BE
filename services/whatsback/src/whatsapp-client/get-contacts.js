const { socketEmit } = require("../../lib/socket-instance");
const { upsertMany } = require("../models/contact");

/**
 * Fetches all contacts from the WhatsApp client, filters out contacts
 * that are not valid users, and inserts them into the contacts table.
 * @param {Client} client - The WhatsApp client instance.
 */
module.exports = async function getContacts(client) {
    socketEmit("logs", "Sync your contacts...");
    const allContacts = await client.getContacts();
    const contacts = [];
    if (allContacts) {
        for (let contact of allContacts) {
            if (
                contact.id.server === "c.us" &&
        contact.isMe === false &&
        contact.isUser === true &&
        contact.isGroup === false &&
        contact.isWAContact === true &&
        contact.isBlocked === false
            ) {
                contacts.push({
                    name: contact.name,
                    number: contact.number,
                    profilePicture: `https://robohash.org/${contact.number}`,
                });
            }
        }

        upsertMany(contacts);
    }
    socketEmit("logs", "Synced your contacts!");
};
