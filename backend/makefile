# Makefile для управления микросервисами

# ===== Переменные =====
# Цвета для вывода
GREEN := \033[0;32m
YELLOW := \033[1;33m
RED := \033[0;31m
NC := \033[0m # No Color

# Базы данных
DB_USER := postgres
DB_PASS := pass
DB_HOST := localhost
DB_PORT := 5432
DB_SSLMODE := disable

# URLs баз данных
USERS_DB_URL := postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/users?sslmode=$(DB_SSLMODE)
ENS_DB_URL := postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/ens?sslmode=$(DB_SSLMODE)
AURORA_DB_URL := postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/aurora?sslmode=$(DB_SSLMODE)

# Пути к миграциям
SSO_MIGRATIONS_PATH := ./SSO/migrations
GATEWAY_MIGRATIONS_PATH := ./Api_Gateway/migrations
AURORA_MIGRATIONS_PATH := ./migrations/aurora

# ===== Команды =====
.PHONY: help run-sso run-gateway run-all migrate-all migrate-sso migrate-gateway migrate-aurora \
        migrate-sso-up migrate-sso-down migrate-sso-force migrate-sso-create \
        migrate-gateway-up migrate-gateway-down migrate-gateway-force migrate-gateway-create \
        migrate-aurora-up migrate-aurora-down migrate-aurora-force migrate-aurora-create \
        db-create db-drop db-list docker-up docker-down test clean doc-gen doc-gen-parsing doc-gen-analitic

# ===== Помощь =====
help:
	@echo "$(GREEN)Доступные команды:$(NC)"
	@echo "$(YELLOW)Запуск сервисов:$(NC)"
	@echo "  make run-sso           - Запустить SSO сервис"
	@echo "  make run-gateway       - Запустить API Gateway"
	@echo "  make run-all           - Запустить все сервисы"
	@echo ""
	@echo "$(YELLOW)Миграции SSO:$(NC)"
	@echo "  make migrate-sso       - Применить все миграции SSO"
	@echo "  make migrate-sso-up    - Применить следующую миграцию SSO"
	@echo "  make migrate-sso-down  - Откатить последнюю миграцию SSO"
	@echo "  make migrate-sso-force version=1 - Принудительно установить версию"
	@echo "  make migrate-sso-create name=init - Создать новую миграцию"
	@echo ""
	@echo "$(YELLOW)Миграции Gateway:$(NC)"
	@echo "  make migrate-gateway   - Применить все миграции Gateway"
	@echo "  make migrate-gateway-up - Применить следующую миграцию Gateway"
	@echo "  make migrate-gateway-down - Откатить последнюю миграцию Gateway"
	@echo "  make migrate-gateway-force version=1 - Принудительно установить версию"
	@echo "  make migrate-gateway-create name=init - Создать новую миграцию"
	@echo ""
	@echo "$(YELLOW)Миграции Aurora:$(NC)"
	@echo "  make migrate-aurora    - Применить все миграции Aurora"
	@echo "  make migrate-aurora-up - Применить следующую миграцию Aurora"
	@echo "  make migrate-aurora-down - Откатить последнюю миграцию Aurora"
	@echo "  make migrate-aurora-force version=1 - Принудительно установить версию"
	@echo "  make migrate-aurora-create name=init - Создать новую миграцию"
	@echo ""
	@echo "$(YELLOW)Базы данных:$(NC)"
	@echo "  make db-list           - Список всех баз данных"
	@echo "  make db-create         - Создать все базы данных"
	@echo "  make db-drop           - Удалить все базы данных (осторожно!)"
	@echo ""
	@echo "$(YELLOW)Docker:$(NC)"
	@echo "  make docker-up         - Запустить все сервисы в Docker"
	@echo "  make docker-down       - Остановить все сервисы в Docker"
	@echo "  make docker-logs       - Просмотр логов Docker"
	@echo ""
	@echo "$(YELLOW)Документация (Python — Sphinx):$(NC)"
	@echo "  make doc-gen-parsing   - Сгенерировать HTML-документацию Parsing_Service"
	@echo "  make doc-gen-analitic  - Сгенерировать HTML-документацию Analitic_Service"
	@echo ""
	@echo "$(YELLOW)Документация (Go — golds):$(NC)"
	@echo "  make doc-gen-sso       - Сгенерировать документацию SSO"
	@echo "  make doc-gen-api       - Сгенерировать документацию API Gateway"
	@echo "  make doc-gen-recs      - Сгенерировать документацию Recommendation Service"
	@echo "  make doc-gen-cache     - Сгенерировать документацию Cache Service"
	@echo "  make doc-gen           - Сгенерировать документацию всех сервисов"
	@echo ""
	@echo "$(YELLOW)Прочее:$(NC)"
	@echo "  make test              - Запустить все тесты"
	@echo "  make clean             - Очистить бинарники"

# ===== Запуск сервисов =====
run-sso:
	@echo "$(GREEN)Запуск SSO сервиса...$(NC)"
	cd SSO && go run ./cmd/sso/main.go --config=./config/local.yaml

run-gateway:
	@echo "$(GREEN)Запуск API Gateway...$(NC)"
	cd Api_Gateway && go run ./cmd/API_Service/main.go --config=./config/local.yaml
run-recommendationService:
	@echo "$(GREEN)Запуск recommendationService...$(NC)"
	cd Recommendation_Service && go run ./cmd/recomendation_service/main.go --config=./config/local.yaml
run-cacheService:
	@echo "$(GREEN)Запуск cacheService...$(NC)"
	cd Cache_Service && go run ./cmd/cacheService/main.go --config=./config/local.yaml

run-all:
	@echo "$(GREEN)Запуск всех сервисов...$(NC)"
	@make -j 4 run-sso run-cacheService run-recommendationService run-gateway

# ===== Миграции SSO =====
migrate-sso:
	@echo "$(GREEN)Применение миграций SSO...$(NC)"
	cd SSO && go run ./cmd/migrator -db-dsn "$(USERS_DB_URL)" -migrations-path ./migrations -migrations-table schema_migrations

migrate-sso-up:
	@echo "$(GREEN)Применение следующей миграции SSO...$(NC)"
	migrate -database "$(USERS_DB_URL)" -path $(SSO_MIGRATIONS_PATH) up 1

migrate-sso-down:
	@echo "$(YELLOW)Откат последней миграции SSO...$(NC)"
	migrate -database "$(USERS_DB_URL)" -path $(SSO_MIGRATIONS_PATH) down 1

migrate-sso-force:
	@echo "$(YELLOW)Принудительная установка версии $(version)...$(NC)"
	migrate -database "$(USERS_DB_URL)" -path $(SSO_MIGRATIONS_PATH) force $(version)

migrate-sso-create:
	@echo "$(GREEN)Создание новой миграции SSO: $(name)...$(NC)"
	migrate create -ext sql -dir $(SSO_MIGRATIONS_PATH) -seq $(name)

# ===== Миграции Gateway =
migrate-gateway:
	@echo "$(GREEN)Применение миграций API Gateway...$(NC)"
	cd Api_Gateway && go run ./cmd/migrator -db-dsn "$(AURORA_DB_URL)" -migrations-path ./migrations -migrations-table schema_migrations

migrate-gateway-up:
	@echo "$(GREEN)Применение следующей миграции Gateway...$(NC)"
	migrate -database "$(AURORA_DB_URL)" -path $(GATEWAY_MIGRATIONS_PATH) up 1

migrate-gateway-down:
	@echo "$(YELLOW)Откат последней миграции Gateway...$(NC)"
	migrate -database "$(ENS_DB_URL)" -path $(GATEWAY_MIGRATIONS_PATH) down 1

migrate-gateway-force:
	@echo "$(YELLOW)Принудительная установка версии $(version)...$(NC)"
	migrate -database "$(ENS_DB_URL)" -path $(GATEWAY_MIGRATIONS_PATH) force $(version)

migrate-gateway-create:
	@echo "$(GREEN)Создание новой миграции Gateway: $(name)...$(NC)"
	migrate create -ext sql -dir $(GATEWAY_MIGRATIONS_PATH) -seq $(name)

# ===== Миграции Aurora =====
migrate-aurora:
	@echo "$(GREEN)Применение миграций Aurora...$(NC)"
	migrate -database "$(AURORA_DB_URL)" -path $(AURORA_MIGRATIONS_PATH) up

migrate-aurora-up:
	@echo "$(GREEN)Применение следующей миграции Aurora...$(NC)"
	migrate -database "$(AURORA_DB_URL)" -path $(AURORA_MIGRATIONS_PATH) up 1

migrate-aurora-down:
	@echo "$(YELLOW)Откат последней миграции Aurora...$(NC)"
	migrate -database "$(AURORA_DB_URL)" -path $(AURORA_MIGRATIONS_PATH) down 1

migrate-aurora-force:
	@echo "$(YELLOW)Принудительная установка версии $(version)...$(NC)"
	migrate -database "$(AURORA_DB_URL)" -path $(AURORA_MIGRATIONS_PATH) force $(version)

migrate-aurora-create:
	@echo "$(GREEN)Создание новой миграции Aurora: $(name)...$(NC)"
	migrate create -ext sql -dir $(AURORA_MIGRATIONS_PATH) -seq $(name)

# ===== Все миграции сразу =====
migrate-all: migrate-sso migrate-gateway
	@echo "$(GREEN)Все миграции применены!$(NC)"

# ===== Управление базами данных =====
db-list:
	@echo "$(GREEN)Список баз данных:$(NC)"
	psql -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) -c "\l" | grep -E "users|ens|aurora"

db-create:
	@echo "$(GREEN)Создание баз данных...$(NC)"
	psql -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) -c "CREATE DATABASE users;" 2>/dev/null || echo "БД users уже существует"
	psql -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) -c "CREATE DATABASE aurora;" 2>/dev/null || echo "БД aurora уже существует"
	psql -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) -c "CREATE DATABASE parser;" 2>/dev/null || echo "БД aurora уже существует"
	@echo "$(GREEN)Готово!$(NC)"

