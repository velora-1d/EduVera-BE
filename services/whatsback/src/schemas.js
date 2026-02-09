const database = require("./database");
const crypto = require("node:crypto");

const create_commands_table = () => {
    const sql = "CREATE TABLE IF NOT EXISTS commands (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, command TEXT UNIQUE NOT NULL, response TEXT NOT NULL)";
    database.prepare(sql).run();
};

const create_contacts_table = () => {
    const sql = `CREATE TABLE IF NOT EXISTS contacts (
    number INTEGER PRIMARY KEY,
    name TEXT
  )`;
    database.prepare(sql).run();
};

const alter_profile_picture_in_contact_table = () => {
    const columns = database.prepare("PRAGMA table_info(contacts)").all();
    const columnExists = columns.some((col) => col.name === "profilePicture");

    if (!columnExists) {
        const contactImage = crypto.randomBytes(4).toString("hex");
        const sql = `
      ALTER TABLE contacts 
      ADD COLUMN profilePicture TEXT NOT NULL DEFAULT 'https://robohash.org/${contactImage}'
    `;
        database.prepare(sql).run();
    }
};

const create_groups_table = () => {
    const sql = `CREATE TABLE IF NOT EXISTS groups (
    groupId TEXT PRIMARY KEY,
    groupName TEXT,
    totalParticipants INTEGER
  )`;
    database.prepare(sql).run();
};

const create_jobs_table = () => {
    database
        .prepare(
            `
  CREATE TABLE IF NOT EXISTS jobs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    job_name TEXT,
    job_trigger TEXT,
    target_contact_or_group TEXT,
    message TEXT,
    job_cron_expression TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    deleted_at DATETIME
  )
`
        )
        .run();
};

const alter_job_status_in_jobs_table = () => {
    const columns = database.prepare("PRAGMA table_info(jobs)").all();
    const columnExists = columns.some((col) => col.name === "job_status");

    if (!columnExists) {
        const sql = `
      ALTER TABLE jobs 
      ADD COLUMN job_status INTEGER NOT NULL DEFAULT 1
    `;
        database.prepare(sql).run();
    }
};

const create_job_histories_table = () => {
    database
        .prepare(
            `
CREATE TABLE IF NOT EXISTS job_histories (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  job_name TEXT,
  job_execute_time DATETIME,
  job_complete_time DATETIME,
  job_error_message TEXT
)
`
        )
        .run();
};

const create_table_message_histories = () => {
    database.prepare(
        `CREATE TABLE IF NOT EXISTS message_histories (
  message_id INTEGER PRIMARY KEY AUTOINCREMENT,
  message_target TEXT NOT NULL,
  message_content TEXT,
  message_type TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME
);
`
    ).run();
};

const init = () => {
    create_commands_table();
    create_contacts_table();
    alter_profile_picture_in_contact_table();
    create_groups_table();
    create_jobs_table();
    alter_job_status_in_jobs_table();
    create_job_histories_table();
    create_table_message_histories();
};

module.exports = init;
