package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations executes all pending migrations from the migrations directory
func RunMigrations(db *sql.DB) error {
	// Create atlas_schema_migrations table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS atlas_schema_migrations (
		version bigint PRIMARY KEY,
		description text NOT NULL,
		type integer NOT NULL,
		installed_on timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		success boolean NOT NULL,
		execution_time bigint NOT NULL,
		error_statement text
	);`

	if _, err := db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get all migration files
	migrationFiles := []string{}
	migrationsDir := "migrations"

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			migrationFiles = append(migrationFiles, entry.Name())
		}
	}

	sort.Strings(migrationFiles)

	// Get applied migrations
	rows, err := db.Query("SELECT version FROM atlas_schema_migrations WHERE success = true ORDER BY version")
	if err != nil {
		return fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	appliedVersions := make(map[string]bool)
	for rows.Next() {
		var version int64
		if err := rows.Scan(&version); err != nil {
			return fmt.Errorf("failed to scan version: %w", err)
		}
		appliedVersions[fmt.Sprintf("%d", version)] = true
	}

	// Apply pending migrations
	for _, filename := range migrationFiles {
		version := strings.TrimSuffix(filename, ".sql")

		if appliedVersions[version] {
			log.Printf("✓ Migration %s already applied", version)
			continue
		}

		filepath := filepath.Join(migrationsDir, filename)
		content, err := os.ReadFile(filepath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		// Execute migration
		if _, err := db.Exec(string(content)); err != nil {
			// Record failed migration
			var versionNum int64
			fmt.Sscanf(version, "%d", &versionNum)
			db.Exec("INSERT INTO atlas_schema_migrations (version, description, type, success, execution_time, error_statement) VALUES ($1, $2, $3, $4, $5, $6)",
				versionNum, filename, 1, false, 0, err.Error())
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}

		// Record successful migration
		var versionNum int64
		fmt.Sscanf(version, "%d", &versionNum)
		if _, err := db.Exec("INSERT INTO atlas_schema_migrations (version, description, type, success, execution_time) VALUES ($1, $2, $3, $4, $5)",
			versionNum, filename, 1, true, 0); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", filename, err)
		}

		log.Printf("✓ Applied migration: %s", version)
	}

	return nil
}
