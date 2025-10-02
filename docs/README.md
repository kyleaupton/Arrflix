# Snaggle Docs

## Installation

Target `docker-compose.yml` file (for now):

```yml
services:
  media:
    image: ghcr.io/kyleaupton/snaggle:0.0.1
    ports:
      - "8484:8484"
    volumes:
      - media_data:/data # will def change
    restart: unless-stopped
```
