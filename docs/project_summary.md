# 📜 Project Definition: Online Chess Playing System

## 🌟 Objective

Build a simple online chess playing system to practice and showcase fullstack skills using:

- **React (JavaScript)** for the frontend
- **Go** for the backend, using mostly standard libraries
- **PostgreSQL** as the database

---

## 🛡️ Architecture Overview

### Frontend

- **Framework**: React (plain JavaScript, no TypeScript), Using Vite for project setup
- **Routing**: React Router
- **State Management**: Local component state, React Context, or minimal custom logic
- **UI Libraries**: Keep it minimal; avoid third-party libraries unless absolutely necessary
- **Move Logic**: Use simple custom logic im for validation on the client side, validation happens again on backend side
- **Communication**: REST API with polling (no WebSockets for now)

### Backend

- **Language**: Go
- **Framework**: Standard library only (`net/http`, `database/sql`)
- **Database Driver**: `lib/pq` (PostgreSQL)
- **Password Hashing**: `golang.org/x/crypto/bcrypt` (minimal cryptographic lib usage)
- **Authentication**: Custom session-based (token stored in database)
- **Game Sync**: REST + polling

### Database

- **Engine**: Postgresql

---

## 📁 Project Structure

### Frontend (React)

```
frontend/
├── public/
├── src/
│   ├── components/        # Shared reusable UI components
│   ├── features/
│   │   ├── auth/          # Login/Register forms and logic
│   │   |   |── LoginForm.js
│   │   |   |── RegisterForm.js
│   │   ├── lobby/         # Lobby list and game creation/joining
│   │   |   |── LobbyList.js
│   │   |   |── CreateGameButton.js
│   │   |   |── JoinGameButton.js
│   │   └── game/          # Game board, move history, timers
│   │   |   |── ChessBoard.js
│   │   |   |── MoveHistory.js
│   │   |   |── ResignButton.js
│   ├── pages/             # Route-level components
│   |   |── LoginPage.js
│   |   |── RegisterPage.js
│   |   |── LobbyPage.js
│   |   |── GamePage.js
│   ├── services/          # API interaction (REST)
│   |   |── authService.js
│   |   |── gameService.js
│   |   |── moveService.js
│   ├── utils/             # Utility functions
│   |   |── auth.js
│   |   |── api.js
|   |── chess/
|   |   |── pieces/
|   |   |   |── pawn.js
|   |   |   |── tower.js
|   |   |   |── knight.js
|   |   |   |── bishop.js
|   |   |   |── queen.js
|   |   |   |── king.js
|   |   |── moveValidator.js
|   |   |── board.js
│   └── App.js
└── package.json
```

#### Additional Info
- CSS files should be beside their .js file, so for example if there is a LobbyList.js its styles will be on LobbyList.css on the same folder
- Will use Bootstrap components

### Backend (Go)

```
backend/
├── cmd/
│   └── server/
│       └── main.go                 # Starts HTTP server, initializes DB and routes
├── internal/
│   ├── api/
│   │   ├── auth.go                # login, register, session handling
│   │   ├── game.go                # game creation, join, state fetch, move posting
│   │   └── middleware.go          # auth/session validation middleware
│   ├── db/
│   │   ├── postgres.go            # DB connection setup
│   │   ├── user_repository.go     # User DB operations
│   │   ├── game_repository.go     # Game and move DB operations
│   │   └── session_repository.go  # Session token handling
│   ├── model/
│   │   ├── user.go
│   │   ├── game.go
│   │   └── move.go
│   └── utils/
│       ├── config.go              # Configuration handling, global consts
│       ├── hash.go                # Password hashing/verification
│       ├── token.go               # UUID or token generation
│       └── response.go            # JSON response helpers
├── schema.sql                     # DB schema (version 1)
└── go.mod / go.sum
```

📂 cmd/server/main.go
- Connecting to PostgreSQL
- Setting up routes and starting the HTTP server

📂 internal/api/
- Controllers/handlers grouped by purpose
- Manages HTTP-specific concerns (request parsing, status codes)
- Delegates to db/ and model/ for actual work

📂 internal/db/
- Pure data access logic
- Use database/sql with lib/pq
- Can define query helpers and wrap transactions here if needed

📂 internal/model/
- Struct definitions for User, Game, Move, etc.
- Pure Go, unaware of HTTP or SQL
- Might include basic validation methods (e.g., IsValidMove() if needed)

📂 internal/utils/
- Tiny helpers to avoid clutter in api/
- Focused, independent logic like:
- - GenerateToken() string
- - HashPassword(pw string) ([]byte, error)
- - WriteJSON(w, status, data)
- Configuration and global consts

#### Endpoints

🔐 Authentication Endpoints
Method |  Endpoint       |  Description                      |  Auth Required
POST   |  /api/register  |  Register a new user              |  ❌
POST   |  /api/login	 |  Log in and get session token     |  ❌
POST   |  /api/logout	 |  Invalidate session token         |  ✅
GET	   |  /api/me	     |  Get current logged-in user info  |  ✅

🎮 Game Management
Method | Endpoint              | Description                         | Auth Required
GET    | /api/games            | List all open games (status = open) | ✅
POST   | /api/games            | Create a new game                   | ✅
POST   | /api/games/:id/join   | Join an existing game by ID         | ✅
GET    | /api/games/:id        | Get full game state                 | ✅

♟️ Moves and Gameplay
Method | Endpoint                  | Description                                  | Auth Required
POST   | /api/games/:id/move       | Submit a move (e.g. {"from":"e2","to":"e4"}) | ✅
GET    | /api/games/:id/moves      | Get move history (for polling)               | ✅
POST   | /api/games/:id/resign     | Resign from the game                         | ✅

---

## 🗃 Database Schema

### `users`

Stores player credentials.

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### `games`

Stores metadata about matches.

```sql
CREATE TABLE games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_white_id INTEGER REFERENCES users(id),
    player_black_id INTEGER REFERENCES users(id),
    winner TEXT, -- 'white', 'black', 'draw', or NULL
    created_at TIMESTAMP DEFAULT NOW(),
    finished_at TIMESTAMP
);
```

### `moves`

Tracks the progression of a game.

```sql
CREATE TABLE moves (
    id SERIAL PRIMARY KEY,
    game_id UUID REFERENCES games(id) ON DELETE CASCADE,
    player_id INTEGER REFERENCES users(id),
    move_number INTEGER,
    notation TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### `sessions`

```sql
CREATE TABLE sessions (
    token UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP
);
```

---

## ✅ Features (Initial Scope)

### Authentication

- Register new user
- Login with username + password
- Token-based session management (manual, not JWT)

### Lobby

- Create a new game
- List open games to join
- Join existing game

### Gameplay

- Basic chessboard UI
- Move submission and history
- Validation with custom logic implemented on code (client and server)
- Polling for opponent moves
- Resign

---

## 🧐 AI Agent Coding Notes

- The agent persona is a senior software developer with great knowledge of react js, golang and postgresql.
- Use REST endpoints for all game communication.
- Keep code modular but simple, no ORMs, no external router libraries.
- Use plain JavaScript in React with good component separation.
- Write raw SQL queries for all DB access.
- Assume chess logic (move validation, checkmate detection) is handled in both frontend and backend.
- Use UUIDs for resource identifiers (games, sessions).
- Environment being used for development is Windows unless other is specified.
- node version is 22.13
- go version is 1.22
