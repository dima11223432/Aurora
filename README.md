# Aurora

Microservices architecture with API Gateway, SSO (Single Sign-On), Telegram channel parsing, and analytics for cryptocurrency trading signals.

## Architecture

```
┌─────────────┐     ┌─────────────┐
│   Frontend  │────▶│ API Gateway │
│   (React)   │     │   (Go)      │
└─────────────┘     └──────┬──────┘
                           │
         ┌─────────────────┼─────────────────┐
         │                 │                 │
         ▼                 ▼                 ▼
    ┌─────────┐      ┌──────────┐     ┌────────────┐
    │   SSO   │      │ Recommen-│     │   Cache    │
    │ (Go)    │      │ dation   │     │  Service   │
    └─────────┘      │ (Go)     │     │   (Go)     │
                     └──────────┘     └─────┬──────┘
                                             │
                  ┌──────────────────────────┼──────────────┐
                  ▼                          ▼              ▼
            ┌──────────┐            ┌────────────┐   ┌──────────┐
            │  Parser  │───────────▶│   Kafka    │◀──│Analytics │
            │ (Python) │            │            │   │(Python)  │
            └──────────┘            └────────────┘   └──────────┘

         ┌─────────────┐     ┌─────────────┐
         │ PostgreSQL  │     │    Redis    │
         │ (aurora,    │     │             │
         │  parser)    │     │             │
         └─────────────┘     └─────────────┘
```

## Services

| Service        | Language   | Port  | Description           |
| -------------- | ---------- | ----- | --------------------- |
| Frontend       | React/Vite | 5173  | Web UI                |
| API Gateway    | Go         | 8081  | Main HTTP API         |
| SSO            | Go         | 44044 | Authentication (gRPC) |
| Recommendation | Go         | 44040 | Trading signals       |
| Cache          | Go         | 44045 | Redis caching         |
| Parser         | Python     | 44042 | Telegram parsing      |
| Analytics      | Python     | 44047 | Data analytics        |

## Tech Stack

- **Frontend**: React 19, Vite, Tailwind CSS, React Router, TON Connect, TradingView
- **Backend**: Go 1.21+, Python 3.11+
- **Database**: PostgreSQL 15, Redis 7
- **Message Queue**: Kafka 7.4.0
- **Container**: Docker Compose

## Quick Start

### Prerequisites

- Go 1.21+
- Python 3.11+
- Docker & Docker Compose
- Make

### Local Development

```bash
git clone https://gitlab.informatics.ru/2025-2026/ydex/s103d/final-project-t1918.git
cd final-project-t1918
```

#### 1. Start Infrastructure

```bash
docker compose up -d postgres postgres-parser redis-master zookeeper kafka
```

#### 2. Install Dependencies

```bash
cd backend/SSO
go get github.com/dima11223432/protos@v0.0.10
go mod download

cd ../Api_Gateway
go get github.com/dima11223432/protos@v0.0.10
go mod download

cd ../../frontend
npm install
```

#### 3. Run Migrations

```bash
cd backend

make db-create

make migrate-all
```

### How to Build project?

```bash
docker compose up -d

docker compose logs -f
```

## Environment Variables

## API Endpoints

- **Frontend**: http://localhost:5173
- **API Gateway**: http://localhost:8081
- **SSO gRPC**: localhost:44044
- **Recommendation gRPC**: localhost:44040
- **Cache Service gRPC**: localhost:44045

## Project Structure

```
Aurora/
├── backend/
│   ├── SSO/              # SSO authentication service
│   ├── Api_Gateway/      # Main API gateway
│   ├── Recommendation_Service/  # Trading recommendations
│   ├── Cache_Service/    # Redis caching with Kafka
│   ├── Parsing_Service/  # Telegram channel parser (Python)
│   ├── Analitic_Service/ # Analytics (Python)
│   └── makefile
├── frontend/             # React frontend
├── docker-compose.yaml
├── init.sql
└── README.md
```

## License

Internal project - Informatics 2025-2026

