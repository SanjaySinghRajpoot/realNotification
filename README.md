# Real Notification

A scalable notification service built with Go that handles SMS, Email, and In-App notifications using event-driven architecture with Apache Kafka.

![systemdesignDiagram](https://github.com/SanjaySinghRajpoot/realNotification/assets/67458417/5d3d5761-e05d-465c-b8b8-b7456d92003a)

## Table of Contents
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Architecture Overview](#architecture-overview)
- [Project Structure](#project-structure)
- [How It Works](#how-it-works)
- [API Endpoints](#api-endpoints)
- [Getting Started](#getting-started)
- [Monitoring](#monitoring)
- [Load Testing](#load-testing)
- [Contributing](#contributing)

## Features

- **Multi-channel notifications**: Supports SMS, Email, and In-App (push) notifications
- **Event-driven architecture**: Uses Apache Kafka for async message processing
- **Database sharding**: Notifications are distributed across two Postgres instances based on user ID (even/odd split)
- **Deduplication**: Redis cache prevents duplicate notifications within a 3-hour window
- **Rate limiting**: IP-based rate limiter (30 requests/minute) to prevent abuse
- **Retry mechanism**: CRON job runs every 10 seconds to retry failed notifications
- **Load balancing**: Custom round-robin load balancer with health checks
- **Monitoring**: Prometheus metrics + Grafana dashboards
- **Horizontal scaling**: Docker Compose setup runs 3 instances of the main app

## Tech Stack

| Component | Technology |
|-----------|------------|
| API Framework | Gin (Go) |
| Message Queue | Apache Kafka |
| Database | PostgreSQL (sharded) |
| Cache | Redis |
| Monitoring | Prometheus + Grafana |
| Containerization | Docker |
| Load Testing | k6 |

## Project Structure

```
realNotification/
├── main.go                 # Entry point - sets up routes, kafka, redis, cron
├── config/
│   ├── db.go              # PostgreSQL connection (2 sharded DBs)
│   └── cors_conf.go       # CORS middleware config
├── controller/
│   └── user.go            # Notification endpoint handler
├── models/
│   └── user.go            # Data structures (Notification, Payload, etc.)
├── routes/
│   └── user.go            # Route definitions
├── utils/
│   ├── utils.go           # Kafka producer, Redis ops, DB routing
│   └── constant.go        # Notification type constants (sms/email/inapp)
├── middleware/
│   ├── limiter.go         # IP-based rate limiter
│   └── ginPrometheus.go   # Prometheus metrics middleware
├── consumer/
│   └── main.go            # Kafka consumer - routes to notification services
├── services/
│   ├── sms/main.go        # SMS microservice
│   ├── mail/main.go       # Email microservice
│   └── inapp/main.go      # In-App notification microservice
├── loadbalancer/
│   └── main.go            # Round-robin LB with health checks
└── monitoring/            # Prometheus + Grafana config
```

## How It Works

### 1. Request Flow (Sending a Notification)

```
POST /user/notification
{
  "user_id": [1, 2, 3],
  "type": "sms",           // sms | email | inapp
  "description": "Hello!"
}
```

**Step-by-step:**

1. **Deduplication Check** (`controller/user.go`)
   - Generates a unique key: `{user_id}+{description}`
   - Checks Redis if this notification was sent in the last 3 hours
   - If exists, returns early with "already sent" message

2. **Cache the notification** (`utils/utils.go`)
   - Sets Redis key with 3-hour TTL to prevent duplicates

3. **Database Sharding** (`utils/utils.go`)
   - Routes to DB based on `user_id % 2`:
     - Even user IDs → `DB` (port 5432)
     - Odd user IDs → `DB1` (port 5433)
   - Saves notification with `state: false` (pending)

4. **Publish to Kafka** (`utils/utils.go`)
   - Serializes notification to JSON
   - Publishes to topic based on type (`sms`, `email`, or `inapp`)
   - Uses synchronous delivery with confirmation

### 2. Consumer Processing (`consumer/main.go`)

The Kafka consumer runs as a separate process:

1. Subscribes to topics: `sms`, `email`, `inapp`
2. On message received:
   - Deserializes the notification
   - Routes to appropriate service based on topic:
     - `sms` → `http://localhost:8082/sms`
     - `email` → `http://localhost:8083/mail`
     - `inapp` → `http://localhost:8084/inapp`

### 3. Service Processing (`services/*/main.go`)

Each service:
1. Receives the notification payload
2. (Would send actual notification - placeholder in code)
3. Updates notification `state` to `true` in the correct sharded DB
4. Returns success response

### 4. Retry Mechanism (`utils/utils.go`)

CRON job runs every 10 seconds:
```go
cronJob.AddFunc("@every 10s", func() {
    utils.CheckForNotificationState()
})
```

- Queries all notifications where `state = false`
- Re-publishes them to Kafka
- Ensures eventual delivery of failed notifications

### 5. Rate Limiting (`middleware/limiter.go`)

- Tracks request count per IP in Redis
- Limit: 30 requests per minute
- Decrements count after 1 minute using `time.AfterFunc`
- Returns `429 Too Many Requests` when exceeded

### 6. Load Balancing (`loadbalancer/main.go`)

Custom implementation:
- Round-robin algorithm across 3 app instances
- Health checks every 2 seconds (TCP dial)
- Auto-retry on failure (up to 3 attempts)
- Marks backends as down after consecutive failures

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Health check / welcome message |
| POST | `/user/notification` | Send notification to users |
| GET | `/metrics` | Prometheus metrics |
| GET | `/swagger/*` | Swagger API docs |

### Request/Response Examples

**Send Notification:**
```bash
curl -X POST http://localhost:8081/user/notification \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": [1, 2],
    "type": "sms",
    "description": "Your OTP is 123456"
  }'
```

**Response:**
```json
"Message Delivered Successfully"
```

## Getting Started

### Prerequisites
- Docker & Docker Compose
- Go 1.21+ (for local development)

### Running with Docker

1. Clone the repo:
```bash
git clone https://github.com/SanjaySinghRajpoot/realNotification.git
cd realNotification
```

2. Create `.env` file:
```bash
PASSWORD=12345678
DOCKER_HOST_IP=127.0.0.1
```

3. Start all services:
```bash
docker-compose up --build
```

This spins up:
- 3 main app instances (8079, 8080, 8081)
- 2 PostgreSQL instances (5432, 5433)
- Redis (6379)
- Kafka + Zookeeper (9092, 2181)

4. Start the consumer (separate terminal):
```bash
cd consumer
go run main.go
```

5. Start notification services:
```bash
cd services/sms && go run main.go    # port 8082
cd services/mail && go run main.go   # port 8083
cd services/inapp && go run main.go  # port 8084
```

6. (Optional) Start load balancer:
```bash
cd loadbalancer && go run main.go    # port 3030
```

## Monitoring

### Setup Prometheus & Grafana

1. Update `monitoring/prometheus/config.yml` with your host IP

2. Start monitoring stack:
```bash
cd monitoring
docker-compose up -d --build
```

3. Access dashboards:
   - Prometheus: http://localhost:9090
   - Grafana: http://localhost:3000 (admin/admin)

### Available Metrics

The app exposes standard HTTP metrics via `/metrics`:
- `gin_requests_total` - Total requests by status, method, handler
- `gin_request_duration_seconds` - Request latency histogram
- `gin_request_size_bytes` - Request payload sizes
- `gin_response_size_bytes` - Response payload sizes

## Load Testing

Using k6:

```bash
cd loadtesting
k6 run loadtest.js
```

Test configuration:
- Ramp up to 10 virtual users over 10s
- Hold for 5s
- Ramp down over 15s
- Threshold: 95th percentile latency < 500ms

## Database Schema

```sql
CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    description VARCHAR(1000),
    type VARCHAR(50),          -- 'sms', 'email', 'inapp'
    state BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string (shard 1) | - |
| `DATABASE_URL_1` | PostgreSQL connection string (shard 2) | - |
| `PASSWORD` | Redis password | - |
| `DOCKER_HOST_IP` | Host IP for Kafka | 127.0.0.1 |

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/awesome`)
3. Commit your changes (`git commit -m 'Add awesome feature'`)
4. Push to the branch (`git push origin feature/awesome`)
5. Open a Pull Request

## License

MIT License - see [LICENSE](LICENSE) for details.

## Related Projects

- [loadbalancer](https://github.com/SanjaySinghRajpoot/loadbalancer) - The custom load balancer used in this project

