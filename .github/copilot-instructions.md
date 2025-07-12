## GitHub Copilot Agent Instructions for GopherMate

### Project Architecture & Purpose
GopherMate is a fullstack online chess system for practicing and showcasing fullstack skills. It uses:
- **Frontend:** React (JavaScript, Vite, React Router, no TypeScript)
- **Backend:** Go (standard library only, no frameworks)
- **Database:** PostgreSQL (raw SQL, no ORM)
All communication is via REST APIs. Chess logic (move validation, checkmate detection) is implemented on both frontend and backend.

### Key Conventions & Patterns
- **File Structure:**
  - See `docs/project_summary.md` for canonical structure. Follow the separation: `frontend/` (React), `backend/` (Go), with further breakdown into `components/`, `features/`, `pages/`, `services/`, `chess/` (frontend) and `cmd/`, `internal/api/`, `internal/db/`, `internal/model/`, `internal/utils/` (backend).
  - CSS files are placed next to their corresponding JS files.
- **Identifiers:** Use UUIDs for games and sessions.
- **Database:** Write all SQL queries manually in Go files. Use the schema in `docs/project_summary.md` and `backend/db/schema.sql`.
- **Authentication:** Custom session tokens stored in the database (not JWT).
- **Frontend State:** Use local state, React Context, or minimal custom logic. No Redux or similar libraries.
- **UI:** Use Bootstrap components if needed, but keep third-party libraries minimal.
- **Communication:** Use REST endpoints and polling (no WebSockets).

### Developer Workflows
- **Frontend:**
  - Start dev server: `npm run dev` in `frontend/`
  - Build: `npm run build`
  - Lint: `npm run lint`
  - Main entry: `src/main.jsx`, routes in `src/pages/`
  - API calls: via `src/services/` (e.g., `authService.js`, `gameService.js`)
- **Backend:**
  - Start server: build and run `cmd/main.go` (e.g., `go run ./cmd/main.go`)
  - DB connection: configured via environment variables or `.env` (see `internal/utils/config.go`)
  - Endpoints: defined in `internal/api/` (e.g., `auth.go`, `game.go`)
  - All SQL is raw, in `internal/db/` files
- **Database:**
  - Schema: see `backend/db/schema.sql` and `docs/project_summary.md`
  - Use psql or a GUI for migrations; no migration tool is included

### REST API Endpoints
- See `docs/project_summary.md` for full list. Examples:
  - `POST /api/register` — Register user
  - `POST /api/login` — Login, returns session token
  - `GET /api/games` — List open games
  - `POST /api/games` — Create game
  - `POST /api/games/:id/join` — Join game
  - `POST /api/games/move` — Submit move

### Project-Specific Notes
- **Move validation** is implemented in both frontend (`src/chess/`) and backend (`internal/movevalidation/`).
- **Session tokens** are required for all authenticated endpoints; stored in localStorage on frontend.
- **No ORMs, no external router libraries, no Redux.**
- **Chess logic**: Both client and server validate moves. See `src/chess/` and `internal/movevalidation/`.
- **Environment:**
  - Node.js: 22.13
  - Go: 1.22
  - OS: Windows

### Best Practices
- Keep code modular and simple; separate UI, API, DB, and utility logic.
- Use comments for non-obvious logic, especially chess rules and move validation.
- Avoid unnecessary dependencies.
- Reference `docs/project_summary.md` for canonical structure and conventions.

---
These instructions are for GitHub Copilot Agent Mode. Use them to ensure code completions and suggestions are consistent with the GopherMate project standards.
