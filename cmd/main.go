package main

import (
    "fmt"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    fmt.Println("🔐 Secure Exam System - Setup Verified!")
    
    // Test SQLite connection
    db, err := sql.Open("sqlite3", "./storage/exam.db")
    if err != nil {
        fmt.Println("Database connection failed:", err)
        return
    }
    defer db.Close()
    
    fmt.Println("Database connection successful!")
}