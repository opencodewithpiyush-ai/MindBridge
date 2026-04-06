# MindBridge API

A clean architecture REST API server that proxies chat requests to an external WebSocket service.

## Tech Stack

- **Framework**: Gin (Go web framework)
- **Architecture**: Clean Architecture
- **WebSocket**: gorilla/websocket

## Project Structure

```
MindBridge/
├── cmd/server/         # Entry point
├── config/            # Configuration
├── domain/            # Entities & interfaces
├── application/       # DTOs & use cases
├── infrastructure/   # Generators & repositories
├── presentation/      # HTTP handlers
└── utils/             # Utilities
```

## Running the Server

```bash
go run ./cmd/server
```

Server runs at `http://127.0.0.1:5000`

## API Endpoints

| Method | Endpoint        | Description                     |
|--------|-----------------|---------------------------------|
| GET    | `/`             | API info                        |
| GET    | `/models`       | List available models          |
| POST   | `/chat`         | Send chat message (sync)       |
| POST   | `/chat/stream` | Send chat message (SSE stream) |
| POST   | `/upload`      | Upload image/file              |
| GET    | `/file/:key`   | Get uploaded file URL         |

## API Usage

### GET /

```bash
curl http://127.0.0.1:5000/
```

### GET /models

```bash
curl http://127.0.0.1:5000/models
```

### POST /chat (Sync Response)

```bash
curl -X POST http://127.0.0.1:5000/chat \
  -H "Content-Type: application/json" \
  -d '{"query": "Hello", "model": "gateway-claude-opus-4-1"}'
```

### POST /chat/stream (Real-time Streaming)

```bash
curl -X POST http://127.0.0.1:5000/chat/stream \
  -H "Content-Type: application/json" \
  -H "Accept: text/event-stream" \
  -d '{"query": "Hello", "model": "gateway-claude-opus-4-1"}'
```

**Response Events:**
```
event: connected
data: {"status": "connected"}

event: chunk
data: {"chunk": "Hello"}

event: chunk
data: {"chunk": "Hello, how"}

event: done
data: {"title": "...", "response": "..."}
```

### POST /upload (File Upload)

```bash
curl -X POST http://127.0.0.1:5000/upload \
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

### GET /file/:key

```bash
curl http://127.0.0.1:5000/file/chat/files/abc123.jpg
```

**Response:**
```json
{
  "success": true,
  "url": "https://files.use.ai/files/chat/files/abc123.jpg"
}
```

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
| 16 | gateway-llama-3-3-70b-versatile | Llama 3.3 70B   |
| 17 | gateway-deepinfra-kimi-k2  | Kimi K2               |
