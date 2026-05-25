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

# Платформа для агрегации и анализа новостей фондового рынка и криптовалют

---

## Актуальность темы

Тема Fin-tech актуальна, потому что она требует большой компетенции в разработке, архитектуре и безопасности, а также уверенной работы со сторонними API.

---

## Технологическая особенность

Сложно создать проект, отличающийся от остальных проектов в техническом плане, но наш проект будет создаваться по большинству стандартов современной коммерческой разработки, которая позволит сервису быть высоконагруженным и прод пригодным.

---

## Стек

```
• Golang
• Python3
• Django
• gRPC
• Redis
• Kafka
• Postgres
• React.js
• TailwindCss
```

---

## Трудозатраты и роли

Команде T-1918 понадобится около 300 часов на реализацию проекта "Aurora".

| Участник | Роль | Часы |
|----------|------|------|
| Сафронов Дмитрий | Тимлид | 200 часов |
| Иосифов Николай | Старший Разработчик Гибрид | 200 часов |
| Бурдюков Савелий | Middle разработчик | 180 часов |
| Хвастюк Ян | Middle Разработчик | 150 часов |
| Жильцова Екатерина | Junior+ Разработчик | 120 часов |
| Манчуленко Василий | Junior Разработчик | 80 часов |

---

## Актуализация к курсу

Я считаю, что "Aurora" должен считаться итоговым проектом, потому что он задействует область Fin-tech, являющуюся возможно сложнейшей в реализации и безопасности областью, которая требует немалого багажа знаний не только тимлида, но и всей команды. В проекте используются самые новые технологии в промышленной разработке такие как: gRPC, Redis, Kafka, Golang, Python. Исходя из вышеперечисленного, Наш проект должен считаться итоговым.

---

# Техническое задание к проекту

## Проблема

"Aurora" решает проблему поиска и обнаружения паттернов корреляции цены активов: монет на основе TON, а также акций на фондовом рынке.

## Цель

Разработка платформы для агрегации и анализа новостей фондового рынка и криптовалют (экосистема TON) с возможностью принятия торговых решений и совершения покупок через TON Wallet. Проект разбит на три спринта: Core, Intelligence, Trading.

---

## Тематика

```
#Fin-tech, #Telegram, #TON, #TON_Wallet, #T-investment
```

---

## Минимальные функциональные требования

1. Пользователь подключает TON кошелек к платформе
2. Система получает wallet address и привязывает его к user_id
3. Пользователь выбирает Telegram-каналы для отслеживания
4. Пользователь может включать/отключать каналы в своем профиле
5. Система в фоне парсит выбранные Telegram-каналы (только последние 24-48 часов). Парс будет производиться через Python библиотеку Telethon
6. Система формирует ленту новостей только по фондовому рынку
7. Посты в ленте сортируются по времени (новые сверху)
8. Каждая новость показывает дату публикации
9. Каждая новость содержит ссылку на оригинальный пост в Telegram
10. Текст новостей отправляется в очередь на LLM-анализ
11. Анализ выполняется асинхронно (не блокирует интерфейс)
12. Результаты анализа кешируются для одинаковых новостей
13. CNN LSTM модели определяют тональность новости (bullish/bearish/neutral)
14. Результат анализа возвращается в формате JSON

---

## Производительность

- 10,000 DAU (Daily Active Users) можно тестировать через Apache JMeter
- Максимальное время загрузки всех новостей (пак по 15 штук) — 3 секунды (Можно проводить проверку на фронтенде)

---

## Метрики пользователей

- Большинство пользователей в режиме Read-only

---

## География

- Аудитория: СНГ

---

## Варианты расширения функционала

1. Пользователь может покупать монеты TON через интерфейс (без автотрейдинга)
2. Покупка требует подтверждения в TON кошельке
3. История транзакций сохраняется в профиле
4. Система отображает графики цен акций по тикерам
5. Система отображает графики цен монет TON
6. На графиках отображаются метки новостей по датам
7. Система отправляет уведомления на основе аналитики (например, "Акция X может вырасти на Y%")

---