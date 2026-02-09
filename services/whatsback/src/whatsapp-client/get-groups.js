const { upsertMany } = require("../models/group");

/**
 * Retrieves all group chats from the WhatsApp client, transforms them into group objects,
 * and inserts or replaces them in the database.
 * 
 * @param {Client} client - The WhatsApp client instance.
 * @async
 * @returns {Promise<void>} - A promise that resolves when the operation is complete.
 */

module.exports = async function getGroups(client) {
    const chats = await client.getChats();
    const groupChats = chats.filter((chat) => chat.isGroup);
    const groups = groupChats.map((group) => ({
        groupId: group.id._serialized,
        groupName: group.name,
        totalParticipants: group.participants.length,
    }));

    upsertMany(groups);
};

