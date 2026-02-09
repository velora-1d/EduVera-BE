const contacts = require("../models/contact");
const groups = require("../models/group");
const job = require("../models/job");

const displaySendMessageToUser = (_, res) => {
    const contactsToDisplay = contacts.paginate(undefined, 10, 1);
    const totalContacts = contacts.count();
    res.render("message", {
        title: "Send Message",
        pathname: "send-message",
        contactsToDisplay,
        totalContacts,
    });
};

const displaySendMessageToGroup = (_, res) => {
    const groupsToDisplay = groups.paginate(undefined, 10, 1);
    const totalGroups = groups.count();
    res.render("group", {
        title: "Send Group Message",
        pathname: "send-group-message",
        groupsToDisplay,
        totalGroups,
    });
};

const displayScheduleMessage = (_, res) => {
    const contactsToDisplay = contacts.paginate(undefined, 10, 1);
    const totalContacts = contacts.count();
    const groupsToDisplay = groups.paginate(undefined, 10, 1);
    const totalGroups = groups.count();
    const schedules = job.paginate(undefined, 10, 1);
    res.render("schedule", {
        title: "Schedule Message",
        pathname: "schedule-message",
        contacts: contactsToDisplay,
        total_contacts: totalContacts,
        groups: groupsToDisplay,
        total_groups: totalGroups,
        schedules,
    });
};

module.exports = {
    displaySendMessageToUser,
    displaySendMessageToGroup,
    displayScheduleMessage,
};
