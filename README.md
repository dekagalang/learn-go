# CRUD API - Go + Gin + GORM + PostgreSQL + Atlas

Production-ready CRUD API dengan declarative database schema using **Atlas** (mirip Drizzle ORM).

## 📁 Struktur Project

```
crud-api/
├── database/
│   ├── database.go           # Database connection
│   └── migration.go          # Atlas migration runner
├── migrations/
│   └── 20260307000001.sql    # Atlas-generated migrations
├── models/
│   └── user.go               # User model
├── handlers/
│   └── user_handler.go       # CRUD handlers
├── schema.hcl                # Declarative schema (like Drizzle)
├── atlas.hcl                 # Atlas configuration
├── docker-compose.yml        # PostgreSQL & PgAdmin
├── main.go                   # Entry point
└── go.mod                    # Go module
```

## 🚀 Cara Menjalankan

### 1. Start Database
```bash
docker-compose up -d
```

Verifikasi database berjalan:
- **PostgreSQL:** localhost:5433
- **PgAdmin:** http://localhost:5051 (admin@example.com / admin123)

### 2. Run Server
```bash
go run main.go
```

Atau gunakan binary yang sudah compiled:
```bash
./crud-api
```

Server berjalan di `http://localhost:8080`

## 📊 Atlas Declarative Migrations

### Apa itu Atlas?

Atlas adalah modern migration tool yang **mirip dengan Drizzle ORM di Node.js**.

**Perbedaan dengan golang-migrate:**
- golang-migrate: SQL-first (tulis SQL manual)
- Atlas: Schema-first (deklaratif seperti Drizzle)

### Schema Definition (schema.hcl)

```hcl
schema "public" {}

table "users" {
  schema = schema.public
  
  column "id" {
    null           = false
    type           = bigserial
    auto_increment = true
  }
  
  column "name" {
    null = false
    type = text
  }
  
  column "email" {
    null = false
    type = text
  }
  
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  
  primary_key {
    columns = [column.id]
  }
  
  index "idx_users_email" {
    columns = [column.email]
    unique  = true
  }
}
```

### Migration Files

Atlas auto-generates migration files dari schema diff:

```sql
-- migrations/20260307000001.sql
CREATE TABLE "public"."users" (
  "id" bigserial NOT NULL,
  "name" text NOT NULL,
  "email" text NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "idx_users_email" ON "public"."users" ("email");
```

### Auto-Run Migrations

Migrations auto-run saat server startup:

```go
database.Connect()
if err := database.RunMigrations(database.SqlDB); err != nil {
    log.Fatal("Migration failed:", err)
}
```

## 📡 API Endpoints

### Create User (POST)
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Deka",
    "email": "deka@mail.com"
  }'
```

**Response:**
```json
{
  "id": 1,
  "name": "Deka",
  "email": "deka@mail.com",
  "created_at": "2026-03-07T10:30:00Z",
  "updated_at": "2026-03-07T10:30:00Z"
}
```

### Get All Users (GET)
```bash
curl http://localhost:8080/api/users
```

### Get User by ID (GET)
```bash
curl http://localhost:8080/api/users/1
```

### Update User (PUT)
```bash
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Deka Updated",
    "email": "deka2@mail.com"
  }'
```

### Delete User (DELETE)
```bash
curl -X DELETE http://localhost:8080/api/users/1
```

## 🔧 Konfigurasi Database

**Docker Compose:**
- **Host:** localhost
- **Port:** 5433
- **Database:** crud_api_dev
- **User:** postgres
- **Password:** postgres

**Connection String (di database.go):**
```
postgres://postgres:postgres@localhost:5433/crud_api_dev?sslmode=disable
```

## 📁 Menambah Tabel/Kolom Baru

### Workflow (mirip Drizzle):

#### 1. Edit schema.hcl

```hcl
table "products" {
  schema = schema.public
  
  column "id" {
    type = bigserial
  }
  
  column "name" {
    type = text
  }
  
  column "price" {
    type = decimal
  }
  
  primary_key {
    columns = [column.id]
  }
}
```

#### 2. Generate migration dengan Atlas CLI (opsional)

```bash
# Install Atlas CLI
go install ariga.io/atlas/cmd/atlas@latest

# Generate migration dari schema diff
atlas migrate diff create_products --dev "docker://postgres/16/crud_api_dev?user=postgres&password=postgres"
```

Ini akan create `migrations/20260307000002.sql`

#### 3. Atau manual create migration file

Buat file: `migrations/20260307000002.sql`

```sql
-- Create "products" table
CREATE TABLE "public"."products" (
  "id" bigserial NOT NULL,
  "name" text NOT NULL,
  "price" numeric NOT NULL,
  PRIMARY KEY ("id")
);
```

#### 4. Create Go model

`models/product.go`

```go
package models

import "time"

type Product struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name"`
    Price     float64   `json:"price"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

#### 5. Restart server

```bash
go run main.go
```

Migration auto-run ✓

## 🛠️ Atlas CLI Commands

Setelah install `atlas` CLI:

```bash
# Show current schema in database
atlas schema inspect -u "postgres://postgres:postgres@localhost:5433/crud_api_dev?sslmode=disable"

# Generate migration dari schema diff
atlas migrate diff my_migration -d "file://migrations" \
  --dir "file://migrations"

# Apply migrations
atlas migrate apply -u "postgres://postgres:postgres@localhost:5433/crud_api_dev?sslmode=disable"

# Check migration status
atlas migrate status -u "postgres://postgres:postgres@localhost:5433/crud_api_dev?sslmode=disable"
```

## 🎯 Keunggulan Atlas Setup

✅ **Declarative schema** (seperti Drizzle ORM)  
✅ **Auto-generate migrations** dari schema diff  
✅ **Version control** untuk migrations  
✅ **Schema-first approach** (lebih aman)  
✅ **Auto-run saat startup**  
✅ **Production-ready**  

## 📋 Migration Tracking

Atlas mencatat migrations dalam tabel `atlas_schema_migrations`:

```sql
SELECT * FROM atlas_schema_migrations;
```

Output:
```
 version |          description          | type | installed_on | success | execution_time
---------+-------------------------------+------+--------------+---------+----------------
    1   | 20260307000001.sql            |  1   | 2026-03-07   | t       | 0
```

## 🔄 Reset Database (untuk development)

```bash
# Stop containers & remove volumes
docker-compose down -v

# Start fresh
docker-compose up -d

# Restart app (migrations auto-run)
go run main.go
```

## 📚 Perbandingan: Drizzle vs Atlas

| Fitur | Drizzle (TypeScript) | Atlas (Go) |
|-------|-------|--------|
| Schema definition | TypeScript code | HCL file |
| Auto-generate migrations | ✓ | ✓ |
| Version control | ✓ | ✓ |
| Rollback support | ✓ | ✓ |
| Type-safe queries | ✓ | Partial (GORM) |

## 🐛 Troubleshooting

**Migration gagal?**
1. Pastikan PostgreSQL running: `docker-compose ps`
2. Cek konfigurasi di `database.go`
3. Lihat error log saat startup
4. Cek tabel `atlas_schema_migrations` untuk history

**Ingin lihat schema database?**
```bash
# Dengan Atlas CLI
atlas schema inspect -u "postgres://..."

# Atau di PgAdmin
# http://localhost:5051
```

---

✨ **Production-ready dengan declarative schema management (Atlas = "Drizzle untuk Go")**

# learn-go
# learn-go
# learn-go
