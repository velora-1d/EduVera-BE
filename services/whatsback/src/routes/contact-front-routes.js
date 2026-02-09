const express = require("express");
const { displayContacts } = require("../controllers/contact-front-controller");
const router = express.Router();

router.get("/", displayContacts);

module.exports = router;
