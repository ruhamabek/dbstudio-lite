# DBStudio Lite v0.1.0

Initial release of DBStudio Lite, a lightweight, modular desktop PostgreSQL client built with Go, Next.js, and Wails v2.

DBStudio Lite was developed using strict Test-Driven Development (TDD) and dependency injection principles, decoupling core database domain logic from presentation and transport layers.

---

### Highlights

- **Connection Management & Profiles**:
  - Connect via standard PostgreSQL connection URLs or discrete parameters (Host, Port, User, Password, Database, SSL Mode).
  - Save, list, and delete connection profiles with masked credential previews.
  - Profiles are persisted locally to disk using user-restricted permissions (`0600`) within the standard operating system configuration directory.

- **Schema & Metadata Inspection**:
  - Automatic discovery and categorization of user tables grouped by schema (`public`, custom schemas).
  - Fast table filtering and search.
  - Column metadata inspection including SQL data types, nullability constraints, default expressions, and composite primary key detection.

- **Interactive Tabular Data Viewer**:
  - Browse table records with clean visual handling of `NULL` values, booleans, timestamps, numbers, and structured JSON fields.
  - Telemetry reporting row counts and query execution latency in milliseconds.

- **Monaco SQL Console**:
  - Full-featured SQL editor powered by Monaco Editor.
  - Syntax highlighting, autocomplete, line numbers, bracket matching, and auto-indentation.
  - Keyboard shortcuts (`Ctrl+Enter` / `Cmd+Enter` to execute).
  - Resizable split pane to adjust vertical proportions between editor and results.
  - 100% offline bundling with zero external CDN dependencies.

---

### Architecture & Engineering

- **Interface-Driven Core**: Core database capabilities (`internal/db`) depend exclusively on lean Go interfaces (`Querier`, `Rows`, `Client`, `ConnectionStore`). The entire unit test suite executes in-memory in milliseconds with zero live database dependencies.
- **Concurrency Safety**: Connection lifecycles and client dispatch utilize read-write mutex locks (`sync.RWMutex`) to prevent race conditions during asynchronous frontend operations.
- **Decoupled Driver**: Utilizes `jackc/pgx/v5` through an internal adapter layer, ensuring the engine remains driver-agnostic.
- **Centralized Frontend State**: Employs a dedicated `DatabaseContext` and reusable UI components built on Tailwind CSS v4 and shadcn/ui design tokens.

---

### Assets & Platform Binaries

Pre-compiled standalone binaries generated via GitHub Actions:

| Platform | Architecture | File |
| :--- | :--- | :--- |
| **Linux** | x86_64 (glibc 2.35+) | `DBStudio-Lite-Linux-x86_64` |
| **macOS** | Universal (Intel + Apple Silicon) | `DBStudio-Lite-macOS-Universal.zip` |
| **Windows** | x64 | `DBStudio-Lite-Windows-x64.exe` |

---

### Verification

All backend packages have been tested with Go's race detector enabled:
```bash
go test -v -race ./...
```
