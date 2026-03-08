# Atlas Migration Workflow

Dokumentasi lengkap cara generate migrations dengan Atlas.

## 📋 Struktur Production-Grade dengan Atlas

```
crud-api/
├── schema.hcl              # Source of truth (deklaratif)
├── atlas.hcl               # Atlas config
├── migrations/
│   └── 20260307000001.sql  # Generated SQL (version controlled)
├── database/migration.go   # Migration executor
└── main.go                 # Auto-run migrations
```

---

## 🔄 Workflow Migrasi (Production-Grade)

### Step 1: Edit Schema (Source of Truth)

File: `schema.hcl`

```hcl
schema "public" {}

table "users" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
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

### Step 2: Generate Migration SQL

Atlas akan **otomatis generate SQL** dari diff schema:

```bash
atlas migrate diff --dir file://migrations \
  --to file://schema.hcl \
  --dev-url "docker://postgres/16/dev?search_path=public" \
  --name "create_users_table"
```

**Output:** `migrations/20260307000001_create_users_table.sql`

### Step 3: Review Generated SQL

```bash
cat migrations/20260307000001_create_users_table.sql
```

Output:
```sql
-- Create "users" table
CREATE TABLE "public"."users" (
  "id" bigserial NOT NULL,
  "name" text NOT NULL,
  "email" text NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
-- Create index "idx_users_email"
CREATE UNIQUE INDEX "idx_users_email" ON "public"."users" ("email");
```

### Step 4: Commit & Deploy

```bash
# Git commit migration
git add migrations/
git commit -m "migration: create users table"

# Deploy: migrations auto-run saat startup
go run main.go
```

---

## 📚 Contoh: Menambah Table Baru

Misal ingin tambah table `products`:

### 1. Edit schema.hcl - Tambah Table

```hcl
schema "public" {}

table "users" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
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

# ➕ TABLE BARU
table "products" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "name" {
    null = false
    type = text
  }
  column "price" {
    null = false
    type = decimal
  }
  column "user_id" {
    type = bigint
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_products_users" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }
}
```

### 2. Generate Migration

```bash
atlas migrate diff --dir file://migrations \
  --to file://schema.hcl \
  --dev-url "docker://postgres/16/dev?search_path=public" \
  --name "add_products_table"
```

**Output:** `migrations/20260307000002_add_products_table.sql`

```sql
-- Create "products" table
CREATE TABLE "public"."products" (
  "id" bigserial NOT NULL,
  "name" text NOT NULL,
  "price" numeric NOT NULL,
  "user_id" bigint,
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
-- Create "fk_products_users" foreign key
ALTER TABLE "public"."products" ADD CONSTRAINT "fk_products_users" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE;
```

### 3. Create Go Model

`models/product.go`:

```go
package models

import "time"

type Product struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name"`
    Price     float64   `json:"price"`
    UserID    *uint     `json:"user_id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 4. Commit & Deploy

```bash
git add schema.hcl migrations/
git commit -m "migration: add products table with user reference"

go run main.go  # Auto-execute
```

---

## 🚀 Complete Commands Reference

### Generate Fresh Migration dari Schema

```bash
atlas migrate diff my_migration \
  --dir file://migrations \
  --to file://schema.hcl \
  --dev-url "docker://postgres/16/dev?search_path=public"
```

### Inspect Current Database Schema

```bash
atlas schema inspect \
  -u "postgres://postgres:postgres@localhost:5433/crud_api_dev?sslmode=disable"
```

### Check Migration Status

```bash
# Lihat applied migrations
SELECT * FROM atlas_schema_migrations ORDER BY version DESC;
```

---

## 💡 Best Practices

### ✅ DO:

1. Edit `schema.hcl` dulu (source of truth)
2. Generate migration dengan Atlas
3. Review SQL sebelum commit
4. Version control migrations
5. Test di dev environment terlebih dahulu

### ❌ DON'T:

1. Edit migration SQL manual (error-prone)
2. Run DDL queries langsung di production
3. Forget to test rollback
4. Commit without reviewing SQL

---

## 🔄 Workflow Comparison

### ❌ OLD STYLE (golang-migrate)

```
1. Tulis 0001_create_users.up.sql
2. Tulis 0001_create_users.down.sql
3. Run migrations
4. Potensi error & inconsistency
```

### ✅ ATLAS WORKFLOW (Production-Grade)

```
1. Edit schema.hcl (deklaratif)
2. Atlas generate SQL (automatic)
3. Review & commit
4. Deploy / migrations auto-run
5. Guaranteed consistency
```

---

## 📊 Current Setup

**Current schema state:**
- ✅ `schema.hcl` - defined
- ✅ `migrations/20260307000001.sql` - generated
- ✅ Auto-migration runner - ready
- ✅ Production-ready - yes

**To add new tables:**
1. Edit `schema.hcl`
2. Run: `atlas migrate diff --dir file://migrations --to file://schema.hcl --dev-url "docker://postgres/16/dev?search_path=public" --name "your_change"`
3. Deploy
