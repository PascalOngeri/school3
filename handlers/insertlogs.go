package handlers

import (
    "database/sql"
    "log"
   
)

// InsertLog records user activities into the logs table
func InsertLog(db *sql.DB, user string, activity string) error {
    // Get the current timestamp
   
    // Prepare the SQL insert statement
    sqlInsert := "INSERT INTO logs (user, activities) VALUES (?, ?)"
    _, err := db.Exec(sqlInsert, user, activity)
    if err != nil {
        log.Printf("Error inserting log: %v", err)
        return err
    }

    return nil
}
