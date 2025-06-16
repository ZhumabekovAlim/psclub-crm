package migrations

import (
	"bufio"
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Run executes SQL migration files located in dir. It creates a schema_migrations
// table to track applied versions. Files are executed in lexical order.
func Run(db *sql.DB, dir string) error {
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
        version VARCHAR(255) PRIMARY KEY,
        applied_at DATETIME NOT NULL
    )`); err != nil {
		return err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		version := strings.TrimSuffix(e.Name(), filepath.Ext(e.Name()))
		var count int
		if err := db.QueryRow(`SELECT COUNT(1) FROM schema_migrations WHERE version=?`, version).Scan(&count); err != nil {
			return err
		}
		if count > 0 {
			continue
		}
		path := filepath.Join(dir, e.Name())
		if err := applyFile(db, path); err != nil {
			return err
		}
		if _, err := db.Exec(`INSERT INTO schema_migrations (version, applied_at) VALUES (?, ?)`, version, time.Now()); err != nil {
			return err
		}
	}
	return nil
}

func applyFile(db *sql.DB, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	var stmt strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "--") || trimmed == "" {
			continue
		}
		stmt.WriteString(line)
		if strings.Contains(line, ";") {
			query := strings.TrimSpace(stmt.String())
			query = strings.TrimSuffix(query, ";")
			if query != "" {
				if _, err := db.Exec(query); err != nil {
					return err
				}
			}
			stmt.Reset()
		} else {
			stmt.WriteString(" ")
		}
	}
	if s := strings.TrimSpace(stmt.String()); s != "" {
		if _, err := db.Exec(s); err != nil {
			return err
		}
	}
	return scanner.Err()
}
