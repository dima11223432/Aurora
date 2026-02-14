Отлично! Вот README.md только с миграциями и запуском через make, без Docker:

````markdown
# Название проекта

Микросервисная архитектура с API Gateway и SSO (Single Sign-On).

## 📋 Требования

- Go 1.21+
- PostgreSQL 14+
- Make
- [golang-migrate](https://github.com/golang-migrate/migrate) для миграций

```bash
# Установка migrate (если не установлен)
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```
````

## 📦 Установка зависимостей

### Важно! Устанавливаем правильную версию protos

```bash
# В корне проекта
go get github.com/dima11223432/protos@v0.0.10
go mod download

# Для SSO сервиса
cd SSO && go get github.com/dima11223432/protos@v0.0.10 && go mod download && cd ..

# Для API Gateway
cd Api_Gateway && go get github.com/dima11223432/protos@v0.0.10 && go mod download && cd ..
```

Или одной командой через make:

```bash
make deps
```

## 🗄 Настройка PostgreSQL

### Создание баз данных

```bash
# Подключитесь к PostgreSQL
psql -U postgres -h localhost -p 5432

# Создайте базы данных
CREATE DATABASE users;
CREATE DATABASE aurora;

# Выйдите из psql
\q
```

Или через make:

```bash
make db-create
```

### Проверка создания баз

```bash
psql -U postgres -h localhost -p 5432 -c "\l" | grep -E "users|ens|aurora"
```

## 📊 Миграции

### Структура миграций

Проект использует [golang-migrate](https://github.com/golang-migrate/migrate) для управления миграциями:

- `SSO/migrations/` - миграции для базы `users` (SSO сервис)
- `Api_Gateway/migrations/` - миграции для базы `aurora` (API Gateway)

### Применение миграций

```bash
# Применить все миграции для всех баз
make migrate-all

# Или по отдельности:

# SSO миграции (база users)
make migrate-sso

# API Gateway миграции (база ens)
make migrate-gateway

# Aurora миграции (база aurora)
make migrate-aurora
```

### Управление миграциями

```bash
# Создать новую миграцию для SSO
make migrate-sso-create name=название_миграции

# Создать новую миграцию для Gateway
make migrate-gateway-create name=название_миграции

# Создать новую миграцию для Aurora
make migrate-aurora-create name=название_миграции

# Откатить последнюю миграцию SSO
make migrate-sso-down

# Откатить последнюю миграцию Gateway
make migrate-gateway-down

# Проверить статус миграций
make migrate-sso-version
make migrate-gateway-version
```

### Если миграции не применяются

Если вы видите сообщение "no migrations to apply", проверьте:

```bash
# 1. Проверьте наличие файлов миграций
ls -la Api_Gateway/migrations/
ls -la SSO/migrations/

# 2. Убедитесь, что файлы имеют правильный формат:
#    000001_init.up.sql
#    000001_init.down.sql

# 3. Принудительно сбросьте версию (если нужно)
migrate -database "postgres://postgres:pass@localhost:5432/ens?sslmode=disable" \
  -path ./Api_Gateway/migrations force 0

# 4. Примените миграции заново
make migrate-gateway
```

## 🚀 Запуск сервисов

### Запуск SSO сервиса

```bash
# Из корня проекта
make run-sso

# Или вручную:
cd SSO
go run ./cmd/sso/main.go --config=./config/local.yaml
```

Сервис запустится на порту 44044 (gRPC).

### Запуск API Gateway

```bash
# Из корня проекта
make run-gateway

# Или вручную:
cd Api_Gateway
go run ./cmd/main.go --config=./config/local.yaml
```

Сервис запустится на порту 8080 (HTTP).

### Запуск всех сервисов одновременно

```bash
# Запуск в отдельных терминалах
make run-all

# Или вручную:
# Терминал 1:
cd SSO && go run ./cmd/sso/main.go --config=./config/local.yaml

# Терминал 2:
cd Api_Gateway && go run ./cmd/main.go --config=./config/local.yaml
```

## 🔧 Полный список Make команд

```bash
# Справка
make help

# Установка зависимостей
make deps              # установить все Go-зависимости (включая protos v0.0.10)

# Базы данных
make db-create         # создать все БД (users, ens, aurora)
make db-list           # показать список БД
make db-drop           # удалить все БД (осторожно!)

# Миграции
make migrate-all       # применить все миграции
make migrate-sso       # миграции для users DB (SSO)
make migrate-gateway   # миграции для ens DB (Gateway)
make migrate-aurora    # миграции для aurora DB
make migrate-sso-down  # откат последней миграции SSO
make migrate-gateway-down # откат последней миграции Gateway

# Создание миграций
make migrate-sso-create name=init      # создать миграцию для SSO
make migrate-gateway-create name=init  # создать миграцию для Gateway
make migrate-aurora-create name=init   # создать миграцию для Aurora

# Запуск сервисов
make run-sso          # запустить SSO
make run-gateway      # запустить API Gateway
make run-all          # запустить все сервисы

# Тестирование
make test             # все тесты
make test-sso         # тесты SSO
make test-gateway     # тесты Gateway

# Утилиты
make fmt              # форматирование кода
make lint             # запуск линтера
make clean            # очистка
```

## 📝 Конфигурация

### SSO сервис (`SSO/config/local.yaml`)

```yaml
env: "local"
database_url: "postgres://postgres:pass@localhost:5432/users?sslmode=disable"
grpc:
  port: 44044
  timeout: 1h
```

### API Gateway (`Api_Gateway/config/local.yaml`)

```yaml
server:
  port: 8080
database:
  ens:
    url: "postgres://postgres:pass@localhost:5432/ens?sslmode=disable"
sso:
  addr: "localhost:44044"
```

## ✅ Проверка работоспособности

```bash
# 1. Запустите SSO
make run-sso
# Должны увидеть: "starting gRPC server" и "gRPC server listening"

# 2. В другом терминале запустите Gateway
make run-gateway

# 3. Проверьте health endpoint
curl http://localhost:8080/health
# Ожидаемый ответ: {"status":"ok"} или подобное
```

## 🐛 Возможные проблемы и решения

### Проблема: "no migrations to apply"

**Решение:** Проверьте формат файлов миграций (должны быть 000001\_\*.up.sql и .down.sql)

### Проблема: "cannot find package github.com/dima11223432/protos"

**Решение:** Убедитесь, что установлена правильная версия:

```bash
go get github.com/dima11223432/protos@v0.0.10
```

### Проблема: "connection refused" к PostgreSQL

**Решение:** Проверьте, запущен ли PostgreSQL:

```bash
pg_isready -U postgres -h localhost -p 5432
```

## 📂 Структура проекта

```
.
├── Api_Gateway/          # API Gateway микросервис
│   ├── cmd/             # точка входа
│   ├── config/          # конфигурация
│   ├── internal/        # внутренний код
│   └── migrations/      # миграции для ens БД
├── SSO/                  # SSO микросервис
│   ├── cmd/             # точка входа
│   ├── config/          # конфигурация
│   ├── internal/        # внутренний код
│   └── migrations/      # миграции для users БД
├── migrations/           # общие миграции
│   └── aurora/          # миграции для aurora БД
├── makefile             # главный makefile
└── README.md
```

## 📌 Важные замечания

1. **Версия protos**: строго используйте `v0.0.10`
2. **PostgreSQL**: должен быть запущен до применения миграций
3. **Порты**:
   - SSO: 44044 (gRPC)
   - API Gateway: 8080 (HTTP)
   - PostgreSQL: 5432
4. **Очередность запуска**: сначала SSO, потом API Gateway

---

> **Далее**: Для запуска через Docker смотрите раздел "Docker" в полной версии документации.

```

Этот README.md фокусируется только на:
- ✅ Установке зависимостей (с правильной версией protos v0.0.10)
- ✅ Создании баз данных
- ✅ Миграциях
- ✅ Запуске через make

Docker часть пока не включена. Устраивает?
```
