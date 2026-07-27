# real-time-forum

A single-page forum web application with real-time private messaging. Users can register, log in, create posts under categories, comment on posts, and chat privately with other users over WebSockets, with online/offline presence and message history — all without a page reload.

## Collaboration note

This was a paired/group project completed with a partner (originally hosted under the Gitea account `malalawi` at learn.reboot01.com). Full commit history from both contributors is preserved in this repository.

## About the exercise

This is a 01-edu / Reboot01 piscine project that builds on an earlier, simpler forum exercise by adding real-time features: a Single Page Application front end (a single `index.html` with all navigation handled in JavaScript), private messaging with live delivery over WebSockets, and online/offline user tracking — while restricting the front end to vanilla JavaScript, HTML, and CSS (no frameworks).

## How it works

- The Go backend (`main.go`, `api/`) registers routes through a central route registry (`api/core_route_registry.go`, `api/core_routes.go`) and serves both the static SPA assets and JSON API endpoints/WebSocket upgrades from the same process.
- Authentication uses hashed passwords (`internal/security/hash.go`) and cookie-backed sessions (`internal/cookie`), with users able to log in by nickname or email.
- Data (users, sessions, posts, categories, comments, private messages) is stored in SQLite via `internal/database`, with schema creation and seed logic in `schema_create.go` and `schema_seed.go`.
- Real-time messaging runs through `api/websocket` (`manager.go`, `client.go`, `broadcast.go`, `user_manager.go`) using `gorilla/websocket`: each connected client is tracked by the manager, which broadcasts new messages and presence changes to the relevant online users.
- The front end (`static/js`) is organized as a small MVC-style framework of its own: `models/` fetch and hold data, `views/` render markup, `controllers/` wire the two together per feature (posts, comments, private chat, notifications, user list, login/register, settings), and `modules/router.js` handles in-page navigation without reloading.
- Private chat history loads the most recent messages first and fetches older ones in batches as the user scrolls, using throttling/debouncing in the chat controller/module.

## Tech stack

- Go (module `reboot01.com/js/realtime-forum`, Go 1.23)
- `github.com/gorilla/websocket` for WebSocket connections
- `github.com/mattn/go-sqlite3` with SQLite for persistence
- `golang.org/x/crypto` (bcrypt) for password hashing
- `github.com/google/uuid` for ID generation
- Vanilla JavaScript, HTML, and CSS for the single-page front end (no frontend framework)

## Project structure

```
real-time-forum/
├── main.go                     # entry point: config, DB init, starts the HTTP/WebSocket server
├── go.mod / go.sum             # Go module and dependencies
├── Makefile                    # build/run/test/audit/migration commands
├── api/                        # HTTP layer
│   ├── core_server.go, core_routes.go, core_route_registry.go
│   ├── core_middleware.go, core_context.go, core_errors.go, core_helpers.go
│   ├── handler_*.go            # login, register, posts, comments, categories, user list/profile, websocket
│   └── websocket/              # connection manager, client, broadcast, presence tracking
├── internal/
│   ├── database/               # SQLite schema, seed data, and per-entity queries (user, post, comment, message, session, category)
│   ├── cookie/, security/      # session cookies, password hashing, UUIDs
│   ├── request/, response/     # JSON request/response helpers
│   ├── validator/, env/, time/, version/
├── static/
│   ├── index.html              # the single SPA page
│   ├── css/                    # base, variables, layout, component styles
│   ├── images/                 # SVG icons
│   └── js/
│       ├── models/ views/ controllers/   # per-feature MVC-style modules
│       └── modules/             # router, pagination, websocket client, utils
└── LICENSE.md
```

## Build & run

```bash
go mod tidy
go run main.go
```

or via the Makefile:

```bash
make run      # builds to /tmp/bin/api and runs it
make build    # builds only
make test     # go test -v -race -buildvcs ./...
make audit    # vet, staticcheck, govulncheck, formatting checks
```

Configuration is read from environment variables with defaults: `BASE_URL` (`http://localhost:8080`), `HTTP_PORT` (`8080`), `DB_DSN` (`db.sqlite`). Running with `-recreate` drops and recreates the SQLite database with seed data. Once running, the app is served at `http://localhost:8080`.
