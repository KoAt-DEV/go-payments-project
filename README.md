# GoPayments API

[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://go.dev)
[![Fiber](https://img.shields.io/badge/Fiber-v3-00ADD8?logo=go)](https://gofiber.io)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A modern, production-oriented payment processing backend written in Go. This project is a portfolio/example app built with best practices in mind.

## Tech Stack

- Go 1.23+
- Fiber v3 (ultra-fast HTTP framework)
- Zerolog (fast structured logging)
- caarlos0/env (zero-boiler config from env vars)
- godotenv (.env file support)

## Contents

- [Overview](#overview)
- [Quick Start](#quick-start)
- [Requirements](#requirements)
- [Running Locally](#running-locally)
- [Configuration](#configuration)
- [Project Structure](#project-structure)
- [Features](#features)
- [Roadmap](#roadmap)
- [License](#license)

## Overview

This repository demonstrates a clean, layered Go backend using Fiber, structured logging, environment-driven configuration, and middleware for request tracing and logging.

## Quick Start

Clone the repository and prepare a local environment:

```bash
git clone https://github.com/KoAt-DEV/go-payments-portfolio-project.git
cd go-payments-portfolio-project
cp .env.example .env
# Edit .env to configure DB, Redis, JWT, PORT, etc.
```

Run the API for development:

```bash
# from project root, run the api command
go run ./cmd/api

# or build and run the binary
go build -o bin/gopayments ./cmd/api
./bin/gopayments
```

By default the server listens on the port specified in `PORT` (e.g. 3000). Check `GET /health` after startup.

## Requirements

- Go 1.23 or newer
- (Optional) PostgreSQL and Redis if you want the full stack locally

Run `go mod tidy` to fetch dependencies.

## Configuration

All configuration is read from environment variables (see `.env.example`). Main groups:

- `APP_*` — application settings
- `POSTGRES_*` — database connection
- `REDIS_*` — cache/session store
- `JWT_*` — authentication keys and durations
- `METRICS_*` — Prometheus/metrics settings

The code under `internal/config` contains the config definitions and parsing logic.

## Running with Docker

This project includes a `Dockerfile` and `docker-compose.yml` to run the API together with Postgres and Redis for local development.

Quick Docker steps:

```bash
# build and start api + postgres + redis
docker-compose up --build

# stop
docker-compose down
```

Notes about the compose setup:

- The API container is named `gpp-api` and is built from the repository `Dockerfile`.
- The compose file maps the API port `3000:3000` by default and uses the `.env` file for configuration. The service healthcheck calls `http://localhost:3000/health`.
- The `Dockerfile` uses a multi-stage build (golang:1.25-alpine builder → alpine runtime) and copies the `.env` into the image. The final image exposes port `8080` in the Dockerfile, but the application reads the port from the `APP_PORT` / `PORT` environment variable — prefer to set `APP_PORT` to match the port you want to expose (default 3000 in `.env.example`).

## Databases (Postgres & Redis)

The included `docker-compose.yml` defines Postgres and Redis services for local development. Key defaults (see `.env.example` and `docker-compose.yml`):

- Postgres (container name `gpp-postgres`)
	- Image: `postgres:16`
	- Default credentials (also shown in `.env.example`):
		- POSTGRES_USER=gppuser
		- POSTGRES_PASSWORD=gpppass
		- POSTGRES_DB=gppdb
	- Container internal port: `5432`
	- Host port mapped in the compose file: `5434:5432` (so from the host you can connect to Postgres on port 5434)
	- The application by default connects to the Postgres host `postgres` on port `5432` (container network). The `.env.example` uses `POSTGRES_HOST=postgres` and `POSTGRES_PORT=5432`.

- Redis (container name `gpp-redis`)
	- Image: `redis:7`
	- Default internal port: `6379`
	- Host port mapping: `6379:6379`
	- The `.env.example` uses `REDIS_ADDR=redis:6379` so the app connects to `redis` inside the compose network.

If you prefer to run Postgres/Redis locally (not in Docker), update the `.env` values accordingly and ensure the `POSTGRES_HOST` / `REDIS_ADDR` point to reachable host addresses.

## Project Structure

- `cmd/api` — application entrypoint
- `internal/config` — configuration management
- `internal/logger` — logging setup (Zerolog wrapper)
- `internal/middleware` — request ID, request logger, etc.

This structure follows Go conventions for internal packages and a single `cmd` binary.

## Features

- Layered, extendable architecture
- Environment-based configuration with `.env` support
- Structured logging via Zerolog (pretty console in dev, JSON in prod)
- Unique request IDs per request
- Request logging middleware with slow-request warnings
- Prometheus-compatible metrics namespace (`gopayments_api`) prepared
- JWT authentication scaffolding
- Production-grade graceful shutdown with context-aware resource cleanup (Postgres, Redis, Fiber)
- Zero connection leaks even under Kubernetes rolling updates
- Health-check endpoint

## Roadmap

- User registration & login (JWT)
- Refresh token rotation
- Role-based access control
- Payment processing endpoints (Stripe/PayPal/webhooks)
- Rate limiting
- Prometheus metrics endpoint
- OpenAPI/Swagger documentation
- Docker + docker-compose
- CI/CD via GitHub Actions
- Unit & integration tests

## License

This project is licensed under the MIT License — see the `LICENSE` file.

---

Made with ❤️ by @KoAt-DEV