db-drop:
	@echo "$(RED)ВНИМАНИЕ! Это удалит все базы данных!$(NC)"
	@read -p "Вы уверены? (y/N): " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		psql -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) -c "DROP DATABASE IF EXISTS users WITH (FORCE);"; \
		psql -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) -c "DROP DATABASE IF EXISTS ens WITH (FORCE);"; \
		psql -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) -c "DROP DATABASE IF EXISTS aurora WITH (FORCE);"; \
		echo "$(GREEN)Базы данных удалены$(NC)"; \
	else \
		echo "$(YELLOW)Операция отменена$(NC)"; \
	fi

# ===== Docker команды =====
docker-up:
	@echo "$(GREEN)Запуск Docker контейнеров...$(NC)"
	docker-compose up -d
	@echo "$(GREEN)Контейнеры запущены. Просмотр логов: make docker-logs$(NC)"

docker-down:
	@echo "$(YELLOW)Остановка Docker контейнеров...$(NC)"
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-rebuild:
	@echo "$(GREEN)Пересборка и запуск Docker контейнеров...$(NC)"
	docker-compose up -d --build

# ===== Тестирование =====
SERVICES :=  Recommendation_Service SSO

test-all:
	@echo "Running tests for all services with internal coverage..."
	@for service in $(SERVICES); do \
		if [ -f "$$service/go.mod" ]; then \
			echo "------------------------------------------"; \
			echo "Testing service: $$service"; \
			cd $$service && \
			go test -v -gcflags="-N -l" -coverprofile=coverage.out -coverpkg=./internal/... ./... && \
			echo "Coverage report created: $$service/coverage.out" && \
			echo -n "Total coverage for $$service: " && \
			go tool cover -func=coverage.out | grep total | awk '{print $$3}' && \
			cd ..; \
		fi \
	done
	@echo "------------------------------------------"
	@echo "All tests finished!"
