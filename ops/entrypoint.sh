# --- Build API (Go) ---
FROM golang:1.22 AS api-build
WORKDIR /app
COPY backend/go.mod backend/go.sum ./backend/
RUN cd backend && go mod download
COPY backend ./backend
RUN cd backend && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/api ./cmd/api

# --- Build Web (Vite SPA) ---
FROM node:20 AS web-build
WORKDIR /web
COPY web/package*.json ./
RUN npm ci
COPY web/ .
RUN npm run build

# --- Runtime: Postgres + Nginx + supervisord ---
# Using postgres:16 so we can leverage its robust init flow
FROM postgres:16-bookworm

# Core tools
RUN apt-get update && apt-get install -y --no-install-recommends \
      nginx supervisor dumb-init curl ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Nginx (serve SPA, proxy /api -> 127.0.0.1:8080)
COPY --from=web-build /web/dist /usr/share/nginx/html
COPY ops/nginx/default.conf /etc/nginx/conf.d/default.conf

# API binary
COPY --from=api-build /out/api /usr/local/bin/api

# Supervisor config + helper
COPY ops/supervisor/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY ops/scripts/wait-for-db.sh /usr/local/bin/wait-for-db.sh
RUN chmod +x /usr/local/bin/wait-for-db.sh

# Expose the “signature” port for Nginx and the local Postgres port (optional)
EXPOSE 8484 5432

# Volumes:
# - PGDATA already defined in base image (/var/lib/postgresql/data)
# - Optional: a place for app logs if you want
VOLUME ["/var/lib/postgresql/data"]

# Entrypoint: supervisord under dumb-init (PID1)
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["supervisord", "-n", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
