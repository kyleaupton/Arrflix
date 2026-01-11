# Getting Started

This guide will help you get Arrflix up and running on your system.

## Prerequisites

Before you begin, make sure you have:

- **Docker** and **Docker Compose** installed
- A **TMDB API key** (see below)
- Access to your media library directories

## Acquiring a TMDB API Key

Arrflix requires a TMDB (The Movie Database) API key to fetch metadata for movies and TV shows.

1. Visit [TMDB](https://www.themoviedb.org/) and create a free account
2. Go to your [Account Settings](https://www.themoviedb.org/settings/api)
3. Request an API key under the "API" section
4. Copy your API key - you'll need it for the installation step

## Installation

Create a `docker-compose.yml` file with the following configuration:

```yml
services:
  arrflix:
    image: ghcr.io/kyleaupton/arrflix:latest
    environment:
      TMDB_API_KEY: your_api_key_here
    ports:
      - "8484:8484"
    volumes:
      - /path/to/libraries/:/data
      - arrflix_pg_data:/var/lib/postgresql/data
    restart: unless-stopped

volumes:
  arrflix_pg_data:
```

### Configuration Steps

1. **Replace `your_api_key_here`** with your actual TMDB API key in the `docker-compose.yml` file.

2. **Update the volume path** - Replace `/path/to/libraries/` with the actual path to your media library directory.

3. **Start Arrflix**:

   ```bash
   docker compose up -d
   ```

4. **Access the web interface** at `http://localhost:8484`

## Next Steps

Once Arrflix is running, you can:

- Configure your media libraries
- Add indexers and download clients
- Start downloading your favorite shows and movies
