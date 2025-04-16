# Rust Full-Stack Application

A modern full-stack web application built with Rust, using:

- Axum for the web framework
- HTMX for dynamic UI updates
- PostgreSQL for the database
- Redis for caching
- Tera for templating
- Tailwind CSS for styling

## Prerequisites

- Rust (latest stable version)
- PostgreSQL
- Redis
- Docker (optional, for running dependencies)

## Setup

1. Install dependencies:

   ```bash
   cargo build
   ```

2. Set up the database:

   ```bash
   createdb gitlab_redux
   ```

3. Configure environment variables:
   - Copy `.env.example` to `.env`
   - Update the database and Redis URLs if needed

## Running the Application

1. Start the development server:

   ```bash
   cargo run
   ```

2. Access the application at `http://localhost:3000`

## Development

- The application uses hot-reloading for templates
- HTMX is used for dynamic UI updates without writing JavaScript
- PostgreSQL is used for persistent storage
- Redis is used for caching and session management

## Project Structure

```
src/
├── handlers/     # Route handlers
├── models/       # Data models
├── services/     # Business logic
├── templates/    # Tera templates
└── utils/        # Utility functions
static/
└── assets/       # Static assets (CSS, JS, images)
    ├── stylesheets/
    ├── javascripts/
    └── images/
```

## Static Assets

The application serves static assets from the `static/assets` directory:

- CSS files are in `static/assets/stylesheets/`
- JavaScript files are in `static/assets/javascripts/`
- Images are in `static/assets/images/`

These assets are served at the `/assets` URL path.

## API Endpoints

- `GET /` - Home page
- `GET /health` - Health check endpoint

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request
