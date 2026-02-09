const express = require("express");
const helmet = require("helmet");
const hpp = require("hpp");
const { createServer } = require("node:http");
const { Server } = require("socket.io");
const cors = require("cors");
require("dotenv").config();

const { client, setSocketManager } = require("./src/whatsapp-client");
const { parseOrigins, bannerCLI } = require("./src/helper");
const path = require("node:path");
const expressLayouts = require("express-ejs-layouts");
const schemaMigrations = require("./src/schemas");
const socketInstance = require("./lib/socket-instance");

// ===================================================
// API Routes
// ===================================================
const commandRoutes = require("./src/routes/api/command-routes");
const contactRoutes = require("./src/routes/api/contact-routes");
const cronRoutes = require("./src/routes/api/cron-routes");
const groupRoutes = require("./src/routes/api/group-routes");
const messageRoutes = require("./src/routes/api/message-routes");
const jobRoutes = require("./src/routes/api/job-routes");
const sessionRoutes = require("./src/routes/api/session-routes");

// ===================================================
// Frontend Routes
// ===================================================
const indexFrontRoutes = require("./src/routes/index-front-routes");
const commandFrontRoutes = require("./src/routes/command-front-routes");
const messageFrontRoutes = require("./src/routes/message-front-routes");
const contactFrontRoutes = require("./src/routes/contact-front-routes");

const helmetPolicies = require("./src/middlewares/helmet-policies");
const apiRateLimiter = require("./src/middlewares/api-rate-limiter");
const { additionalCors, apiCors } = require("./src/middlewares/cors-options");
const appLayout = require("./src/middlewares/app-layout");
const socketOptions = require("./src/socket-options");

schemaMigrations();

const app = express();

// ===================================================
// Middlewares
// ===================================================
app.use("/api", additionalCors);
app.use(cors(apiCors));
app.use(helmet.contentSecurityPolicy(helmetPolicies));
app.disable("x-powered-by");
app.use(hpp());
app.use("/api", apiRateLimiter);

app.use(express.static(path.join(__dirname, "public")));
app.set("views", path.join(__dirname, "src", "views"));
app.use(expressLayouts);
app.set("layout extractScripts", true);
app.set("view engine", "ejs");
app.use(appLayout);

app.use(express.json());

// ===================================================
// Frontend routes
// ===================================================
app.use("/", indexFrontRoutes);
app.use("/commands", commandFrontRoutes);
app.use("/contacts", contactFrontRoutes);
app.use("/message", messageFrontRoutes);
app.get("/profile", (_, res) =>
    res.render("profile", { title: "Profile", pathname: "profile" })
);

// ===================================================
// API routes
// ===================================================
app.use("/api/command", commandRoutes);
app.use("/api/contacts", contactRoutes);
app.use("/api/cron-next-runs", cronRoutes);
app.get("/health", (_, res) => res.status(200).send("OK"));
app.use("/api/groups", groupRoutes);
app.use("/api/message", messageRoutes);
app.use("/api/jobs", jobRoutes);
app.use("/api/session", sessionRoutes);

// ===================================================
// Error Handlers
// ===================================================
app.use((_, res) => {
    res.status(404).render("404", { title: "404 Not Found" });
});
/* eslint-disable */
app.use((error, _, res, __) => {
    console.error('Unhandled error:', error);
    res.status(500).json({ message: 'Internal Server Error' });
});
/* eslint-enable */

// ===================================================
// HTTP & Socket Server Instance
// ===================================================
const httpServer = createServer(app);
const io = new Server(httpServer, {
    cors: {
        origin: parseOrigins(process.env.SOCKET_IO_CORS_ORIGIN),
        methods: ["GET", "POST"],
        credentials: true,
    },
});
socketInstance.setIo(io);

const connectedSockets = [];
io.on("connection", async (socket) =>
    socketOptions(socket, connectedSockets, client)
);
setSocketManager(connectedSockets);

const PORT = process.env.APP_PORT || 5001;
httpServer.listen(PORT, "0.0.0.0", bannerCLI(PORT));
