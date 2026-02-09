FROM ghcr.io/puppeteer/puppeteer:24.7.0

# Skip Chromium download - already included in base image
ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true \
    PUPPETEER_EXECUTABLE_PATH=/usr/bin/google-chrome-stable \
    NODE_ENV=production

USER root

# Install additional utilities
RUN apt-get update && apt-get install -y --no-install-recommends \
    netcat-openbsd \
    dumb-init \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /usr/src/app

# Copy package files first for better caching
COPY package*.json ./
RUN npm ci --omit=dev && npm cache clean --force

# Copy source code
COPY . .

RUN chmod +x wait-for

# Create data directory for whatsapp-web.js auth
RUN mkdir -p /usr/src/app/.wwebjs_auth && \
    chmod -R 777 /usr/src/app/.wwebjs_auth

EXPOSE 5001

CMD ["dumb-init", "node", "server.js"]
