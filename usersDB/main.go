package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Variables injected by -X flag
var appVersion = "unknown"
var gitVersion = "unknown"
var lastCommitTime = "unknown"
var lastCommitHash = "unknown"
var lastCommitUser = "unknown"
var buildTime = "unknown"

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err := migrateDB(cfg); err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	rest, err := NewUserRest(buildInfo(), cfg)
	if err != nil {
		log.Fatal(err)
	}

	router.HandleFunc("/api/users", rest.Users).Methods("GET")
	router.HandleFunc("/api/users", rest.AddUser).Methods("POST")
	router.HandleFunc("/api/users/{id}", rest.GetUser).Methods("GET")
	router.HandleFunc("/health", rest.Health).Methods("GET")
	router.HandleFunc("/api/error", rest.Err).Methods("POST")
	router.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{},
	))

	log.Println("starting http server!")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.HTTPPort), NewLogginHandler(router)))
}

func buildInfo() BuildInfo {
	return BuildInfo{
		Version:    appVersion,
		GitVersion: gitVersion,
		BuildTime:  buildTime,
		LastCommit: Commit{
			Author: lastCommitUser,
			ID:     lastCommitHash,
			Time:   lastCommitTime,
		},
	}
}

func migrateDB(cfg *config) error {
	log.Printf("DB connection string %s", cfg.connectionString())

	db, err := sql.Open("mysql", cfg.connectionString())
	if err != nil {
		return fmt.Errorf("can't open db. err: %w", err)
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("can't create migrations driver. err: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", cfg.DBMigrationsPath),
		"mysql",
		driver,
	)

	if err != nil {
		return fmt.Errorf("can't create migrations instance. err: %w", err)
	}

	schemaVerBefore, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("can't get current schem version. err: %w", err)
	}

	if dirty {
		return fmt.Errorf("current version of migration [%d] is dirty. please fix it manually", schemaVerBefore)
	}

	if err = m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			return nil
		}
		return err
	}

	schemaVerAfter, dirty, err := m.Version()
	if err != nil {
		return errors.Wrap(err, "can't get new schema version")
	}

	if dirty {
		return fmt.Errorf("new version of migration [%d] is dirty. please fix it manually", schemaVerAfter)
	}

	log.Printf("migration done to version: %d", schemaVerAfter)

	return nil

}