test:
	@echo "$(GREEN)Запуск тестов...$(NC)"
	go test ./... -v

test-sso:
	@echo "$(GREEN)Запуск тестов SSO...$(NC)"
	cd SSO && go test ./... -v

test-gateway:
	@echo "$(GREEN)Запуск тестов Gateway...$(NC)"
	cd Api_Gateway && go test ./... -v

test-coverage:
	@echo "$(GREEN)Запуск тестов с покрытием...$(NC)"
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

# ===== Очистка =====
clean:
	@echo "$(YELLOW)Очистка бинарников...$(NC)"
	find . -type f -name "*.exe" -delete
	find . -type f -name "*.test" -delete
	find . -type f -name "coverage.out" -delete
	@echo "$(GREEN)Готово!$(NC)"

# ===== Документация (go doc) =====
define print_doc
	@echo "$(GREEN)========== $(1) ==========$(NC)"
	cd $(2) && for pkg in $$(go list ./internal/... ./cmd/... 2>/dev/null); do \
		rel=".$${pkg#$$(go list -m)}"; \
		echo "$(YELLOW)--- $$rel ---$(NC)"; \
		go doc ./$$rel 2>/dev/null | head -15; \
		echo ""; \
	done
endef

doc-api:
	$(call print_doc,API Gateway,Api_Gateway)

