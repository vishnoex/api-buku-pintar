package main

import (
	"buku-pintar/pkg/config"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	configPath := flag.String("config", "./config.json", "path to application config file")
	migrationsPath := flag.String("path", "./migrations", "path to migration files")
	command := flag.String("command", "up", "migration command: up, down, steps, force, version, drop")
	steps := flag.Int("steps", 1, "number of migration steps for down or steps command")
	forceVersion := flag.Int("version", 0, "migration version for force command")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dbConfig := cfg.GetDatabaseConfig()
	databaseURL := mysqlMigrationURL(dbConfig)
	sourceURL := "file://" + strings.TrimRight(*migrationsPath, "/")

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		log.Fatalf("failed to initialize migrations: %v", err)
	}
	defer func() {
		sourceErr, databaseErr := m.Close()
		if sourceErr != nil {
			log.Printf("failed to close migration source: %v", sourceErr)
		}
		if databaseErr != nil {
			log.Printf("failed to close migration database: %v", databaseErr)
		}
	}()

	if err := runMigrationCommand(m, *command, *steps, *forceVersion); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
}

func mysqlMigrationURL(db config.DatabaseConfig) string {
	params := db.Params
	if params == "" {
		params = "parseTime=true"
	}

	return fmt.Sprintf(
		"mysql://%s:%s@tcp(%s:%s)/%s?%s",
		url.QueryEscape(db.User),
		url.QueryEscape(db.Password),
		db.Host,
		db.Port,
		db.Name,
		params,
	)
}

func runMigrationCommand(m *migrate.Migrate, command string, steps int, forceVersion int) error {
	switch strings.ToLower(strings.TrimSpace(command)) {
	case "up":
		return ignoreNoChange(m.Up())
	case "down":
		if steps <= 0 {
			return errors.New("steps must be greater than zero for down command")
		}
		return ignoreNoChange(m.Steps(-steps))
	case "steps":
		if steps == 0 {
			return errors.New("steps cannot be zero")
		}
		return ignoreNoChange(m.Steps(steps))
	case "force":
		if forceVersion <= 0 {
			return errors.New("version must be greater than zero for force command")
		}
		return m.Force(forceVersion)
	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			if errors.Is(err, migrate.ErrNilVersion) {
				fmt.Fprintln(os.Stdout, "migration version: nil")
				return nil
			}
			return err
		}
		fmt.Fprintf(os.Stdout, "migration version: %d dirty: %t\n", version, dirty)
		return nil
	case "drop":
		return m.Drop()
	default:
		return fmt.Errorf("unsupported migration command: %s", command)
	}
}

func ignoreNoChange(err error) error {
	if errors.Is(err, migrate.ErrNoChange) {
		fmt.Fprintln(os.Stdout, "no migration changes to apply")
		return nil
	}
	return err
}
