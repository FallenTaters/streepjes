# Streepjes

Streepjes is a custom POS app for two specific sports clubs and a hobby project written in pure Go.

The frontend uses server-side rendered HTML templates with vanilla JavaScript for interactivity. The backend uses PostgreSQL.

## Development

### Requirements

* go >= 1.26
* [just](https://github.com/casey/just) (command runner)
* [entr](https://eradman.com/entrproject/) (for hot-reloading during development)
* go tooling:
    * [enumer](https://github.com/dmarkham/enumer) (installed as a tool dependency via `go.mod`)

### Run locally

1. Copy `.env.dev.example` to `.env.dev` and configure your PostgreSQL connection string
2. `just run` (watches for file changes and auto-restarts)

### Build

`just build`

This produces `./bin/streepjes`. CGO is not required.

### Build container

`just container`

## Settings

Configuration via environment variables (or `.env.dev` for local development):

```
STREEPJES_PORT=80
STREEPJES_DB_CONNECTION_STRING=postgres://...
STREEPJES_DISABLE_SECURE=false
```
