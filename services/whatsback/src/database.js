const Database = require("better-sqlite3");
require("dotenv").config();

const dbPath =
  process.env.NODE_ENV === "production"
      ? process.env.DB_PATH
      : "./whatsback.db";

const database = new Database(dbPath, {
    busyTimeout: 7000,
});
database.pragma("journal_mode = WAL");

module.exports = database;
