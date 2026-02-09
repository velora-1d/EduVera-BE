const express = require("express");
const { getPaginateGroup } = require("../../controllers/api/group-controller");
const router = express.Router();

// GET /api/groups
router.get("/", getPaginateGroup);

module.exports = router;