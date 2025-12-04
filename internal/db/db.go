package db

import (
    "context"
    "database/sql"
    "fmt"
    "time"

    _ "github.com/lib/pq" // PostgreSQL driver

    "github.com/example/bridge-service/internal/config"
)

// Connect establishes a connection to the Postgres database using the
// provided configuration. A simple ping is performed to verify the
// connection. The caller is responsible for closing the returned DB.
func Connect(cfg config.DBConfig) (*sql.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
    )
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }
    // Set reasonable connection parameters
    db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(time.Hour)
    // Validate connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := db.PingContext(ctx); err != nil {
        return nil, err
    }
    return db, nil
}