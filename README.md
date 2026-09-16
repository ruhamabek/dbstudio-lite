# DBStudio Lite

[![Total Downloads](https://img.shields.io/github/downloads/ruhamabek/dbstudio-lite/total?style=flat-square&color=000000&label=Downloads)](https://github.com/ruhamabek/dbstudio-lite/releases)
[![Latest Release](https://img.shields.io/github/v/release/ruhamabek/dbstudio-lite?style=flat-square&color=000000&label=Release)](https://github.com/ruhamabek/dbstudio-lite/releases)
[![License](https://img.shields.io/github/license/ruhamabek/dbstudio-lite?style=flat-square&color=000000)](https://github.com/ruhamabek/dbstudio-lite/blob/main/LICENSE)

DBStudio Lite is a lightweight, cross-platform desktop database client for **PostgreSQL**, **MySQL**, and **SQLite**, built with Go and Next.js, packaged via Wails v2. The project emphasizes modular software design, test-driven development (TDD), strict dependency injection, and clean separation between backend domain logic and user interface presentation.

## Screenshots

### Table Data Exploration
![Table Data Exploration](screenshots/2.jpg)

### SQL Console with Monaco Editor
![SQL Console](screenshots/3.jpg)

### Connection Setup and Profile Management
![Connection Dialog](screenshots/4.png)

### Saved Connection Profiles
![Saved Connections](screenshots/1.jpg)

## Features

- Multi-Database Support: Connect to PostgreSQL, MySQL, and SQLite databases seamlessly using connection URIs, file paths, or discrete parameters.
- Persistent Profiles: Securely save, list, and delete connection profiles on local disk (`0600` file permissions in standard OS configuration directories).
- Polymorphic Schema Inspection: Automatically discover and group tables by schema across all supported engines (`information_schema` for PostgreSQL and MySQL, `sqlite_master` and `PRAGMA table_info` for SQLite).
- Column and Key Metadata: Inspect column definitions, SQL data types, nullability constraints, default expressions, and composite primary keys.
- Tabular Data Viewer: Explore table records with formatted handling of NULL values, booleans, timestamps, numbers, and structured JSON fields.
- Monaco SQL Editor: Execute arbitrary SQL queries with integrated syntax highlighting, line numbers, autocomplete, keyboard shortcuts (`Ctrl+Enter` or `Cmd+Enter`), and resizable split view. Bundled 100% offline with zero CDN dependencies.
- Execution Telemetry: Track row count and execution duration in milliseconds for every query.
- Thread-Safe Architecture: Concurrency-safe backend connection management using read-write mutex locks (`sync.RWMutex`).

## Architecture and Engineering Principles

DBStudio Lite was built from the ground up using strict Test-Driven Development (TDD) practices:

### 1. Interface-Driven Core (`internal/db`)
The core inspection and execution engine depends exclusively on lean Go interfaces (`Querier`, `Rows`, `Client`, `SchemaInspector`, `Connector`, `ConnectionStore`). It contains zero direct dependencies on external network sockets or third-party drivers. As a result, the entire unit test suite runs in-memory using lightweight fakes in milliseconds without requiring external Docker containers or live database instances.

### 2. Dependency Injection
Every component receives its dependencies explicitly through constructors or factory functions:
- `SchemaInspector` is constructed via `NewInspector(driver DriverType, db Querier)`.
- `Runner` accepts a `Querier` interface.
- `ConnectionManager` accepts a functional `Connector` constructor (`func(ctx, cfg) (Client, error)`).
- `App` receives `ConnectionManager` and `ConnectionStore`.

### 3. The Adapter Pattern & Pure-Go Drivers
Third-party database drivers are isolated behind domain adapters:
- **PostgreSQL**: `github.com/jackc/pgx/v5` (`pgxpool`).
- **MySQL**: `github.com/go-sql-driver/mysql` (pure Go).
- **SQLite**: `modernc.org/sqlite` (pure Go, CGO-free, enabling cross-compilation across Linux, macOS Universal, and Windows without C toolchains).

Rather than leaking driver data structures into application business logic, thin adapter layers reconcile driver signatures with internal domain interfaces.

### 4. Structured Frontend State Management
The frontend avoids fragmented component state by centralizing connection, schema, and query lifecycles in a dedicated `DatabaseContext`. All styling utilizes semantic CSS design tokens defined in `globals.css` with zero hardcoded color values.

## Tech Stack

- Backend: Go 1.25, Wails v2
- Database Drivers:
  - PostgreSQL: `jackc/pgx/v5` (`pgxpool`)
  - MySQL: `go-sql-driver/mysql`
  - SQLite: `modernc.org/sqlite`
- Frontend: Next.js 15, React 19, TypeScript
- Editor: Monaco Editor (`@monaco-editor/react`)
- UI Primitives: Tailwind CSS v4, shadcn/ui, Lucide Icons

## Project Structure

```text
dbstudio-lite/
├── app.go                      # Wails IPC application bindings and service layer
├── main.go                     # Desktop application entrypoint and options
├── internal/
│   └── db/                     # Core database package (100% test-driven)
│       ├── config.go           # DSN building and multi-engine URL parsing
│       ├── config_test.go      # Configuration and URL parsing tests
│       ├── drivers.go          # Multi-driver connector dispatcher
│       ├── generic_sql.go      # database/sql adapter for MySQL and SQLite
│       ├── inspector.go        # PostgreSQL, MySQL, and SQLite schema inspectors
│       ├── inspector_test.go   # Multi-engine schema inspection unit tests
│       ├── runner.go           # Dynamic tabular query execution and timing
│       ├── runner_test.go      # Query execution unit tests
│       ├── manager.go          # Thread-safe connection lifecycle management
│       ├── manager_test.go     # Connection management unit tests
│       ├── store.go            # Connection profile persistence implementation
│       ├── store_test.go       # Storage unit tests
│       └── postgres.go         # pgx driver adapter implementation
├── frontend/
│   ├── src/
│   │   ├── app/                # Next.js App Router (layout, page, globals.css)
│   │   ├── components/
│   │   │   ├── connection/     # Multi-database connection bar and modal
│   │   │   ├── console/        # Monaco SQL editor and resizable split view
│   │   │   ├── sidebar/        # Table search and schema grouping sidebar
│   │   │   ├── ui/             # Reusable shadcn UI components
│   │   │   └── viewer/         # Tabular data and structure viewers
│   │   └── context/            # Centralized DatabaseContext state
│   └── wailsjs/                # Generated Go-to-TypeScript runtime bindings
├── website/                    # Standalone Next.js marketing landing page (for Vercel)
└── screenshots/                # Application documentation screenshots
```

## Getting Started

### Prerequisites

Ensure the following tools are installed on your machine:
- Go (version 1.22 or higher)
- Node.js (version 18 or higher) and npm
- Wails CLI v2 (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### Running in Development Mode

To start the application with live reloading for both Go backend code and Next.js frontend assets:

```bash
wails dev
```

### Running Backend Unit Tests

All unit tests are isolated and run in-memory:

```bash
go test -v -race ./...
```

### Building for Production

To compile an optimized, standalone desktop binary:

```bash
wails build
```

The output executable will be placed in the `build/bin/` directory.

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). See the [LICENSE](LICENSE) file for complete license terms.
