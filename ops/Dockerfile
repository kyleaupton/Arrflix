# --- Build API ---
FROM golang:1.24 AS api-build
WORKDIR /app
COPY backend/go.mod backend/go.sum ./backend/
RUN cd backend && go mod download
COPY backend ./backend
RUN cd backend && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/api ./cmd/api

# --- Build Web (Vue SPA) ---
FROM node:20 AS web-build
WORKDIR /web
COPY web/package*.json ./
RUN npm ci
COPY web/ .
RUN npm run build

# --- Runtime (Nginx + API + Postgres + Supervisor) ---
FROM debian:stable-slim

# install deps: nginx, curl, dumb-init, postgres, and supervisor
RUN apt-get update && \
    apt-get install -y --no-install-recommends nginx postgresql postgresql-contrib ca-certificates dumb-init curl supervisor && \
    rm -rf /var/lib/apt/lists/*

# Make Postgres client/server binaries easy to find on PATH
# Debian installs them under /usr/lib/postgresql/<major>/bin
RUN set -eu; \
    pg_bin_base="/usr/lib/postgresql"; \
    if [ -d "$pg_bin_base" ]; then \
      for vdir in "$pg_bin_base"/*; do \
        if [ -d "$vdir/bin" ]; then \
          ln -sf "$vdir/bin/initdb" /usr/local/bin/initdb || true; \
          ln -sf "$vdir/bin/pg_ctl" /usr/local/bin/pg_ctl || true; \
          ln -sf "$vdir/bin/postgres" /usr/local/bin/postgres || true; \
          ln -sf "$vdir/bin/psql" /usr/local/bin/psql || true; \
        fi; \
      done; \
    fi

# Copy built artifacts
COPY --from=web-build /web/dist /usr/share/nginx/html
COPY ops/nginx/default.conf /etc/nginx/conf.d/default.conf
COPY --from=api-build /out/api /usr/local/bin/api
COPY ops/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Supervisor + helper scripts
COPY ops/supervisor/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY ops/scripts/wait-for-db.sh /usr/local/bin/wait-for-db.sh
COPY ops/scripts/run-postgres.sh /usr/local/bin/run-postgres.sh
RUN chmod +x /usr/local/bin/wait-for-db.sh /usr/local/bin/run-postgres.sh

# Prepare Postgres data directory
RUN mkdir -p /var/lib/postgresql/data && chown -R postgres:postgres /var/lib/postgresql
VOLUME ["/var/lib/postgresql/data", "/data"]

EXPOSE 8484

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/entrypoint.sh"]
