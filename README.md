# Arrflix

Arrflix is an **early-stage, self-hosted media management project**.

It’s an experiment in simplifying personal media automation — with a focus on being easier to understand, easier to reason about, and less fragile than many existing setups.

This project is **actively evolving** and not yet production-ready.

## Project Status

Arrflix is:

- Experimental
- Opinionated
- Incomplete
- Subject to breaking changes

Documentation and features will change as the project evolves.  
If you’re looking for something stable and polished today, Arrflix probably isn’t there yet.

If you’re comfortable experimenting or following along with an evolving project, welcome.

## Quick Start

If you’re interested in trying Arrflix, start with the documentation:

👉 **Introduction & Overview**  
https://kyleaupton.github.io/Arrflix/guide/introduction.html

From there, you can continue to the **Getting Started** guide for installation instructions.

👉 **Getting Started / Installation**  
https://kyleaupton.github.io/Arrflix/guide/getting-started.html

## Development Setup

If you’re here to hack on Arrflix, the dev setup is lightweight.

### Requirements

- Docker
- Docker Compose
- A TMDB API key

### Local Development

1. Clone the repository:

   ```bash
   git clone https://github.com/kyleaupton/Arrflix.git
   cd Arrflix
   ```

2. Create a `.env` file:

   ```env
   TMDB_API_KEY=your_tmdb_api_key_here
   MEDIA_LIBRARIES=/path/to/test/media
   ```

3. Start the development stack:
   ```bash
   docker compose up
   ```

That’s it. The backend, frontend, database, and supporting services run together via Docker Compose.

## Documentation

All project documentation lives here:

👉 https://kyleaupton.github.io/Arrflix/

Expect documentation to lag behind implementation at times — this is normal for the current stage of the project.

## Contributing & Feedback

There’s no strict roadmap yet. The project is still finding its shape.

Feedback and discussions are encouraged, but pull requests may be declined until the project’s core design has stabilized (especially larger ones).

## License

License not finalized yet.
