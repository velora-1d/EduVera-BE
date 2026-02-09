> [!CAUTION]
> **Whatsback Web** and [whatsapp-web.js](https://github.com/pedroslopez/whatsapp-web.js) are not officially supported by WhatsApp. Use this project at your own risk.

<p align="center">
    <img src="https://i.imgur.com/X1bM7I6.png" />
</p>
<p align="center">
  <img src="https://github.com/darkterminal/whatsback-web/actions/workflows/release.yml/badge.svg" alt="Whatsback Web GitHub Action" />
  <img src="https://img.shields.io/github/tag/darkterminal/whatsback-web" alt="Whatsback Web Tag" />
  <img src="https://img.shields.io/github/v/release/darkterminal/whatsback-web" alt="Whatsback Web Release" />
  <img src="https://img.shields.io/github/v/tag/darkterminal/whatsback-web?label=package" alt="Whatsback Web Package Registry" />
  <img src="https://ghcr-badge.egpl.dev/darkterminal/whatsback-web/size?color=%2344cc11&tag=latest&label=image+size" alt="Whatsback Web Image Size" />
  <a href="https://discord.gg/ZQPEtcWyBh">
    <img alt="Discord" src="https://img.shields.io/discord/1343952367615217764?style=flat&logo=discord&label=discord">
  </a>
</p>

Whatsback Provider is a simple WhatsApp provider that offers basic functionality such as predefined static commands, sending messages to contacts or groups, and listing all contacts. This project leverages the unofficial [whatsapp-web.js](https://github.com/pedroslopez/whatsapp-web.js) package to interface with WhatsApp Web.

> [!IMPORTANT]
> **It is not guaranteed you will not be blocked by using this method. WhatsApp does not allow bots or unofficial clients on their platform, so this shouldn't be considered totally safe.**

## Features

- **Predefined Static Commands:**  Quickly execute common commands without the need for manual input.
- **Send Message to Contact:**  Programmatically send message directly to individual contacts.
- **Send Message to Group:**  Programmatically send message to the groups.
- **List All Contacts:**  Retrieve and display a list of all contacts available on the WhatsApp account.
- **Schedule Message:**  Schedule a message to be sent at a later time.

---


<h4 align="center">Leave Your Reviews & Start ⭐ This Repository</h4>
<p align="center">
    <a href="https://www.producthunt.com/posts/whatsback-web?embed=true&utm_source=badge-featured&utm_medium=badge&utm_souce=badge-whatsback&#0045;web" target="_blank"><img src="https://api.producthunt.com/widgets/embed-image/v1/featured.svg?post_id=917168&theme=light&t=1740688609385" alt="Whatsback&#0032;Web - Host&#0032;Your&#0032;Own&#0032;Whatsapp&#0032;Provider&#0032;&#0045;&#0032;Free&#0032;and&#0032;Open&#0032;Source | Product Hunt" style="width: 250px; height: 54px;" width="250" height="54" /></a>
</p>

---

## Table of Contents

- [⚠️ Caution](#caution)
- [✨ Features](#features)
- [🚀 Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Key Configuration Options](#key-configuration-options)
- [🐳 Docker Installation](#docker-installation)
  - [Docker CLI](#docker-cli)
  - [docker-compose.yml](#docker-composeyml)
- [💻 Source Installation](#source-installation)
- [🏃 Running the Project](#running-the-project)
- [🔌 Available REST API](#available-rest-api)
- [🔒 Security Considerations](#security-considerations)
- [🤝 Contributing](#contributing)
- [💖 Donate or Sponsoring](#donate-or-sponsoring)
- [📢 Disclaimer](#disclaimer)
- [📄 License](#license)

## Getting Started

### Prerequisites

- Node.js (v20 or later recommended)
- npm (comes with Node.js)
- A valid WhatsApp account

### Key Configuration Options

| Environment Variable | Description              | Default                  |
|----------------------|--------------------------|------------------------- |
| `NODE_ENV`           | Runtime environment       | `production`            |
| `APP_PORT`           | Application port          | `5001`                  |
| `UI_PORT`            | External exposed port     | `8169`                  |
| `DB_PATH`            | Path to SQLite database   | `/data/database.sqlite` |
| `TZ`                 | Set your default timezone | `Asia/Jakarta` default `UTC` |

### Docker Installation

#### Docker CLI

1. **Pull the image:**

   ```bash
   docker pull ghcr.io/darkterminal/whatsback-web:latest
   ```

2. **Create Network**

  ```bash
  docker network create whatsback-net
  ```

3. **Create Volume**

  ```bash
  docker volume create whatsback-db
  ```

4. **Run Whatsback Application Container**

  ```bash
  docker run -d \
    --name whatsback-app-provider \
    --network whatsback-net \
    -p 8169:5001 \
    -e NODE_ENV=production \
    -e APP_PORT=5001 \
    -e DB_PATH=/data/database.sqlite \
    -v whatsback-db:/data \
    ghcr.io/darkterminal/whatsback-web:latest
  ```

5. **Run Whatsback Cronjob Container**

  ```bash
  docker run -d \
    --name whatsback-app-cronjob \
    --network whatsback-net \
    -e NODE_ENV=production \
    -e APP_HOST=whatsback-app-provider \
    -e DB_PATH=/data/database.sqlite \
    -v whatsback-db:/data \
    ghcr.io/darkterminal/whatsback-web:latest \
    sh -c "./wait-for whatsback-app-provider:5001 -t 120 -- node cronjob.js"
  ```

#### docker-compose.yml

```yaml
services:
  app:
    image: ghcr.io/darkterminal/whatsback-web:latest
    container_name: whatsback-app-provider
    ports:
      - "${UI_PORT:-8169}:5001"
    environment:
      - NODE_ENV=production
      - APP_PORT=${APP_PORT:-5001}
      - DB_PATH=/data/database.sqlite
    volumes:
      - db-data:/data
    networks:
      - app_net
    healthcheck:
      test: ["CMD-SHELL", "curl -f http://localhost:$$APP_PORT/health || exit 1"]
      interval: 15s
      timeout: 10s
      retries: 5

  cronjob:
    image: ghcr.io/darkterminal/whatsback-web:latest
    container_name: whatsback-app-cronjob
    environment:
      - NODE_ENV=production
      - APP_HOST=app
      - DB_PATH=/data/database.sqlite
    volumes:
      - db-data:/data
    command: sh -c "./wait-for app:5001 -t 120 -- node cronjob.js"
    networks:
      - app_net
    depends_on:
      app:
        condition: service_healthy

networks:
  app_net:
    driver: bridge

volumes:
  db-data:
```

### Source Installation

```bash
git clone https://github.com/darkterminal/whatsback-web.git
cd whatsback-web
npm install
```
   
### Running the Project

For development, start the server with:

```bash
npm run dev
```

For production, start the server with:

```bash
NODE_ENV=production node server.js
```

Your server should start on the port defined in the `.env` file (default is 5001).

## Available REST API

You can read the REST API documentation [here](WHATSBACK-API.md)

## Security Considerations

- This project uses middleware like Helmet, express-rate-limit, and hpp to help protect against common web vulnerabilities.
- Be aware that using an unofficial API (whatsapp-web.js) can carry risks with regard to WhatsApp's terms of service.
- Whatsback Web and whatsapp-web.js are not officially supported by WhatsApp. Use this project at your own risk.
- It is not guaranteed you will not be blocked by using this method. WhatsApp does not allow bots or unofficial clients on their platform, so this shouldn't be considered totally safe.

## Contributing

Contributions are welcome! Please open issues or pull requests to improve the project.

## Donate or Sponsoring

You can support the maintainer of this project through the button and links below:

<a href="https://github.com/sponsors/darkterminal">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://i.imgur.com/IcHNW1L.png">
    <source media="(prefers-color-scheme: light)" srcset="https://i.imgur.com/Yzbwovb.png">
    <img alt="Shows a black logo in light color mode and a white one in dark color mode." src="https://i.imgur.com/IcHNW1L.png">
  </picture>
</a>

- [Support via GitHub Sponsors](https://github.com/sponsors/darkterminal)
- [Support via Saweria](https://saweria.co/darkterminal)

## Disclaimer

This project is not affiliated, associated, authorized, endorsed by, or in any way officially connected with WhatsApp or any of its subsidiaries or its affiliates. The official WhatsApp website can be found at [whatsapp.com](https://whatsapp.com). "WhatsApp" as well as related names, marks, emblems and images are registered trademarks of their respective owners. Also it is not guaranteed you will not be blocked by using this method. WhatsApp does not allow bots or unofficial clients on their platform, so this shouldn't be considered totally safe.

## License

```
Copyright 2025 Imam Ali Mustofa

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
<!-- GitAds-Verify: EW24FYC84IIZIJ3DV3UFDHDWUFLF7YTQ -->

## GitAds Sponsored
[![Sponsored by GitAds](https://gitads.dev/v1/ad-serve?source=tursodatabase/turso-driver-laravel@github)](https://gitads.dev/v1/ad-track?source=tursodatabase/turso-driver-laravel@github)
