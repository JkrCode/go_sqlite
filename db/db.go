package db

import (
	"database/sql"
	"fmt"
	"go_sqlite_demo/models"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func Run(params *models.EnvironmentParams) (*sql.DB, error) {

    dbPath := filepath.Join(params.HomeDir, "go_projects", "go_sqlite", "db", "mydatabase.db")
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }

    // Read the schema file
    schemaPath := filepath.Join(params.HomeDir, "go_projects", "go_sqlite","db", "schema", "schema.sql")
    schemaSQL, err := os.ReadFile(schemaPath)
    if err != nil {
		return nil, err

    }

    // Execute the schema SQL
    _, err = db.Exec(string(schemaSQL))
    if err != nil {
		return nil, err
    }
	return db, nil
}


func InsertMessage(db *sql.DB, message models.Message,params *models.EnvironmentParams)(int64, error){
stmt, err := prepareQuery(db, "insert_message.sql", params)
    if err != nil {
        return 0, err
    }
    defer stmt.Close()

    result, err := stmt.Exec(message.Severity , message.DescriptionText, message.ReceivedDateTime.Format("2006.01.02 15.04.05)"))
    if err != nil {
        return 0, err
    }
    return result.LastInsertId()
}


func prepareQuery(db *sql.DB, filename string, params *models.EnvironmentParams) (*sql.Stmt, error) {
	
    queryPath := filepath.Join(params.HomeDir, "go_projects", "go_sqlite","db", "schema", filename)
    querySQL, err := os.ReadFile(queryPath)
    if err != nil {
        return nil, fmt.Errorf("error reading query file: %w", err)
    }

    stmt, err := db.Prepare(string(querySQL))
    if err != nil {
        return nil, fmt.Errorf("error preparing query: %w", err)
    }

    return stmt, nil
}

