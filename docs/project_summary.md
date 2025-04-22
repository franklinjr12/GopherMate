# 📜 Project Definition: Online Chess Playing System

## 🌟 Objective

Build a simple online chess playing system to practice and showcase fullstack skills using:

- **React (JavaScript)** for the frontend
- **Go** for the backend, using mostly standard libraries
- **PostgreSQL** as the database

---

## 🛡️ Architecture Overview

### Frontend

- **Framework**: React (plain JavaScript, no TypeScript)
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
│   │   ├── lobby/         # Lobby list and game creation/joining
│   │   └── game/          # Game board, move history, timers
│   ├── pages/             # Route-level components
│   ├── services/          # API interaction (REST)
│   ├── utils/             # Utility functions
│   └── App.js
└── package.json
```

### Backend (Go)

```
backend/
├── cmd/
│   └── server/            # Main entry point (main.go)
├── internal/
│   ├── api/               # HTTP handlers grouped by domain (auth, game, etc.)
│   ├── db/                # Database access logic using raw SQL
│   ├── model/             # Core domain structs
│   ├── utils/             # Helper functions (e.g., hashing, UUIDs)
└── schema.sql             # PostgreSQL schema
```

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

- Use REST endpoints for all game communication.
- Keep code modular but simple, no ORMs, no external router libraries.
- Use plain JavaScript in React with good component separation.
- Write raw SQL queries for all DB access.
- Assume chess logic (move validation, checkmate detection) is handled in both frontend and backend.
- Use UUIDs for resource identifiers (games, sessions).
- Environment being used for development is Windows unless other is specified.
- node version is 22.13
- go version is 1.22
