# Streepjes

Custom POS (Point of Sale) app for sports clubs. Tracks orders placed at club bars, member debt, monthly billing, and leaderboards. Written in pure Go with server-side rendered HTML.

## Architecture

```
domain/                  Pure business entities (no dependencies)
├── club.go              Club enum (Parabool, Gladiators, Calamari)
├── authdomain/          User, Role, Permission
└── orderdomain/         Order, Member, Item, Category, Price, Status

api/                     Shared payload types (used by handlers and templates)

backend/
├── application/         Business logic services
│   ├── auth/            auth.Service — login, sessions, user CRUD
│   └── order/           order.Service — orders, members, catalog, billing, leaderboard
├── global/settings/     Config struct and defaults
└── infrastructure/
    ├── repo/            Repository interfaces (Order, Member, User, Catalog)
    │   ├── postgres/    PostgreSQL implementations + migrations
    │   └── mockdb/      Mock implementations for tests
    └── router/          HTTP handlers and routing (chi)

templates/               Go html/template files (embedded via embed.FS)
├── base.html            Base layout
├── nav.html             Navigation component
├── *.html               Page templates
└── admin/*.html         Admin page templates

static/files/            Static assets (embedded), BeerCSS framework
```

**Dependency flow:** `domain` ← `application` ← `infrastructure` ← `main.go`

- Domain layer has zero imports from other project packages
- Application services depend on repository _interfaces_, not implementations
- Infrastructure implements interfaces and wires everything together

## Domain Model

- **Club** — top-level tenant (`Parabool`, `Gladiators`, `Calamari`). Separates users, members, orders. Items have per-club pricing.
- **User** — staff operating the POS (bartenders/admins). Has `Role` (NotAuthorized, Bartender, Admin).
- **Member** — club member who places orders (customer). Tracks `LastOrder`.
- **Order** — purchase record: club, bartender, member, contents (JSON `[]Line`), price (cents), status.
- **Status** — order lifecycle: `Open` → `Billed` → `Paid` (or `Cancelled`).
- **Item** — product in catalog, belongs to a `Category`. Has `PriceGladiators`, `PriceParabool`, `PriceCalamari`.
- **Price** — `int` in cents, formatted as currency.

## Auth

- Cookie-based token auth (`auth_token` cookie)
- 24h token duration, refreshed on activity
- Passwords hashed with bcrypt
- Roles: `Bartender` (PermissionBarStuff), `Admin` (BarStuff + AdminStuff)
- Default admin users created on first run (see `main.go`)

## Routes

### Public
| Method | Path | Description |
|--------|------|-------------|
| GET | `/login` | Login page |
| POST | `/login` | Login handler |
| GET | `/version` | Build version info |
| GET | `/static/*` | Static assets (compressed) |

### Authenticated (any role)
| Method | Path | Description |
|--------|------|-------------|
| GET | `/logout` | Logout |
| POST | `/active` | Refresh session |
| GET | `/profile` | Profile page |
| POST | `/profile/password` | Change password |
| POST | `/profile/name` | Change display name |

### Bartender (PermissionBarStuff)
| Method | Path | Description |
|--------|------|-------------|
| GET | `/order` | Order page (main POS UI) |
| GET | `/api/member/{id}` | Member details JSON |
| POST | `/api/order` | Place order JSON |
| GET | `/history` | Order history page |
| POST | `/history/{id}/delete` | Cancel order |
| GET | `/leaderboard` | Leaderboard page |

### Admin (PermissionAdminStuff) — `/admin` prefix
| Method | Path | Description |
|--------|------|-------------|
| GET/POST | `/admin/users` | User management |
| POST | `/admin/users/{id}/delete` | Delete user |
| GET/POST | `/admin/members` | Member management |
| POST | `/admin/members/{id}/delete` | Delete member |
| GET | `/admin/catalog` | Catalog management |
| POST | `/admin/catalog/category` | Create/update category |
| POST | `/admin/catalog/category/{id}/delete` | Delete category |
| POST | `/admin/catalog/item` | Create/update item |
| POST | `/admin/catalog/item/{id}/delete` | Delete item |
| GET | `/admin/billing` | Billing overview |
| GET | `/admin/download` | Download billing CSV |

## Frontend

