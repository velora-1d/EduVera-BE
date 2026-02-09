const database = require("./src/database");
const cron = require("node-cron");
const { serverLog } = require("./src/helper");
require("dotenv").config();

const APP_HOST =
  process.env.NODE_ENV === "production"
      ? process.env.APP_HOST || "app"
      : "localhost";
const APP_PORT = process.env.APP_PORT || 5001;
const APP_URL = `http://${APP_HOST}:${APP_PORT}`;

const scheduledJobs = new Map();

/**
 * Extracts the name and number from a given job object.
 *
 * @param {Object} job - The job object containing the target contact or group.
 * @returns {[string, string]} An array containing the name and number extracted from the job.
 */
function getDetailContact(job) {
    const detail = job.target_contact_or_group.split("|");
    const name = detail[0];
    const number = detail[1];

    return [name, number];
}

/**
 * Sends a direct message to a specified contact using the WhatsApp API.
 *
 * @param {Object} job - The job object containing details for the message to be sent.
 * @param {string} job.message - The message content to be sent.
 * @param {string} job.target_contact_or_group - The target contact information in the format "name|number".
 *
 * Logs the status of the message sending operation, indicating success or failure.
 */
async function sendMessage(job) {
    const [name, number] = getDetailContact(job);
    const response = await fetch(`${APP_URL}/api/message/send-message`, {
        method: "post",
        headers: {
            "Content-Type": "application/json",
            "WHATSBACK-SOURCE": "cronjob",
        },
        body: JSON.stringify({
            number,
            message: job.message,
        }),
    });
    const result = await response.json();
    if (result.status) {
        serverLog(`send_message: Sending "${job.message}" to ${name} ${number}`);
        return;
    }

    serverLog(
        `send_message: Failed to send "${job.message}" to ${name} ${number}`
    );
}

/**
 * Sends a group message to a specified contact using the WhatsApp API.
 *
 * @param {Object} job - The job object containing details for the message to be sent.
 * @param {string} job.message - The message content to be sent.
 * @param {string} job.target_contact_or_group - The target contact information in the format "name|number".
 *
 * Logs the status of the message sending operation, indicating success or failure.
 */
async function sendGroupMessage(job) {
    const [name, number] = getDetailContact(job);
    const response = await fetch(`${APP_URL}/api/message/send-group-message`, {
        method: "post",
        headers: {
            "Content-Type": "application/json",
            "WHATSBACK-SOURCE": "cronjob",
        },
        body: JSON.stringify({
            number,
            message: job.message,
        }),
    });
    const result = await response.json();
    if (result.status) {
        serverLog(
            `send_group_message: Sending "${job.message}" to ${name} ${number}`
        );
        return;
    }

    serverLog(
        `send_group_message: Failed to send "${job.message}" to ${name} ${number}`
    );
}

/**
 * Logs a job history entry.
 *
 * @param {Object} job - The job object containing details for the job to be logged.
 * @param {string} job.job_name - The name of the job.
 * @param {string} executeTime - The time at which the job was executed.
 * @param {string} completeTime - The time at which the job completed.
 * @param {string} [errorMessage] - An error message if the job failed.
 *
 * Logs the job history entry and any errors encountered.
 */
function logJobHistory(job, executeTime, completeTime, errorMessage) {
    try {
        const stmt = database.prepare(`
      INSERT INTO job_histories (job_name, job_execute_time, job_complete_time, job_error_message)
      VALUES (?, ?, ?, ?)
    `);
        stmt.run(job.job_name, executeTime, completeTime, errorMessage);
        serverLog(`Job history logged for job "${job.job_name}"`);
    } catch (error) {
        serverLog(`Error logging history for job ${job.id}:`, error.message);
    }
}

/**
 * Schedules a job to be executed at the specified cron expression time.
 *
 * @param {Object} job - The job object containing details for the job to be scheduled.
 * @param {string} job.job_name - The name of the job.
 * @param {string} job.job_cron_expression - The cron expression for scheduling the job.
 * @param {string} job.job_trigger - The trigger type for the job (e.g. "send_message").
 *
 * Logs the scheduling of the job, and any errors encountered while executing the job.
 */
function scheduleJob(job) {
    if (!cron.validate(job.job_cron_expression)) {
        serverLog(
            `Invalid cron expression for job "${job.job_name}" (ID: ${job.id})`
        );
        return;
    }

    if (scheduledJobs.has(job.id)) {
        return;
    }

    const task = cron.schedule(job.job_cron_expression, () => {
        const executeTime = new Date().toISOString();
        serverLog(
            `Executing job "${job.job_name}" (ID: ${job.id}) at ${executeTime}`
        );

        try {
            if (job.job_trigger === "send_message") {
                sendMessage(job);
            } else if (job.job_trigger === "send_group_message") {
                sendGroupMessage(job);
            } else {
                console.warn(
                    `Unknown trigger "${job.job_trigger}" for job "${job.job_name}"`
                );
            }

            const completeTime = new Date().toISOString();
            logJobHistory(job, executeTime, completeTime);
        } catch (error) {
            const completeTime = new Date().toISOString();
            logJobHistory(job, executeTime, completeTime, error.message);
            serverLog(
                `Error executing job "${job.job_name}" (ID: ${job.id}):`,
                error.message
            );
        }
    });

    scheduledJobs.set(job.id, task);
    serverLog(
        `Scheduled job "${job.job_name}" (ID: ${job.id}) with cron expression: ${job.job_cron_expression}`
    );
}

/**
 * Fetches all scheduled jobs from the database and schedules them.
 *
 * @private
 */
function loadJobs() {
    try {
        const rows = database
            .prepare(
                "SELECT * FROM jobs WHERE target_contact_or_group IS NOT NULL AND message IS NOT NULL AND job_status = 1 AND deleted_at IS NULL AND job_status = 1"
            )
            .all();
        for (let job of rows) {
            if (!scheduledJobs.has(job.id)) {
                scheduleJob(job);
            }
        }
    } catch (error) {
        serverLog("Error fetching jobs:", error.message);
    }
}

/**
 * Starts the cron job service by performing a health check and scheduling
 * jobs.
 *
 * Performs a GET request to the app's health check endpoint every 5 seconds
 * until it succeeds. If the health check fails 5 times, the process exits with
 * code 1.
 *
 * Once the health check succeeds, schedules all jobs in the database and
 * starts the interval to load jobs every 20 seconds.
 */
async function startCronJobs() {
    serverLog("Cron job started!");

    let retries = 5;
    while (retries > 0) {
        try {
            const response = await fetch(`${APP_URL}/health`, {
                method: "GET",
                timeout: 5000,
            });

            if (response.ok) {
                serverLog("Health check successful");
                break;
            }
        } catch (error) {
            serverLog(
                `Health check failed (${retries} retries left):`,
                error.message
            );
            retries--;
            await new Promise((resolve) => setTimeout(resolve, 5000));
        }
    }

    if (retries === 0) {
        serverLog("Failed to connect to app after 5 attempts");
        process.exit(1);
    }

    setInterval(loadJobs, 20000);
    loadJobs();
}

startCronJobs();
