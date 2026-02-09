/**
 * Returns a middleware that validates the request body
 * by checking if the given requiredFields are present
 * and are of the correct type.
 *
 * @param {Object<string, string>} requiredFields
 * @returns {Function} express middleware
 */
module.exports = function validateRequestBody(requiredFields) {
    return (req, res, next) => {
        const errors = [];

        for (const [field, type] of Object.entries(requiredFields)) {
            if (!(field in req.body)) {
                errors.push(`${field} is required.`);
            } else if (type && typeof req.body[field] !== type) {
                errors.push(`${field} must be of type ${type}.`);
            }
        }

        if (errors.length > 0) {
            return res.status(400).json({ errors });
        }

        next();
    };
};
