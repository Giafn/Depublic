# Nama file migrasi & versi
MIGRATION_NAME ?= new_migration
# VERSION ?= 20240615102726

# Direktori untuk menyimpan file migrasi
MIGRATIONS_DIR = ./db/migrations
APP_DIR = ./cmd/app

DB_USER ?= dimskuy
DB_PASSWORD ?= 123
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_NAME ?= depublic

DB_URL = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

ifeq ($(OS),Windows_NT)
    MIGRATE_BIN = $(shell where migrate.exe)
else
    MIGRATE_BIN = $(shell which migrate)
endif

TIMESTAMP = $(shell date +%Y%m%d%H%M%S || powershell -Command "Get-Date -Format yyyyMMddHHmmss")

.PHONY: create-migration
create-migration:
	@echo "Creating new migration: $(TIMESTAMP)_$(MIGRATION_NAME)"
	@echo.>$(MIGRATIONS_DIR)/$(TIMESTAMP)_$(MIGRATION_NAME).up.sql
	@echo.>$(MIGRATIONS_DIR)/$(TIMESTAMP)_$(MIGRATION_NAME).down.sql

.PHONY: migrate-up
migrate-up:
	@echo "Running all up migrations"
	@$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database $(DB_URL) up

.PHONY: migrate-down
migrate-down:
	@echo "Rolling back last migration"
	@$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database $(DB_URL) down 1

.PHONY: migrate-down-all
migrate-down-all:
	@echo "Rolling back all migrations"
	@$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database $(DB_URL) down

.PHONY: migrate-status
migrate-status:
	@echo "Checking migration status"
	@$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database $(DB_URL) version

# .PHONY: migrate-fix
# migrate-fix:
# 	@echo "Fix Dirty Version"
# 	@$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database $(DB_URL) force $(VERSION)

.PHONY: migrate-refresh
migrate-refresh:
	@echo "Refreshing migrations"
	@$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database $(DB_URL) down
	@$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database $(DB_URL) up 

.PHONY: run-server
run-server:
	@echo "Running Go application"
	@go run $(APP_DIR)/main.go

# cara menggunakan command
# make create-migration MIGRATION_NAME=nama_migration (menambahkan file migrasi baru)
# make migrate-up (menjalankan semua migrasi)
# make migrate-down (membatalkan migrasi terakhir)
# make migrate-down-all (membatalkan semua migrasi)
# make migrate-status (melihat status migrasi)
# make migrate-refresh (membuat ulang migrasi terakhir)

# make run-server (menjalankan aplikasi)