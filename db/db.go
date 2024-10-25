package db

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"go_sqlite_demo/models"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema/schema.sql
var schemaSQL string

//go:embed schema/insert_message.sql
var insertMessageSQL string

//go:embed schema/drop_table.sql
var dropTableSQL string

//go:embed schema/get_message.sql
var getMessage string

func Run() (*sql.DB, error) {

	// Get the current working directory (assuming the executable is run from the project root).
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		return nil, err
	}

	// Use a relative path to the database file located in the db folder.
	dbPath := filepath.Join(workingDir, "db", "mydatabase.db")

	// Open the SQLite database.
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("Error opening the database:", err)
		return nil, err
	}

	// Execute the embedded schema SQL from db/schema/schema.sql.
	_, err = db.Exec(schemaSQL)
	if err != nil {
		fmt.Println("Error executing schema SQL:", err)
		return nil, err
	}

	return db, nil
}

// InsertMessage inserts a message into the database using the embedded SQL.
func InsertMessage(db *sql.DB, message models.Message) (int64, error) {
	// Use the embedded SQL for inserting the message.
	stmt, err := db.Prepare(insertMessageSQL)
	if err != nil {
		return 0, fmt.Errorf("error preparing query: %w", err)
	}
	defer stmt.Close()

	// Execute the query with the provided message data.
	result, err := stmt.Exec(message.Severity, message.DescriptionText, message.ReceivedDateTime.Format("2006.01.02 15.04.05"))
	if err != nil {
		return 0, fmt.Errorf("error executing query: %w", err)
	}
	return result.LastInsertId()
}

// ResetDB resets the database by dropping the table using the embedded SQL.
func ResetDB(db *sql.DB) (int64, error) {
	// Use the embedded SQL for dropping the table.
	stmt, err := db.Prepare(dropTableSQL)
	if err != nil {
		return 0, fmt.Errorf("error preparing query: %w", err)
	}
	defer stmt.Close()

	// Execute the query.
	result, err := stmt.Exec()
	if err != nil {
		return 0, fmt.Errorf("error executing query: %w", err)
	}

	// Get the number of rows affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error getting rows affected: %w", err)
	}

	fmt.Println("Rows affected:", rowsAffected)
	return rowsAffected, nil
}
func GetRowsAndPutInChannel(ctx context.Context, conn *sql.DB, ch chan<- models.Message) {
	defer close(ch)

	// Query the messages table for the necessary fields.
	rowPointer, err := conn.Query(getMessage)
	if err != nil {
		fmt.Println("Error getting rows from db:", err)
		return
	}
	defer rowPointer.Close()

	// Iterate over the rows and send messages to the channel.
	for rowPointer.Next() {
		var message models.Message
		var receivedDateTime string

		err := rowPointer.Scan(&message.Severity, &message.DescriptionText, &receivedDateTime)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}

		message.ReceivedDateTime, err = time.Parse("2006.01.02 15.04.05", receivedDateTime)
		if err != nil {
			fmt.Println("Error parsing date from db row:", err)
			return
		}

		select {
		case <-ctx.Done():
			return
		case ch <- message:
		}
	}

	// Check if any error occurred during row iteration.
	if err := rowPointer.Err(); err != nil {
		fmt.Println("Error iterating over rows:", err)
	}
}
