const database = require("../database");
const { serverLog, toSnakeCase } = require("../helper");

const table = "jobs";

module.exports = {
    /**
     * Create a new job
     *
     * @param {string} jobName Name of the job
     * @param {string} jobTrigger Trigger for the job (e.g. "send_message")
     * @param {string} targetContactOrGroup Contact or Group to target
     * @param {string} message Message to send
     * @param {string} jobCronExpression Cron expression for scheduling
     * @returns {number} ID of the created job
     */
    create(
        jobName,
        jobTrigger,
        targetContactOrGroup,
        message,
        jobCronExpression
    ) {
        try {
            const stmt = database.prepare(`
        INSERT INTO ${table} (
          job_name,
          job_trigger,
          target_contact_or_group,
          message,
          job_cron_expression,
          job_status
        ) VALUES (
          @jobName,
          @jobTrigger,
          @targetContactOrGroup,
          @message,
          @jobCronExpression,
          @jobStatus
        )
      `);

            const info = stmt.run({
                jobName,
                jobTrigger,
                targetContactOrGroup,
                message,
                jobCronExpression,
                jobStatus: 1,
            });

            return info.lastInsertRowid;
        } catch (error) {
            serverLog("Error creating job:", error.message);
            return 0;
        }
    },

    /**
     * Finds a job by its ID
     * @param {number} jobId Job ID to find
     * @returns {Object} Job object if found, otherwise null
     */
    findById(jobId) {
        try {
            const stmt = database.prepare(`
        SELECT * FROM ${table} WHERE id = ?
      `);

            const result = stmt.get(jobId);

            return result || null;
        } catch (error) {
            serverLog("Error finding job by ID:", error.message);
            return null;
        }
    },

    /**
     * Finds jobs by its status
     * @param {number} jobStatus Job status to find (0 = disabled, 1 = enabled)
     * @returns {Array} Array of job objects if found, otherwise empty array
     */
    findByStatus(jobStatus) {
        try {
            const stmt = database.prepare(`
        SELECT * FROM ${table} WHERE job_status = ?
      `);

            const result = stmt.all(jobStatus);

            return result || [];
        } catch (error) {
            serverLog("Error finding job by status:", error.message);
            return [];
        }
    },

    /**
     * Soft deletes a job by its ID
     * @param {number} jobId Job ID to delete
     * @returns {Object} Job object if deleted, otherwise null
     */
    softDeleteById(jobId) {
        try {
            const stmt = database.prepare(`
        UPDATE ${table} SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?
      `);

            const result = stmt.run(jobId);

            return result.changes > 0 ? this.findById(jobId) : null;
        } catch (error) {
            serverLog("Error soft deleting job:", error.message);
            return null;
        }
    },

    /**
     * Force deletes a job by its ID
     * @param {number} jobId Job ID to delete
     * @returns {Object} Job object if deleted, otherwise null
     */
    forceDeleteById(jobId) {
        try {
            const stmt = database.prepare(`
        DELETE FROM ${table} WHERE id = ?
      `);

            const result = stmt.run(jobId);

            return result.changes > 0 ? this.findById(jobId) : null;
        } catch (error) {
            serverLog("Error force deleting job:", error.message);
            return null;
        }
    },

    /**
     * Paginates jobs from the table.
     * @param {string} [search=""] - Search term to filter jobs by name.
     * @param {number} [limit=10] - Number of jobs per page.
     * @param {number} [page=1] - Page number to fetch.
     * @returns {Array} - Array of paginated jobs.
     */
    paginate: (search = "", limit = 10, page = 1) => {
        const offset = (page - 1) * limit;
        let sql = `SELECT * FROM ${table} WHERE deleted_at IS NULL`;

        if (search) {
            sql += ` AND job_name LIKE '%${search}%' OR target_contact_or_group LIKE '%${search}%'`;
        }

        sql += " ORDER BY created_at DESC LIMIT ? OFFSET ?";
        try {
            const stmt = database.prepare(sql);
            return stmt.all(limit, offset);
        } catch (error) {
            serverLog("Error paginating jobs:", error.message);
            return null;
        }
    },

    /**
     * Updates a job by its ID
     * @param {number} jobId Job ID to update
     * @param {Object} jobData Job data to update
     * @returns {Object} Job object if updated, otherwise null
     */
    updateById(jobId, jobData) {
        try {
            const stmt = database.prepare(`
        UPDATE ${table} SET ${Object.keys(jobData)
    .map((key) => `${toSnakeCase(key)} = ?`)
    .join(", ")} WHERE id = ?
      `);

            const result = stmt.run(...Object.values(jobData), jobId);

            return result.changes > 0 ? this.findById(jobId) : null;
        } catch (error) {
            serverLog("Error updating job:", error.message);
            return null;
        }
    },

    /**
     * Counts all jobs in the table
     * @returns {number} Count of all jobs
     */
    countAll() {
        try {
            const stmt = database.prepare(
                `SELECT COUNT(*) AS total FROM ${table} WHERE deleted_at IS NULL`
            );
            return stmt.get().total;
        } catch (error) {
            serverLog("Error counting all jobs:", error.message);
            return 0;
        }
    },
};
