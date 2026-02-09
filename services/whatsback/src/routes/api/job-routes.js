const express = require("express");
const {
    getJobs,
    createJob,
    findJobById,
    findJobByStatus,
    updateJob,
    deleteJob,
    forceDeleteJob,
} = require("../../controllers/api/job-controller");
const router = express.Router();

// GET /api/jobs
router.get("/", getJobs);
// GET /api/jobs/:id
router.get("/:id", findJobById);
// GET /api/jobs/:status
router.get("/status/:status", findJobByStatus);
// POST /api/jobs
router.post("/", createJob);
// PUT /api/jobs/:id
router.put("/:id", updateJob);
// DELETE /api/jobs/:id
router.delete("/:id", deleteJob);
// DELETE /api/jobs/force/:id
router.delete("/force/:id", forceDeleteJob);

module.exports = router;
