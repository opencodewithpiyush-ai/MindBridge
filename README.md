# MindBridge API

A clean architecture REST API server with authentication, AI chat, and image editing capabilities.

## Technology Stack

- **Framework**: Gin (Go web framework)
- **Architecture**: Clean Architecture
- **Database**: MongoDB
- **Cache/Session**: Redis
- **Authentication**: JWT + Bcrypt
- **WebSocket**: gorilla/websocket
- **Streaming**: Server-Sent Events (SSE)

## Project Structure

```
MindBridge/
├── cmd/server/                 # Entry point
├── config/                     # Configuration & environment
├── domain/
│   ├── entities/               # Domain entities (User, Chat)
│   └── repositories/           # Interface definitions
├── application/
│   ├── dto/                    # Data Transfer Objects
│   └── usecases/               # Business logic
├── infrastructure/
│   ├── repositories/           # Implementations (MongoDB, WebSocket, JWT, Redis, File)
│   └── generators/             # ID & Email generators
├── presentation/
│   └── handlers/               # HTTP handlers & middleware
├── utils/                      # Validators & logging
└── web/                        # React + TypeScript frontend
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

## Running the Frontend

```bash
cd web
npm install
npm run dev
```

Frontend runs at `http://localhost:5173`

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
| POST | `/chat/stream-raw` | Send chat message (SSE stream, all events incl. tool calls) |
| POST | `/auth/logout` | Logout and invalidate session |
| POST | `/upload` | Upload image/file |

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
    }
  ]
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

**Error Response (401):**
```json
{
  "success": false,
  "error": "invalid email or password"
}
```

### POST /auth/logout (Protected)

```bash
curl -X POST http://127.0.0.1:5000/auth/logout \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Logout successful"
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

### POST /chat/stream-raw (Protected - Full SSE Events)

Streams all raw WebSocket events including tool calls (image generation, editing, etc.).

```bash
curl -X POST http://127.0.0.1:5000/chat/stream-raw \
  -H "Content-Type: application/json" \
  -H "Accept: text/event-stream" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"query": "Edit this image", "model": "gateway-claude-sonnet-4-6", "files": [{"name": "img.jpg", "type": "image/jpeg", "url": "https://files.use.ai/files/..."}]}'
```

**Response Events:**
```
event: connected
data: {"status": "connected"}

event: chunk
data: {"type": "stream-start", ...}

event: chunk
data: {"type": "data-chat-title-update", "data": {"title": "..."}}

event: chunk
data: {"chunk": {"type": "text-delta", "delta": "Hello"}}

event: chunk
data: {"chunk": {"type": "tool-input-available", "toolName": "image-google", "input": {...}}}

event: chunk
data: {"chunk": {"type": "tool-image-google", "output": {"images": [{"url": "..."}]}}}

event: done
data: {"title": "...", "response": "..."}
```

### POST /upload (Protected - File Upload)

```bash
curl -X POST http://127.0.0.1:5000/upload \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@image.jpg" \
  -F "name=image.jpg" \
  -F "type=image/jpeg"
```

**Response:**
```json
{
  "success": true,
  "key": "chat/files/b738979b-...-image.jpg",
  "url": "/files/chat%2Ffiles%2Fb738979b-...-image.jpg"
}
```

---

## Available Models

| ID | Name | Display Name |
|----|------|--------------|
| 1 | gateway-gpt-5-5 | GPT-5.5 (Latest) |
| 2 | gateway-gpt-5-4 | GPT-5.4 |
| 3 | gateway-gpt-5-3 | GPT-5.3 |
| 4 | gateway-gpt-5-1 | GPT-5.1 |
| 5 | gateway-gpt-5 | GPT-5 |
| 6 | gateway-gpt-4o | GPT-4o |
| 7 | gateway-gpt-4o-mini | GPT-4o Mini |
| 8 | gateway-grok-4 | Grok-4 (xAI) |
| 9 | gateway-claude-sonnet-4-6 | Claude Sonnet 4.6 |
| 10 | gateway-claude-opus-4-5 | Claude Opus 4.5 |
| 11 | gateway-claude-opus-4-1 | Claude Opus 4.1 |
| 12 | gateway-deepseek-v4-pro | DeepSeek V4 Pro |
| 13 | gateway-deepseek-v4-flash | DeepSeek V4 Flash |
| 14 | gateway-deepseek-r1 | DeepSeek R1 |
| 15 | gateway-gemini-3-1-pro | Gemini 3.1 Pro |
| 16 | gateway-gemini-3-pro | Gemini 3 Pro |
| 17 | gateway-gemini-2.5-flash | Gemini 2.5 Flash |
| 18 | gateway-qwen-3-max | Qwen 3 Max |
| 19 | gateway-llama-3-3-70b-versatile | Llama 3.3 70B |
| 20 | gateway-deepinfra-kimi-k2 | Kimi K2 |

## Features

- **20 AI Models**: Single gateway for OpenAI, Anthropic, Google, xAI, DeepSeek, and more
- **JWT Authentication**: Register, login, logout with session management via Redis
- **Real-time Streaming**: SSE-based streaming for chat responses
- **Image Editing**: Upload images and edit them via AI (shirt change, glasses, etc.)
- **File Upload**: Support for images and files with preview
- **Validation**: Strict input validation for registration

## Support The Developer

Give it a ⭐. If You Found This Useful.
