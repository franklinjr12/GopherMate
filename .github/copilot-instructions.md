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

These instructions are for GitHub Copilot Agent Mode. Use them to ensure code completions and suggestions are consistent with the GopherMate project standards.
### Extended Reference: Routes, Models, and Schema

#### REST API Endpoints (Full List)
See `docs/project_summary.md` for full list. Examples:
  - `POST /api/register` — Register user
  - `POST /api/login` — Login, returns session token
  - `POST /api/logout` — Logout, invalidate session token
  - `GET /api/me` — Get current logged-in user info
  - `GET /api/games` — List open games
  - `POST /api/games` — Create game
  - `POST /api/games/:id/join` — Join game
  - `GET /api/games/:id` — Get full game state
  - `POST /api/games/:id/move` — Submit move
  - `GET /api/games/:id/moves` — Get move history
  - `POST /api/games/:id/resign` — Resign from game

#### Database Models & Schema
  - `users`: id (SERIAL, PK), username (unique), email, password_hash, created_at
  - `games`: id (UUID, PK), player_white_id, player_black_id, winner, created_at, finished_at
  - `moves`: id (SERIAL, PK), game_id (UUID, FK), player_id, move_number, notation, created_at
  - `sessions`: token (UUID, PK), user_id, created_at, expires_at
  - See `backend/db/schema.sql` and `docs/project_summary.md` for full SQL definitions.

#### Backend Models
  - `internal/model/user.go`: User struct
  - `internal/model/game.go`: Game struct
  - `internal/model/move.go`: Move struct

#### Frontend Models
  - Chess piece logic in `src/chess/pieces/` (pawn.js, knight.js, etc.)
  - Board and move validation in `src/chess/board.js` and `src/chess/moveValidator.js`

#### Frontend Pages & Features
  - Auth: `src/pages/LoginPage.jsx`, `src/pages/RegisterPage.jsx`
  - Lobby: `src/pages/GamesPage.jsx`, game creation/joining in `src/features/lobby/`
  - Game session: `src/pages/GameSessionPage.jsx`, board UI in `src/chess/board.jsx`, move log in `src/pages/MoveLog.jsx`

#### Backend Endpoints
  - Auth: `internal/api/auth.go`
  - Game: `internal/api/game.go`, board state in `internal/api/board.go`
  - Middleware: `internal/api/middleware.go` (session validation)

#### Backend DB Access
  - User: `internal/db/user_repository.go`
  - Game/move: `internal/db/game_repository.go`, `internal/db/move.go`
  - Session: `internal/db/session_repository.go`

#### Backend Utilities
  - Config: `internal/utils/config.go`
  - Password hashing: `internal/utils/hash.go`
  - Token generation: `internal/utils/token.go`
  - JSON response helpers: `internal/utils/response.go`
