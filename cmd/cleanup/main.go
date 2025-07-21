package main

import (
    "database/sql"
    "log"
    "time"
    _ "github.com/lib/pq"
    "github.com/joho/godotenv"
)

func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }

    // Connect to database
    db, err := sql.Open("postgres", "postgres://postgres:password@localhost:5432/reflection_db?sslmode=disable")
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    // Test the connection
    if err := db.Ping(); err != nil {
        log.Fatal("Failed to ping database:", err)
    }

    log.Println("Connected to database successfully")

    // Check current timezone
    var timezone string
    err = db.QueryRow("SELECT current_setting('timezone')").Scan(&timezone)
    if err != nil {
        log.Printf("Error getting timezone: %v", err)
    } else {
        log.Printf("Database timezone: %s", timezone)
    }

    // Check what dates exist in reflection_tracking
    rows, err := db.Query(`
        SELECT reflected_date, COUNT(*) as count
        FROM reflection_tracking
        GROUP BY reflected_date
        ORDER BY reflected_date DESC
        LIMIT 10
    `)
    if err != nil {
        log.Printf("Error querying reflection_tracking: %v", err)
        return
    }
    defer rows.Close()

    log.Println("Existing tracking dates:")
    for rows.Next() {
        var date time.Time
        var count int
        if err := rows.Scan(&date, &count); err != nil {
            log.Printf("Error scanning row: %v", err)
            continue
        }
        log.Printf("  %s: %d records", date.Format("2006-01-02"), count)
    }

    // Get current date in Go
    currentDate := time.Now().Format("2006-01-02")
    log.Printf("Current date in Go: %s", currentDate)

    // Check what PostgreSQL thinks is current date
    var dbCurrentDate time.Time
    err = db.QueryRow("SELECT CURRENT_DATE").Scan(&dbCurrentDate)
    if err != nil {
        log.Printf("Error getting PostgreSQL current date: %v", err)
    } else {
        log.Printf("PostgreSQL current date: %s", dbCurrentDate.Format("2006-01-02"))
    }

    // Clean up any records with wrong dates (older than today)
    result, err := db.Exec(`
        DELETE FROM reflection_tracking 
        WHERE reflected_date < $1
    `, currentDate)
    if err != nil {
        log.Printf("Error cleaning up old records: %v", err)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Printf("Error getting rows affected: %v", err)
    } else {
        log.Printf("Cleaned up %d old tracking records", rowsAffected)
    }

    log.Println("Cleanup completed")
} 