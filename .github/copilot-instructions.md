## GitHub Copilot Agent Instructions for GopherMate

### Project Overview
GopherMate is a fullstack online chess system with a React (JavaScript) frontend, a Go backend (using only standard libraries), and a PostgreSQL database. The project is modular, simple, and uses REST APIs for all communication. No ORMs or external router libraries are used. All chess logic (move validation, checkmate detection) is implemented on both frontend and backend.

### Coding Guidelines

#### General
- Always use a three step process when asked a new implementation:
  1. **Read**: Read any relevant files and code to this implementation.
  2. **Plan**: Outline the approach and structure.
  3. **Implement**: Ask if there are any changes to the plan, if ok then write the code.
- Use the project structure and naming conventions as described in `docs/project_summary.md`.
- Use plain JavaScript for React code (no TypeScript).
- Use Go standard library for backend (no frameworks, no external router libraries).
- Use raw SQL queries for all database access (no ORM).
- Use UUIDs for resource identifiers (games, sessions).
- Use REST endpoints for all game communication.
- Keep code modular, simple, and readable.
- Use Bootstrap components for UI, but keep third-party libraries to a minimum.
- CSS files should be placed next to their corresponding JS files.

#### Frontend (React)
- Use Vite for project setup.
- Use React Router for routing.
- Use local component state, React Context, or minimal custom logic for state management.
- Organize code into `components/`, `features/`, `pages/`, `services/`, `utils/`, and `chess/` as per the project summary.
- Implement chess move validation logic in `src/chess/moveValidator.js` and related files.
- Use REST API calls for all backend communication (no WebSockets).
- Use polling for game state updates.

#### Backend (Go)
- Use `net/http` and `database/sql` from the standard library.
- Use `lib/pq` for PostgreSQL driver.
- Use `golang.org/x/crypto/bcrypt` for password hashing.
- Organize code into `cmd/`, `internal/api/`, `internal/db/`, `internal/model/`, and `internal/utils/` as per the project summary.
- Implement authentication using custom session tokens stored in the database.
- Write all SQL queries manually in the repository files.

#### Database
- Use the schema provided in `docs/project_summary.md` for `users`, `games`, `moves`, and `sessions` tables.
- Use raw SQL for all migrations and queries.

#### Endpoints
- Follow the REST API endpoints and methods as described in the project summary for authentication, game management, and gameplay.

#### Best Practices
- Write clear, concise, and well-documented code.
- Keep logic separated by concern (UI, API, DB, utils, etc.).
- Avoid unnecessary dependencies.
- Use comments to explain non-obvious logic, especially for chess rules and move validation.

### Example File Structure
Refer to `docs/project_summary.md` for the canonical file and folder structure for both frontend and backend.

### Environment
- Node.js version: 22.13
- Go version: 1.22
- OS: Windows

---
These instructions are for GitHub Copilot Agent Mode. Use them to ensure code completions and suggestions are consistent with the GopherMate project standards.
