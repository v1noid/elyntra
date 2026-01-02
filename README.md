# Project elyntra Ngrok Clone

A lightweight ngrok-like tunneling solution built with **Elysia** (TypeScript/Node.js) for the server and **Go** for the client CLI.

## Architecture

- **Server**: Elysia-based HTTP server that manages client registration, tunnel creation, and request forwarding
- **Client**: Go CLI application for users to create and manage tunnels to their local services

## Features

- Client registration and authentication with token-based security
- Create tunnels from local services to public URLs
- List active tunnels
- Close/delete tunnels
- Request forwarding from public URLs to local services
- Simple configuration storage

## Project Structure

```
ngrok-clone/
├── server/                 # Elysia server
│   ├── src/
│   │   └── index.ts       # Main server code
│   ├── package.json       # Server dependencies
│   └── tsconfig.json      # TypeScript config
├── main.go                # Go client CLI
├── go.mod                 # Go module definition
└── README.md              # This file
```

## Setup

### Server (Elysia)

1. Install dependencies:

```bash
cd server
bun install
```

2. Start the development server:

```bash
bun run dev
```

Or run production:

```bash
bun run start
```

The server will run on `http://localhost:3000`

### Client (Go)

1. Build the CLI:

```bash
go build -o ngrok-clone main.go
```

2. Or run directly:

```bash
go run main.go <command> [options]
```

## Usage

### 1. Register a Client

First, register your client with the server:

```bash
./ngrok-clone register -server http://localhost:3000
```

This creates a configuration file at `~/.ngrok-clone/config.json` with your client ID and authentication token.

### 2. Create a Tunnel

Create a tunnel to expose a local service:

```bash
./ngrok-clone tunnel -local http://localhost:8080
```

This will output:

- **Tunnel ID**: Unique identifier for this tunnel
- **Public URL**: The URL to share publicly (e.g., `https://abc123.tunnel.local`)
- **Local URL**: Your local service URL

### 3. List Active Tunnels

View all active tunnels:

```bash
./ngrok-clone list
```

### 4. Close a Tunnel

Close a tunnel when no longer needed:

```bash
./ngrok-clone close -id <tunnelId>
```

## API Endpoints

### Server API

- **POST** `/auth/register` - Register a new client
- **POST** `/tunnels/create` - Create a tunnel
- **GET** `/tunnels/:tunnelId` - Get tunnel information
- **GET** `/client/:clientId/tunnels` - List client's tunnels
- **DELETE** `/tunnels/:tunnelId` - Close a tunnel
- **ALL** `/t/:tunnelId/*` - Forward requests to local service
- **GET** `/health` - Health check

## Configuration

Client configuration is stored at `~/.ngrok-clone/config.json`:

```json
{
  "clientId": "your-client-id",
  "token": "your-auth-token",
  "server": "http://localhost:3000"
}
```

## How It Works

1. **Registration**: Client registers with the server and receives a unique ID and token
2. **Tunnel Creation**: Client sends a tunnel creation request with local URL and receives a public URL
3. **Request Forwarding**: Requests to the public URL (`/t/:tunnelId/*`) are forwarded to the local service
4. **Tunnel Management**: Clients can list and close tunnels as needed

## Future Enhancements

- WebSocket support for real-time tunnel management
- Multiple tunnel support per client
- Bandwidth monitoring and logging
- Custom domain support
- Traffic inspection dashboard
- TLS/SSL certificate management
- Rate limiting and authentication options
- Database for persistent tunnel storage
