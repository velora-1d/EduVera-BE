# Walkthrough - EduVera Backend Migration Deployment

## WhatsMeow Migration Status
**Completed Successfully** on 2026-02-10

### 1. Key Changes
- **Evolution API Removed**: The external Evolution API service has been completely removed from the architecture.
- **WhatsMeow Integrated**: WhatsApp functionality is now embedded directly in the Go backend using `whatsmeow` library.
- **SQLite Persistence**: WhatsApp sessions are stored in `wa_sessions` volume on the VPS.
- **CGO Enabled**: Dockerfile updated to support SQLite/CGO builds.

### 2. Verification Steps

#### Backend Health Check
```bash
curl https://api-eduvera.ve-lora.my.id/
# Response: {"app":"EduVera API","status":"running","version":"1.0.0"}
```

#### WhatsApp Re-pairing (REQUIRED)
Since the underlying engine changed, **all previous WhatsApp sessions are invalid**. You must re-pair them.

1.  **Login to Owner Dashboard**: [https://eduvera.ve-lora.my.id/](https://eduvera.ve-lora.my.id/)
2.  Navigate to **WhatsApp Integration** menu.
3.  Click **Connect** (or disconnect first if it shows old status).
4.  **Scan the QR Code** with your real device.
5.  Wait for status to change to **Connected**.

### 3. Troubleshooting
If pairing fails:
- Check backend logs: `ssh ubuntu@43.156.132.218 "docker logs -f --tail 100 eduvera_backend"`
- Look for `[WHATSMEOW]` logs.
- Ensure 8000 port is accessible (it is via Nginx reverse proxy).

### 4. Comparison
| Feature | Old (Evolution) | New (WhatsMeow) |
| :--- | :--- | :--- |
| **Architecture** | External HTTP Service | Native Go Library |
| **Memory Usage** | High (Node.js + Chrome) | Low (Go + SQLite) |
| **Latency** | Network overhead | Zero latency |
| **Reliability** | Dependent on external container | Integrated with app lifecycle |