doc-sso:
	$(call print_doc,SSO,SSO)

doc-recs:
	$(call print_doc,Recommendation Service,Recommendation_Service)

doc-cache:
	$(call print_doc,Cache Service,Cache_Service)

doc-all: doc-sso doc-api doc-recs doc-cache
	@echo "$(GREEN)Документация всех сервисов выведена выше.$(NC)"

# ===== Документация (HTML — Sphinx-like) =====
DOCS_DIR := docs

doc-gen: doc-gen-sso doc-gen-api doc-gen-recs doc-gen-cache 
	@echo "$(GREEN)Генерация HTML-документации завершена.$(NC)"
	@echo "$(GREEN)Откройте docs_build/index.html в браузере.$(NC)"

PYTHON := /Users/dm_sfr/.pyenv/shims/python3

doc-gen-parsing:
	@echo "$(GREEN)Генерация документации Parsing Service...$(NC)"
	$(PYTHON) -m sphinx.ext.apidoc -o docs_src/parsing/_autogen Parsing_Service/ --force --separate || true
	$(PYTHON) -m sphinx -b html docs_src/parsing/ docs_build/parsing/ -q

doc-gen-analitic:
	@echo "$(GREEN)Генерация документации Analitic Service...$(NC)"
	$(PYTHON) -m sphinx.ext.apidoc -o docs_src/analitic/_autogen Analitic_Service/ --force --separate || true
	$(PYTHON) -m sphinx -b html docs_src/analitic/ docs_build/analitic/ -q

doc-gen-sso:
	@echo "$(GREEN)Генерация документации SSO...$(NC)"
	@command -v golds >/dev/null 2>&1 || go install go101.org/golds@latest
	cd SSO && golds -gen -nouses -dir=../$(DOCS_DIR)/sso ./... 2>&1 | grep "Docs are generated"

doc-gen-api:
	@echo "$(GREEN)Генерация документации API Gateway...$(NC)"
	@command -v golds >/dev/null 2>&1 || go install go101.org/golds@latest
	cd Api_Gateway && golds -gen -nouses -dir=../$(DOCS_DIR)/api-gateway ./... 2>&1 | grep "Docs are generated"

doc-gen-recs:
	@echo "$(GREEN)Генерация документации Recommendation Service...$(NC)"
	@command -v golds >/dev/null 2>&1 || go install go101.org/golds@latest
	cd Recommendation_Service && golds -gen -nouses -dir=../$(DOCS_DIR)/recommendation ./... 2>&1 | grep "Docs are generated"

doc-gen-cache:
	@echo "$(GREEN)Генерация документации Cache Service...$(NC)"
	@command -v golds >/dev/null 2>&1 || go install go101.org/golds@latest
	cd Cache_Service && golds -gen -nouses -dir=../$(DOCS_DIR)/cache ./... 2>&1 | grep "Docs are generated"

doc-serve:
	@echo "$(GREEN)Запуск сервера HTML-документации на http://localhost:56789$(NC)"
	@command -v golds >/dev/null 2>&1 || go install go101.org/golds@latest
	@echo "$(YELLOW)Нажмите Ctrl+C для остановки$(NC)"
	golds -dir=$(DOCS_DIR)/sso,$(DOCS_DIR)/api-gateway,$(DOCS_DIR)/recommendation,$(DOCS_DIR)/cache

# ===== Утилиты =====
deps:
	@echo "$(GREEN)Установка зависимостей...$(NC)"
	go mod download
	cd SSO && go mod download
	cd Api_Gateway && go mod download

fmt:
	@echo "$(GREEN)Форматирование кода...$(NC)"
	go fmt ./...
	cd SSO && go fmt ./...
	cd Api_Gateway && go fmt ./...

lint:
	@echo "$(GREEN)Запуск линтера...$(NC)"
	golangci-lint run ./...
