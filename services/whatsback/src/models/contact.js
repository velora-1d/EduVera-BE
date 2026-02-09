const database = require("../database");
const { serverLog } = require("../helper");

const table = "contacts";

module.exports = {
    /**
     * Retrieves a contact from the database by its number.
     * @param {string} number - The phone number of the contact to retrieve.
     * @returns {Object} The contact object, or undefined if not found.
     */
    find: (number) => {
        const stmt = database.prepare(
            `SELECT * FROM ${table} WHERE number = ?`
        );
        return stmt.get(number);
    },

    /**
     * Returns the total number of contacts in the table.
     * @returns {number} - Total number of contacts.
     */
    count: () => {
        const stmt = database.prepare(
            `SELECT COUNT(number) AS total FROM ${table}`
        );
        const result = stmt.get();
        return result.total;
    },

    /**
     * Returns an iterator for all contacts in the table.
     * @returns {Iterator} - Iterator for all contacts.
     */
    iterate: () => {
        const stmt = database.prepare(
            `SELECT * FROM ${table} ORDER BY name ASC`
        );
        const iterator = stmt.iterate();

        // Convert the iterator to an array
        const contacts = [];
        for (const row of iterator) {
            contacts.push(row);
        }

        return contacts;
    },

    /**
     * Paginates contacts from the table.
     * @param {string} search - Search contact name
     * @param {number} limit - Number of contacts per page.
     * @param {number} page - Number of contacts to skip.
     * @returns {Array} - Array of paginated contacts.
     */
    paginate: (search = "", limit = 10, page = 0) => {
        const offset = (page - 1) * limit;

        if (search) {
            let sql = `SELECT * FROM ${table} WHERE name LIKE '%${search}%' ORDER BY name ASC LIMIT ? OFFSET ?`;
            const stmt = database.prepare(sql);
            return stmt.all(limit, offset);
        }

        let sql = `SELECT * FROM ${table} ORDER BY name ASC LIMIT ? OFFSET ?`;
        const stmt = database.prepare(sql);
        return stmt.all(limit, offset);
    },

    /**
     * Inserts or replaces multiple contacts in the table.
     * @param {Array} contacts - Array of contact objects with `name` and `number` properties.
     */
    insertOrReplaceMany: (contacts) => {
        const insertOrReplace = database.prepare(`
      INSERT OR REPLACE INTO ${table} (name, number)
      VALUES (@name, @number)
    `);

        const insertTransaction = database.transaction((contacts) => {
            for (const contact of contacts) {
                insertOrReplace.run(contact);
            }
        });

        try {
            insertTransaction(contacts);
            serverLog(
                `${contacts.length} contacts inserted or replaced successfully.`
            );
        } catch (error) {
            serverLog("Error inserting or replacing contacts:", error);
            throw error;
        }
    },

    /**
     * Upserts multiple contacts in the table.
     * @param {Array} contacts - Array of contact objects with `name` and `number` properties.
     */
    upsertMany: (contacts) => {
        const upsertStatement = database.prepare(`
      INSERT INTO contacts (name, number) 
      VALUES (@name, @number)
      ON CONFLICT(number) DO UPDATE SET name = excluded.name
    `);

        // Define transaction once
        const upsertTransaction = database.transaction((contacts) => {
            for (const contact of contacts) {
                upsertStatement.run(contact);
            }
        });

        if (!contacts.length) return;

        try {
            upsertTransaction(contacts);
            serverLog(`${contacts.length} contacts upserted successfully.`);
        } catch (error) {
            serverLog("Error upserting contacts:", error);
            console.error(error);
        }
    },
};
