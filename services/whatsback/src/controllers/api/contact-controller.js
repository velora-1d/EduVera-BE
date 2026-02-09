const contacts = require("../../models/contact");

/**
 * GET /api/contacts
 * @summary Fetch a page of contacts (with or without a search term)
 * @param {string} [search=""] - Search term to filter contacts by name.
 * @param {number} [perPage=10] - Number of contacts per page.
 * @param {number} [page=1] - Page number to fetch.
 * @returns {object} - JSON object with a `status` property (boolean), a `message` property (string), and a `data` property (object with `contacts` array and `totalContacts` number).
 * @throws {Error} - If there is an error with the database query.
 */
const getContacts = async (req, res) => {
    try {
        const searchTerm = req.query.search || "";
        const perPage = req.query.perPage || 10;
        const page = req.query.page || 1;

        const contactsToDisplay = contacts.paginate(searchTerm, perPage, page);
        const totalContacts = contacts.count();
        res.status(200).json({
            status: true,
            message: "Contacts fetched successfully",
            data: {
                contacts: contactsToDisplay,
                totalContacts,
            },
        });
    } catch (error) {
        console.log(error);
        res.status(500).json({
            status: false,
            message: "Internal Server Error",
        });
    }
};

module.exports = {
    getContacts,
};
