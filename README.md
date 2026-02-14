                                  
````markdown

                                     
#     /\                              
#    /  \  _   _ _ __ ___  _ __ __ _  
#   / /\ \| | | | '__/ _ \| '__/ _` | 
#  / ____ \ |_| | | | (_) | | | (_| | 
# /_/    \_\__,_|_|  \___/|_|  \__,_| 
                                     
                                     

Микросервисная архитектура с API Gateway и SSO (Single Sign-On).

## 📋 Требования

- Go 1.21+
- PostgreSQL 14+
- Make
- [golang-migrate](https://github.com/golang-migrate/migrate) для миграций

```bash
git clone https://gitlab.informatics.ru/2025-2026/ydex/s103d/final-project-t1918.git
```

```bash
# Установка migrate (если не установлен)
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```
````

## 📦 Установка зависимостей

### Важно! Устанавливаем правильную версию protos

```bash

# Для SSO сервиса
cd SSO && go get github.com/dima11223432/protos@v0.0.10 && go mod download && cd ..

# Для API Gateway
cd Api_Gateway && go get github.com/dima11223432/protos@v0.0.10 && go mod download && cd ..
```

## 🗄 Настройка PostgreSQL

### Создание баз данных

```bash
make db-create
```

### Проверка создания баз

```bash
psql -U postgres -h localhost -p 5432 -c "\l" | grep -E "users|aurora"
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

# API Gateway миграции (база aurora)
make migrate-gateway

# Aurora миграции (база aurora)
make migrate-aurora
```

### Если миграции не применяются

Если вы видите сообщение "no migrations to apply", проверьте:

```bash
# 1. Проверьте наличие файлов миграций
ls -la Api_Gateway/migrations/
ls -la SSO/migrations/

# 2. Убедитесь, что файлы имеют правильный формат:
#    1_init.up.sql
#    1_init.down.sql

# 3. Принудительно сбросьте версию (если нужно)
migrate -database "postgres://postgres:pass@localhost:5432/aurora?sslmode=disable" \
  -path ./Api_Gateway/migrations force 0

# 4. Примените миграции заново
make migrate-gateway
```

## 🚀 Запуск сервисов

### Запуск SSO сервиса

```bash
# Из backend
make run-sso

# Или вручную:
cd SSO
go run ./cmd/sso/main.go --config=./config/local.yaml
```

Сервис запустится на порту 44044 (gRPC).

### Запуск API Gateway

```bash
# из backed
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
