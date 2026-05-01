# MindBridge API

A clean‑architecture REST API built with Go, Gin, MongoDB, and Redis. It provides authentication, AI chat (OpenAI‑compatible model names), image editing, and real‑time streaming.

---

## Technology Stack
- **Language**: Go 1.26
- **Web Framework**: Gin
- **Database**: MongoDB
- **Cache / Session Store**: Redis
- **Authentication**: JWT + Bcrypt
- **Real‑time**: Server‑Sent Events (SSE) & WebSocket
- **Frontend**: React + TypeScript (in `web/`)
- **CI/CD**: GitHub Actions with lint, test, Docker build, and Render deploy

---

## Project Structure
```
MindBridge/
├── cmd/server/                 # main entry point
├── config/                     # env handling & model mapping
├── domain/                     # entities & repository interfaces
├── application/                # DTOs & use‑cases
├── infrastructure/             # concrete implementations (Mongo, Redis, JWT, etc.)
├── middleware/                 # request‑ID, rate‑limit, etc.
├── presentation/               # HTTP handlers
│   └── handlers/
├── utils/                      # validators & logging helpers
├── web/                        # React frontend
└── .claude/                    # internal Claude Code files (ignored by Git)
```

---

## Setup & Run Locally
1. **Clone & install dependencies**
   ```bash
   git clone https://github.com/yourorg/MindBridge.git
   cd MindBridge
   go mod tidy
   ```
2. **Create environment file**
   ```bash
   cp .env.example .env
   # edit .env with your credentials
   ```
3. **Run the API server**
   ```bash
   go run ./cmd/server
   ```
   The server listens on `http://127.0.0.1:5000` (configurable via `SERVER_HOST` & `SERVER_PORT`).
4. **Run the frontend**
   ```bash
   cd web
   npm install
   npm run dev
   ```
   Frontend will be available at `http://localhost:5173`.

---

## Docker
```bash
# Build the image
docker build -t mindbridge .
# Run with Docker Compose (includes Redis & Mongo placeholders)
docker compose up --build
```
The API will be reachable at `http://0.0.0.0:5000`.

---

## API Overview
### Authentication (public)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST   | `/auth/register` | Register a new user |
| POST   | `/auth/login`    | Obtain JWT token |
| POST   | `/auth/logout`   | Invalidate session |

### Public Endpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | `/`      | API health/info |
| GET    | `/models`| List OpenAI‑compatible model names |

### Protected Endpoints (JWT required)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST   | `/chat/stream-raw` | Chat with AI (SSE) |
| POST   | `/upload`          | Upload image/file |

---

## Model Mapping
The API accepts **OpenAI‑compatible** model identifiers (e.g., `gpt-4o`, `claude-sonnet-4-6`). Internally these map to gateway names via `config/model_mapping.go`. See the **Available Models** table in the API docs for the full list.

---

## Testing
```bash
go test ./...   # run all unit tests
```
Tests cover middleware, handlers, and use‑case logic. CI runs `golangci-lint`, tests, and Docker build automatically.

---

## Contributing
1. Fork the repository.
2. Create a feature branch.
3. Write tests for new behavior (TDD is encouraged).
4. Ensure `go test ./...` and `golangci-lint run` pass.
5. Open a Pull Request.

---

## License
MIT License – see the `LICENSE` file.

---

## Security
For security‑related concerns, see the [SECURITY.md](SECURITY.md) file.
