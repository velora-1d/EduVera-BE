const express = require("express");
const { getContacts } = require("../../controllers/api/contact-controller");
const router = express.Router();

// GET /api/contacts
router.get("/", getContacts);

module.exports = router;
