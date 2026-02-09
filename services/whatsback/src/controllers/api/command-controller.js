const command = require("../../models/command"); // Assuming you have a command model

/**
 * GET /api/command
 * Retrieves all available commands.
 * @param {import("express").Request} req - The request object.
 * @param {import("express").Response} res - The response object.
 * @returns {Promise<void>}
 */
const getAllCommands = async (req, res) => {
    try {
        const commands = command.iterate();
        res.status(200).json({
            status: true,
            data: {
                commands,
            },
        });
    } catch (error) {
        res.status(500).json({
            status: false,
            message: `API Error - ${error}`,
        });
    }
};

/**
 * POST /api/command
 * Creates a new command.
 * @param {import("express").Request} req - The request object.
 * @param {import("express").Response} res - The response object.
 * @returns {Promise<void>}
 */
const createCommand = async (req, res) => {
    try {
        const { command: command_name, response } = req.body;
        const save = command.save({ command_name, response });

        if (save > 0) {
            res.status(200).json({
                status: true,
                message: "Command saved successfully",
            });
        } else {
            res.status(500).json({
                status: false,
                message: "Failed to save command",
            });
        }
    } catch (error) {
        res.status(500).json({
            status: false,
            message: `API Error - ${error}`,
        });
    }
};

/**
 * PUT /api/command/:command_id
 * Updates an existing command.
 * @param {import("express").Request} req - The request object.
 * @param {import("express").Response} res - The response object.
 * @returns {Promise<void>}
 */
const updateCommand = async (req, res) => {
    try {
        const { command: command_name, response } = req.body;
        const { command_id } = req.params;
        const update = command.update(command_id, { command_name, response });

        if (update > 0) {
            res.status(200).json({
                status: true,
                message: "Command updated successfully",
            });
        } else {
            res.status(500).json({
                status: false,
                message: "Failed to update command",
            });
        }
    } catch (error) {
        res.status(500).json({
            status: false,
            message: `API Error - ${error}`,
        });
    }
};

/**
 * DELETE /api/command/:command_id
 * Deletes an existing command.
 * @param {import("express").Request} req - The request object.
 * @param {import("express").Response} res - The response object.
 * @returns {Promise<void>}
 */
const deleteCommand = async (req, res) => {
    try {
        const { command_id } = req.params;
        const deleted = command.delete(command_id);

        if (deleted > 0) {
            res.status(200).json({
                status: true,
                message: "Command deleted successfully",
            });
        } else {
            res.status(500).json({
                status: false,
                message: "Failed to delete command",
            });
        }
    } catch (error) {
        res.status(500).json({
            status: false,
            message: `API Error - ${error}`,
        });
    }
};

module.exports = {
    getAllCommands,
    createCommand,
    updateCommand,
    deleteCommand,
};
