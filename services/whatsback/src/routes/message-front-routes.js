const express = require("express");
const {
    displaySendMessageToUser,
    displaySendMessageToGroup,
    displayScheduleMessage,
} = require("../controllers/message-front-controller");
const router = express.Router();

router.get("/send", displaySendMessageToUser);
router.get("/send-group", displaySendMessageToGroup);
router.get("/schedule", displayScheduleMessage);

module.exports = router;
