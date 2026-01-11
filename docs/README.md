# Arrflix Docs

## Installation

Here's what a docker-compose.yml file looks like:

```yml
services:
  arrflix:
    image: ghcr.io/kyleaupton/arrflix:0.0.1
    environment:
      TMDB_API_KEY: ${TMDB_API_KEY}
    ports:
      - "8484:8484"
    volumes:
      - /path/to/libraries/:/data
      - arrflix_pg_data:/var/lib/postgresql/data
    restart: unless-stopped

volumes:
  arrflix_pg_data:
```
