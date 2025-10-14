# Snaggle Docs

## Installation

Target `docker-compose.yml` file (for now):

```yml
services:
  snaggle:
    image: ghcr.io/kyleaupton/snaggle:0.0.1
    environment:
      TMDB_API_KEY: ${TMDB_API_KEY}
    ports:
      - "8484:8484"
    volumes:
      - snaggle_pg_data:/var/lib/postgresql/data
    restart: unless-stopped

volumes:
  snaggle_pg_data:
```

## Local development

Run a single dev container with supervisord orchestrating Postgres, API (live reload), Vite dev server, and Nginx fronting on port 8484.

### Start (dev)

```bash
docker compose -f ops/docker-compose.yml -f ops/docker-compose.dev.yml up --build
```

- App: `http://localhost:8484` (proxied through Nginx)
- API: `http://localhost:8484/api/`

### Environment

Set at least `TMDB_API_KEY` in your shell or an `.env` file. Optional: `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `TZ`.

### Notes

- First run will install `web/node_modules` inside the container; it's persisted via the `web_node_modules` volume.
- Hot reload: Vite serves the SPA; Nginx proxies `/` to Vite and `/api/` to the Go API. No CORS needed.
- Database data is persisted in the `snaggle_pg_data` volume.
