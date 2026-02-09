const express = require("express");
const router = express.Router();
const validateRequestBody = require("../../middlewares/validate-request-body");
const {
    getAllCommands,
    createCommand,
    updateCommand,
    deleteCommand,
} = require("../../controllers/api/command-controller");

// GET /api/command
router.get("/", getAllCommands);

// POST /api/command
router.post(
    "/",
    validateRequestBody({ command: "string", response: "string" }),
    createCommand
);

// PUT /api/command/:command_id
router.put(
    "/:command_id",
    validateRequestBody({ command: "string", response: "string" }),
    updateCommand
);

// DELETE /api/command/:command_id
router.delete("/:command_id", deleteCommand);

module.exports = router;
