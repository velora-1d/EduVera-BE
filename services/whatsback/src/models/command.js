const database = require("../database");

const table = "commands";

module.exports = {
    /**
   * Counts the total number of records in the "commands" table.
   * @returns {number} The total count of records in the table.
   */
    count: () => {
        const stmt = database.prepare(`SELECT COUNT(id) AS total FROM ${table}`);
        const result = stmt.get();
        return result.total;
    },

    /**
   * Retrieves all records from the "commands" table as an iterable.
   * Useful for handling large datasets efficiently.
   * @returns {IterableIterator<Object>} An iterator yielding each row as an object.
   */
    iterate: () => {
        const stmt = database.prepare(`SELECT * FROM ${table}`);
        const iterator = stmt.iterate();

        // Convert the iterator to an array
        const commands = [];
        for (const row of iterator) {
            commands.push(row);
        }

        return commands;
    },

    /**
   * Inserts a new record into the "commands" table.
   * @param {Object} data The record data to be inserted.
   * @param {string} data.command The command trigger text.
   * @param {string} data.response The response text when the command is triggered.
   * @returns {number} The number of rows inserted.
   */
    save: (data) => {
        const stmt = database.prepare(
            `INSERT INTO ${table} (command, response) VALUES (?, ?)`
        );
        const info = stmt.run(data.command_name, data.response);

        return info.changes;
    },

    /**
   * Updates an existing record in the "commands" table.
   * @param {number} id The ID of the record to update.
   * @param {Object} data The updated data for the record.
   * @param {string} data.command_name The new command trigger text.
   * @param {string} data.response The new response text when the command is triggered.
   * @returns {number} The number of rows affected by the update.
   */
    update: (id, data) => {
        const stmt = database.prepare(
            `UPDATE ${table} SET command = ?, response = ? WHERE id = ?`
        );
        const info = stmt.run(data.command_name, data.response, id);

        return info.changes;
    },

    /**
   * Deletes a record from the "commands" table.
   * @param {number} id The ID of the record to delete.
   * @returns {number} The number of rows affected by the deletion.
   */
    delete: (id) => {
        const stmt = database.prepare(`DELETE FROM ${table} WHERE id = ?`);
        const info = stmt.run(id);

        return info.changes;
    },
};
