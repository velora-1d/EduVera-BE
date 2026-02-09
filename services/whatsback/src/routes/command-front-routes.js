const express = require("express");
const { displayCommand } = require("../controllers/command-front-controller");
const router = express.Router();

router.get("/", displayCommand);

module.exports = router;
