# MindBridge API

A clean architecture REST API server with authentication and chat functionality.

## Tech Stack

- **Framework**: Gin (Go web framework)
- **Architecture**: Clean Architecture
- **Database**: MongoDB
- **Cache/Session**: Redis
- **Authentication**: JWT + Bcrypt
- **WebSocket**: gorilla/websocket

## Project Structure

```
MindBridge/
├── cmd/server/          # Entry point
├── config/               # Configuration
├── domain/              # Entities & interfaces
├── application/         # DTOs & use cases
├── infrastructure/      # MongoDB repositories & JWT service
├── presentation/        # HTTP handlers & middleware
└── utils/               # Utilities & validators
```

## Configuration

Create a `.env` file in the root directory by copying `.env.example`:

```bash
cp .env.example .env
```

Then update the values in `.env` with your credentials:

```env
# MongoDB
MONGO_USERNAME=your_username
MONGO_PASSWORD=your_password
MONGO_CLUSTER=cluster0.xxx.mongodb.net
MONGO_DB=mindbridge

# Redis (Cloud)
REDIS_HOST=redis-xxx.cloud.redislabs.com
REDIS_PORT=12345
REDIS_USERNAME=default
REDIS_PASSWORD=your_redis_password

# JWT
JWT_SECRET=your-secret-key-change-in-production

# Server
SERVER_HOST=127.0.0.1
SERVER_PORT=5000

# External Services
WEBSOCKET_URL=wss://agents.use.ai/agents/budget-agent
FILE_UPLOAD_URL=https://files.use.ai/upload
FILE_BASE_URL=https://files.use.ai
```

## Running the Server

```bash
go run ./cmd/server
```

Server runs at `http://127.0.0.1:5000`

## API Endpoints

### Authentication (Public)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/register` | Register new user |
| POST | `/auth/login` | Login and get token |

### Public Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | API info |
| GET | `/models` | List available models |

### Protected Endpoints (Requires JWT Token)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/chat` | Send chat message (sync) |
| POST | `/chat/stream` | Send chat message (SSE stream) |
| POST | `/upload` | Upload image/file |
| GET | `/file/:key` | Get uploaded file URL |

---

## API Usage

### POST /auth/register

**Request:**
```json
{
  "name": "Piyush Makwana",
  "username": "piyushmakwana",
  "email": "piyush@example.com",
  "password": "Piyush@1234"
}
```

**Validations:**
- `name`: Only letters and spaces (no numbers or special characters, no leading spaces)
- `username`: Must start with a letter, can contain letters and numbers (no special characters)
- `email`: Valid email format (temp email domains blocked)
- `password`: 8+ chars with uppercase, lowercase, number, special char. Must not contain name, username, or email

**Success Response (201):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "token": "eyJhbGci...",
    "user": {
      "id": "65f...",
      "name": "Piyush Makwana",
      "username": "piyushmakwana",
      "email": "piyush@example.com"
    }
  }
}
```

**Error Response (400) - Validation Error:**
```json
{
  "success": false,
  "errors": [
    {
      "field": "name",
      "message": "Name must contain only letters (no numbers, spaces, or special characters)"
    },
    {
      "field": "password",
      "message": "Password must be 8+ chars with uppercase, lowercase, number, special char. Must not contain your name, username, or email"
    }
  ]
}
```

**Error Response (400) - Already Registered:**
```json
{
  "success": false,
  "error": "email already registered"
}
```

### POST /auth/login

**Request:**
```json
{
  "email": "piyush@example.com",
  "password": "Piyush@1234"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGci...",
    "user": {
      "id": "65f...",
      "name": "Piyush Makwana",
      "username": "piyushmakwana",
      "email": "piyush@example.com"
    }
  }
}
```

**Error Response (401) - Invalid Credentials:**
```json
{
  "success": false,
  "error": "invalid email or password"
}
```

### GET /

```bash
curl http://127.0.0.1:5000/
```

### GET /models

```bash
curl http://127.0.0.1:5000/models
```

### POST /chat (Protected - Sync Response)

```bash
curl -X POST http://127.0.0.1:5000/chat \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"query": "Hello", "model": "gateway-claude-opus-4-1"}'
```

### POST /chat/stream (Protected - Real-time Streaming)

```bash
curl -X POST http://127.0.0.1:5000/chat/stream \
  -H "Content-Type: application/json" \
  -H "Accept: text/event-stream" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"query": "Hello", "model": "gateway-claude-opus-4-1"}'
```

**Response Events:**
```
event: connected
data: {"status": "connected"}

event: chunk
data: {"chunk": "Hello"}

event: done
data: {"title": "...", "response": "..."}
```

### POST /upload (Protected - File Upload)

```bash
curl -X POST http://127.0.0.1:5000/upload \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@image.jpg" \
  -F "type=image/jpeg"
```

**Response:**
```json
{
  "success": true,
  "key": "chat/files/...",
  "url": "/files/chat%2Ffiles%2F..."
}
```

### GET /file/:key (Protected)

```bash
curl http://127.0.0.1:5000/file/chat/files/abc123.jpg \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "url": "https://files.use.ai/files/chat/files/abc123.jpg"
}
```

---

## Available Models

| ID | Name                        | Display Name           |
|----|-----------------------------|------------------------|
| 1  | gateway-gpt-5-4            | GPT-5.4 (Latest)      |
| 2  | gateway-gpt-5-3            | GPT-5.3               |
| 3  | gateway-gpt-5-1            | GPT-5.1               |
| 4  | gateway-gpt-5              | GPT-5                 |
| 5  | gateway-gpt-4o             | GPT-4o                |
| 6  | gateway-gpt-4o-mini        | GPT-4o Mini           |
| 7  | gateway-grok-4             | Grok-4 (xAI)          |
| 8  | gateway-claude-sonnet-4-6  | Claude Sonnet 4.6     |
| 9  | gateway-claude-opus-4-5    | Claude Opus 4.5       |
| 10 | gateway-claude-opus-4-1    | Claude Opus 4.1       |
| 11 | gateway-deepseek-r1        | DeepSeek R1           |
| 12 | gateway-gemini-3-1-pro     | Gemini 3.1 Pro        |
| 13 | gateway-gemini-3-pro       | Gemini 3 Pro          |
| 14 | gateway-gemini-2.5-flash   | Gemini 2.5 Flash      |
| 15 | gateway-qwen-3-max         | Qwen 3 Max            |
| 16 | gateway-llama-3-3-70b-versatile | Llama 3.3 70B     |
| 17 | gateway-deepinfra-kimi-k2  | Kimi K2               |