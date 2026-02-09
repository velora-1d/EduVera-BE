const express = require("express");
const Cron = require("../../libs/cron-parser");
const router = express.Router();

// GET /api/cron-next-runs
router.post("/", (req, res) => {
    const { exp } = req.body;
    const cronInstance = new Cron(exp, { tz: process.env.TZ || "UTC" });
    res.json({
        description: cronInstance.translate(),
        nexRuns: cronInstance.getNextRuns(),
    });
});

module.exports = router;
