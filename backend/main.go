package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	dev := os.Getenv("DEV") == "true"
	migrationSourceURL := os.Getenv("MIGRATION_SOURCE_URL")

	dbHandle, err := initDB(dbPath, migrationSourceURL)
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}
	defer dbHandle.Close()

	// Create an HTTP server
	mux := http.NewServeMux()

	if dev {
		log.Println("Running in development mode, setting up proxy to locally running dev React app")
		targetURL, err := url.Parse("http://localhost:3000")
		if err != nil {
			log.Fatal(err)
		}
		// Create a reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		// Set up the reverse proxy for all other requests
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			log.Println("Proxying request to React app")
			proxy.ServeHTTP(w, r)
		})
	} else {
		mux.Handle("/", http.FileServer(http.Dir("../frontend/build"))) // Serve the built React app
	}

	httpServer, err := NewHTTPServer(dbHandle, mux)
	if err != nil {
		log.Fatalf("failed to create HTTP server: %v", err)
	}
	defer httpServer.Close()

	httpServer.Start()
}

// Initialize database connection.
func initDB(dbPath, migrationSourceURL string) (*sql.DB, error) {
	dbHandle, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Couldn't get DB handle with path %s", dbPath)

	}
	driver, err := sqlite3.WithInstance(dbHandle, &sqlite3.Config{})
	if err != nil {
		log.Fatalf("Can't get DB driver for migrations: %s", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		migrationSourceURL,
		"sqlite3",
		driver,
	)
	if err != nil {
		log.Fatalf("Can't get migrate instance: %s", err)
	}
	version, _, err := m.Version()
	if err != nil {
		log.Printf("Can't get DB version! %s", err)
	}
	log.Println("DB version is", version)
	err = m.Migrate(1)
	if errors.Is(err, migrate.ErrNoChange) {
		log.Println("No migrations to run")
	} else if err != nil {
		log.Fatalf("Can't run migrations: %s", err)
	}
	return dbHandle, nil
}
