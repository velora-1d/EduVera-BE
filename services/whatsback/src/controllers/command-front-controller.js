const command = require("../models/command");

const displayCommand = (_, res) => {
    const commands = command.iterate();
    const totalCommand = command.count();
    res.render("commands", { title: "Commands", commands, totalCommand, pathname: "commands" });
};

module.exports = {
    displayCommand,
};