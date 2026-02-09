module.exports = {
    directives: {
        defaultSrc: ["'self'"],
        blockAllMixedContent: [],
        frameAncestors: ["'self'"],
        scriptSrc: [
            "'self'",
            "https://cdn.tailwindcss.com",
            "https://cdnjs.cloudflare.com",
            "https://cdn.jsdelivr.net",
            "'unsafe-inline'",
        ],
        styleSrc: [
            "'self'",
            "https://cdn.tailwindcss.com",
            "https://cdnjs.cloudflare.com",
            "https://fonts.googleapis.com",
            "'unsafe-inline'",
        ],
        imgSrc: [
            "'self'",
            "data:",
            "https://robohash.org",
            "https://pps.whatsapp.net",
        ],
        connectSrc: ["'self'"],
        fontSrc: [
            "'self'",
            "https://cdnjs.cloudflare.com",
            "https://fonts.gstatic.com",
        ],
    },
};
