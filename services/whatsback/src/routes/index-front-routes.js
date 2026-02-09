const express = require("express");
const message_history = require("../models/message_history");
const router = express.Router();

router.get("/", (req, res) => {
    const totalDirectMessages = message_history.countDirectMessage();
    const totalGroupMessages = message_history.countGroupMessage();

    res.render("index", {
        pathname: "dashboard",
        totalDirectMessages,
        totalGroupMessages,
    });
});

router.get("/user-manual", (req, res) => {
    res.render("user-manual", {
        pathname: "user-manual",
        title: "User Manual",
    });
});

router.get("/integration", (req, res) => {
    res.render("integration", {
        pathname: "integration",
        title: "Integration",
    });
});

module.exports = router;
