const contact = require("../models/contact");

const displayContacts = (_, res) => {
    const contacts = contact.iterate();
    const totalContacts = contact.count();
    res.render("contacts", { title: "Contacts", pathname: "contacts", contacts, totalContacts });
};

module.exports = {
    displayContacts,
};
