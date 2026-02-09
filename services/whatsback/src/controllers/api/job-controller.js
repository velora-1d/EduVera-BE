const { serverLog } = require("../../helper");
const job = require("../../models/job");

module.exports = {
    /**
     * GET /api/jobs
     * @summary Fetch a page of jobs (with or without a search term).
     * @param {import("express").Request} req - The request object, containing query parameters for search, perPage, and page.
     * @param {import("express").Response} res - The response object.
     * @returns {void} - Responds with a JSON object containing a success status and an array of jobs if successful, otherwise an error message.
     * @throws {Error} - If there is an error with the database query.
     */
    getJobs: (req, res) => {
        try {
            const searchTerm = req.query.search || "";
            const perPage = req.query.perPage || 10;
            const page = req.query.page || 1;

            const jobsToDisplay = job.paginate(searchTerm, perPage, page);

            res.status(200).json({
                success: true,
                data: jobsToDisplay,
                total: job.countAll(),
            });
        } catch (error) {
            res.status(500).json({
                status: false,
                message: `API Error - ${error}`,
            });
        }
    },

    /**
     * Updates an existing job by its ID.
     *
     * @param {import("express").Request} req - The request object, containing the job ID in params and the job data in body.
     * @param {import("express").Response} res - The response object.
     * @returns {void} - Responds with a JSON object containing success status and updated job data if successful, otherwise an error message.
     * Logs the update operation and any errors encountered.
     */
    updateJob: (req, res) => {
        const jobId = req.params.id;
        const jobData = req.body;
        try {
            const store = job.updateById(jobId, jobData);
            serverLog(`Job with ID ${jobId} updated`);
            return res.status(200).json({
                success: true,
                data: store,
            });
        } catch (error) {
            serverLog(`Error: ${error}`);
            return res.status(500).json({
                status: false,
                message: `API Error - ${error}`,
            });
        }
    },

    /**
     * Soft deletes a job by its ID.
     *
     * @param {import("express").Request} req - The request object, containing the job ID in params.
     * @param {import("express").Response} res - The response object.
     * @returns {void} - Responds with a JSON object containing success status and deleted job data if successful, otherwise an error message.
     * Logs the delete operation and any errors encountered.
     */
    deleteJob: (req, res) => {
        const jobId = req.params.id;
        try {
            const store = job.softDeleteById(jobId);
            serverLog(`Job with ID ${jobId} deleted`);
            return res.status(200).json({
                success: true,
                data: store,
            });
        } catch (error) {
            serverLog(`Error: ${error}`);
            return res.status(500).json({
                status: false,
                message: `API Error - ${error}`,
            });
        }
    },

    /**
     * Force deletes a job by its ID.
     *
     * @param {import("express").Request} req - The request object, containing the job ID in params.
     * @param {import("express").Response} res - The response object.
     * @returns {void} - Responds with a JSON object containing success status and deleted job data if successful, otherwise an error message.
     * Logs the delete operation and any errors encountered.
     */
    forceDeleteJob: (req, res) => {
        const jobId = req.params.id;
        try {
            const store = job.forceDeleteById(jobId);
            serverLog(`Job with ID ${jobId} force deleted`);
            return res.status(200).json({
                success: true,
                data: store,
            });
        } catch (error) {
            serverLog(`Error: ${error}`);
            return res.status(500).json({
                status: false,
                message: `API Error - ${error}`,
            });
        }
    },

    /**
     * Creates a new job in the database.
     *
     * @param {Object} jobData - The job data.
     * @param {string} jobData.job_name - The name of the job.
     * @param {string} jobData.job_trigger - The trigger type for the job (e.g., "send_message").
     * @param {string} jobData.target_contact_or_group - The target contact or group for the job.
     * @param {string} jobData.message - The message content for the job.
     * @param {string} jobData.job_cron_expression - The cron expression for scheduling the job.
     * @returns {boolean} - Returns true if the job was successfully created, otherwise false.
     */
    createJob: (req, res) => {
        const jobData = req.body;
        try {
            const store = job.create(
                jobData.job_name,
                jobData.job_trigger,
                jobData.target_contact_or_group,
                jobData.message,
                jobData.job_cron_expression
            );

            serverLog(`New job created with name ${jobData.job_name}`);

            return res.status(200).json({
                success: true,
                data: store,
            });
        } catch (error) {
            serverLog(`Error: ${error}`);
            return res.status(500).json({
                status: false,
                message: `API Error - ${error}`,
            });
        }
    },

    /**
     * Finds a job by its ID.
     *
     * @param {number} id - The job ID.
     * @returns {(Object|null)} - The job object if found, otherwise null.
     */
    findJobById: (req, res) => {
        const id = req.params.id;
        try {
            const jobData = job.findById(id);
            if (jobData) {
                serverLog(`Job with ID ${id} found`);
                return res.status(200).json({
                    success: true,
                    data: jobData,
                });
            }
            serverLog(`Job with ID ${id} not found`);
            return res.status(404).json({
                success: false,
                message: `Job with ID ${id} not found`,
            });
        } catch (error) {
            serverLog(`Error finding job with ID ${id}: ${error}`);
            return res.status(500).json({
                status: false,
                message: `API Error - ${error}`,
            });
        }
    },

    /**
     * Finds jobs by their status.
     *
     * @param {number} status - The status of the jobs to be found (e.g., 0 for disabled, 1 for enabled).
     * @returns {(Array|null)} - An array of job objects if found, otherwise null.
     */
    findJobByStatus: (req, res) => {
        const status = req.params.status;
        try {
            const jobData = job.findByStatus(status);
            if (jobData) {
                serverLog(`Job with status ${status} found`);
                return res.status(200).json({
                    success: true,
                    data: jobData,
                });
            }
            serverLog(`Job with status ${status} not found`);
            return res.status(404).json({
                success: false,
                message: `Job with status ${status} not found`,
            });
        } catch (error) {
            serverLog(`Error finding job with status ${status}: ${error}`);
            return res.status(500).json({
                status: false,
                message: `API Error - ${error}`,
            });
        }
    },
};
