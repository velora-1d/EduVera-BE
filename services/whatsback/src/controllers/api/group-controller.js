const groups = require("../../models/group");

/**
 * GET /api/groups
 * @summary Fetch a page of groups (with or without a search term)
 * @param {string} [search=""] - Search term to filter groups by name.
 * @param {number} [perPage=10] - Number of groups per page.
 * @param {number} [page=1] - Page number to fetch.
 * @returns {object} - JSON object with a `status` property (boolean), a `message` property (string), and a `data` property (array of group objects).
 * @throws {Error} - If there is an error with the database query.
 */
const getPaginateGroup = async (req, res) => {
    try {
    // Pagination variables
        const terms = req.query.search || "";
        const perPage = req.query.perPage || 10; // Number of groups per page
        let page = req.query.page || 1; // Current page number

        const groupsToDisplay = groups.paginate(terms, perPage, page);

        res.status(200).json({
            status: true,
            message: `Fetch group page ${page}`,
            data: groupsToDisplay,
        });
    } catch (error) {
        res.status(500).json({
            status: false,
            message: `API Error - ${error}`,
        });
    }
};

module.exports = {
    getPaginateGroup,
};
