# Go Redis-like Key-Value Store

A lightweight, Redis-inspired key-value store built in go with persistence, TTL support, and distributed deployment capabilities.

## Features

### Core Functionality

- **In-memory key-value storage** with thread safe operations
- **Time-To-Live (TTL) system** - automatic key expiration
- **Persistence** via Append-Only File (AOF) logging
- **Crash recovery** by replaying logged commands
- **Multi-client support** with concurrent connections

### Supported Commands

- `SET key value` - Store key-value pair (6-hour default TTL)
- `GET key` - Retrieve value (returns error if expired)
- `DEL key` - Delete key and return deleted value
- `SETEX key seconds value` - Set key with custom expiration time
- `exit` - Terminate client session

### Distributed Deployment

- **Docker containerization** for easy deployment
- **Multi-instance setup** with Docker Compose
- **Load balancing** using nginx
- **Horizontal scaling** support

## 🛠️ Technical Architecture

### Storage Engine

- **sync.Map** for thread-safe concurrent access
- **TimedValue struct** wrapping values with expiration metadata
- **Passive expiration** checking on key access

### Networking

- **TCP server** on port 6379 (Redis-compatible)
- **Goroutine-per-connection** model for concurrency
- **Custom text protocol** with space-delimited commands

### Persistence

- **Write-Ahead Logging (WAL)** for durability
- **AOF (Append-Only File)** format for command logging
- **Recovery mechanism** replays commands on startup

## 📋 Prerequisites

- Go 1.19 or higher
- Docker and Docker Compose (for distributed setup)

## 🚀 Quick Start

### Single Instance

1. **Clone and build:**

```bash
git clone <repository-url>
cd go-redis-kvstore
go build -o redis-server .
```

2. **Run the server:**

```bash
./redis-server
```

3. **Connect with telnet:**

```bash
telnet localhost 6379
```

4. **Try some commands:**

```
SET mykey hello
GET mykey
SETEX tempkey 10 temporary
DEL mykey
exit
```

### Distributed Setup (Docker)

1. **Build and run multiple instances:**

```bash
docker-compose up --build
```

2. **Connect through load balancer:**

```bash
telnet localhost 6379
```

The load balancer will distribute your connections across multiple Redis instances.

## 📁 Project Structure

```
.
├── main.go              # Main server implementation
├── Dockerfile           # Container configuration
├── docker-compose.yml   # Multi-instance setup
├── nginx.conf           # Load balancer configuration
├── log.txt             # AOF persistence file (auto-generated)
└── README.md           # This file
```

## 🏗️ Architecture Overview

```
Client Connections
       ↓
  Load Balancer (nginx)
       ↓
┌─────────────────────────┐
│  Redis Instance Pool    │
├─────────┬─────────┬─────┤
│ godis1  │ godis2  │ ... │
│ :6379   │ :6379   │     │
└─────────┴─────────┴─────┘
       ↓
  Persistent Storage
    (AOF Files)
```

## ⚙️ Configuration

### Default Settings

- **Default TTL:** 6 hours
- **Server Port:** 6379
- **AOF File:** `log.txt`
- **Load Balancer:** Round-robin distribution

### Docker Compose Services

- **godis1, godis2, godis3:** Redis instances
- **nginx:** Load balancer on port 6379