- Server-side rendered Go templates (`templates/`)
- All templates extend `base.html` and include `nav.html`
- CSS: [BeerCSS](https://www.beercss.com/) (Material Design), dark theme
- JS: vanilla, minimal — `order.js` for POS interactivity, `activity.js` for session keep-alive
- Template rendering: `templates.Render(w, "page.html", data)` — executes `"base"` block
- Custom template funcs: `clubClass` (club → CSS class), `add` (integer addition)

## Database

- PostgreSQL, raw SQL (no ORM), `lib/pq` driver
- Repository pattern with interfaces in `backend/infrastructure/repo/`
- Migrations: `backend/infrastructure/repo/postgres/migrations/NNNN.sql`
  - Auto-applied on startup via version table
  - To add: create next numbered `.sql` file (e.g. `0003.sql`)
- Tables: `users`, `categories`, `items`, `members`, `orders`, `version`

## Configuration

Via TOML config file, environment variables (`STREEPJES_` prefix), or CLI flags. Priority: flags > env > file > defaults.

| Setting | Env var | Default |
|---------|---------|---------|
| `port` | `STREEPJES_PORT` | `80` |
| `db_connection_string` | `STREEPJES_DB_CONNECTION_STRING` | `postgresql://postgres@127.0.0.1:5432?sslmode=disable` |
| `disable_secure` | `STREEPJES_DISABLE_SECURE` | `false` |
| `log_level` | `STREEPJES_LOG_LEVEL` | `info` |
| `tls_cert_path` | `STREEPJES_TLS_CERT_PATH` | `streepjes.pem` |
| `tls_key_path` | `STREEPJES_TLS_KEY_PATH` | `key.pem` |

See `streepjes.example.toml` for config file format.

## Development

### Requirements

- Go >= 1.26
- [just](https://github.com/casey/just) (command runner)
- [entr](https://eradman.com/entrproject/) (hot-reload, for `just run`)
- PostgreSQL (or use `docker-compose up -d`)

### Commands

| Command | Description |
|---------|-------------|
| `just run` | Dev server with hot-reload (watches `.go`, `.html`, `.js`) |
| `just run-once` | Run without file watching |
| `just test` | Run tests with coverage |
| `just lint` | Run golangci-lint |
| `just generate` | Run `go generate ./...` (enumer) |
| `just build` | Build binary to `./bin/streepjes` |
| `just container` | Build container image with podman |

### Setup

1. Start PostgreSQL: `docker-compose up -d`
2. Copy `streepjes.example.toml` to config or set env vars in `.env.dev`
3. `just run`
4. Login with default admin credentials (logged on first startup)

## Testing

- Standard Go tests, co-located with source (`*_test.go`)
- Mock repos in `backend/infrastructure/repo/mockdb/` — structs with function fields
- Assertion library: `git.fuyu.moe/Fuyu/assert`
- Run: `just test`

## How to Implement a New Feature

### Adding a new domain entity

1. Define type in `domain/orderdomain/` (or new subdomain)
2. If enum, add `//go:generate go tool enumer ...` directive, run `just generate`
3. Add migration in `backend/infrastructure/repo/postgres/migrations/NNNN.sql`

### Adding a new repository

1. Define interface in `backend/infrastructure/repo/` (e.g. `repo.MyEntity`)
2. Implement in `backend/infrastructure/repo/postgres/`
3. Add mock in `backend/infrastructure/repo/mockdb/`
4. Wire in `main.go` (create repo, pass to service)

### Adding a new service method

1. Add method to service interface (`backend/application/{service}/`)
2. Implement on the `service` struct in the same package
3. Write tests using mockdb

### Adding a new page

1. Create template in `templates/` (or `templates/admin/`)
2. Register in `templates/templates.go` `pageFiles` slice
3. Template must define blocks expected by `base.html`
4. Add route in appropriate file under `backend/infrastructure/router/`:
   - `public.go` — unauthenticated
   - `auth.go` — any authenticated user
   - `bartender.go` — requires PermissionBarStuff
   - `admin.go` — requires PermissionAdminStuff
5. Write handler function returning `http.HandlerFunc` (closure over service)
6. Use `render(w, "template.html", data)` to render, `newPageData(r, "pagename")` for nav state

### Adding a new API endpoint

1. Add route in `bartender.go` or `admin.go`
2. For JSON request bodies: use `chio.JSON(handler)` wrapper
3. For JSON responses: use `chio.WriteJSON(w, status, data)`

### Adding a new migration

1. Create `backend/infrastructure/repo/postgres/migrations/NNNN.sql` (next number)
2. Migrations auto-apply on startup — no manual step needed

### Code generation

- Enums use [enumer](https://github.com/dmarkham/enumer) for JSON/SQL/String methods
- Add `//go:generate go tool enumer -json -sql -linecomment -type MyType` to source file
- Run `just generate`
