# Nama file migrasi
MIGRATION_NAME ?= new_migration

# Direktori untuk menyimpan file migrasi
MIGRATIONS_DIR = ./db/migrations
APP_DIR = ./cmd/app

DB_USER ?= postgres
DB_PASSWORD ?= Kons123
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_NAME ?= depublic

# URL koneksi basis data
DB_URL = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Nama binari golang-migrate
ifeq ($(OS),Windows_NT)
	MIGRATE_BIN = $(shell powershell -Command "Get-Command migrate | Select-Object -ExpandProperty Definition")
else
	MIGRATE_BIN = $(shell which migrate)
endif

# Mendapatkan timestamp saat ini
ifeq ($(OS),Windows_NT)
	TIMESTAMP = $(shell powershell -Command "Get-Date -Format yyyyMMddHHmmss")
else
	TIMESTAMP = $(shell date +%Y%m%d%H%M%S)
endif

# Membuat migrasi baru dengan timestamp tanpa sequential
.PHONY: create-migration
create-migration:
	@echo "Creating new migration: $(TIMESTAMP)_$(MIGRATION_NAME)"
ifeq ($(OS),Windows_NT)
	@type NUL > $(MIGRATIONS_DIR)/$(TIMESTAMP)_$(MIGRATION_NAME).up.sql
	@type NUL > $(MIGRATIONS_DIR)/$(TIMESTAMP)_$(MIGRATION_NAME).down.sql
else
	@touch $(MIGRATIONS_DIR)/$(TIMESTAMP)_$(MIGRATION_NAME).up.sql
	@touch $(MIGRATIONS_DIR)/$(TIMESTAMP)_$(MIGRATION_NAME).down.sql
endif

# Menjalankan semua migrasi
.PHONY: migrate-up
migrate-up:
	@echo "Running all up migrations"
	@echo $(MIGRATE_BIN)
	@echo $(DB_URL)
	@"$(MIGRATE_BIN)" -path "$(MIGRATIONS_DIR)" -database "$(DB_URL)" up

# Membatalkan migrasi terakhir
.PHONY: migrate-down
migrate-down:
	@echo "Rolling back last migration"
	@"$(MIGRATE_BIN)" -path "$(MIGRATIONS_DIR)" -database "$(DB_URL)" down 1

# Membatalkan semua migrasi
.PHONY: migrate-down-all
migrate-down-all:
	@echo "Rolling back all migrations"
	@"$(MIGRATE_BIN)" -path "$(MIGRATIONS_DIR)" -database "$(DB_URL)" down

# Melihat status migrasi
.PHONY: migrate-status
migrate-status:
	@echo "Checking migration status"
	@"$(MIGRATE_BIN)" -path "$(MIGRATIONS_DIR)" -database "$(DB_URL)" version

# Membuat ulang (rollback dan kemudian migrasi ulang) migrasi terakhir
.PHONY: migrate-refresh
migrate-refresh:
	@echo "Refreshing migrations"
	@"$(MIGRATE_BIN)" -path "$(MIGRATIONS_DIR)" -database "$(DB_URL)" down 1
	@"$(MIGRATE_BIN)" -path "$(MIGRATIONS_DIR)" -database "$(DB_URL)" up 1

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
