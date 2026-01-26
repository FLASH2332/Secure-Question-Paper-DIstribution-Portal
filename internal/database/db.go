package database

import (
    "database/sql"
    "fmt"
    "os"
    
    _ "github.com/go-sql-driver/mysql"
    "github.com/joho/godotenv"
)

func Connect() (*sql.DB, error) {
    // Load .env file
    godotenv.Load()
    
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_NAME"),
    )
    
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    return db, nil
}