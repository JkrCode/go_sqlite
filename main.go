package main

import (
	"fmt"
	"go_sqlite_demo/db"
	"go_sqlite_demo/models"
	"os"
	"time"
)



func main() {
	homeDir, err := os.UserHomeDir()
    if err != nil {
        fmt.Println("Error getting homedir path:", err)
        return
    }

	params := models.EnvironmentParams{HomeDir: homeDir}

	// Step 1: Initialize the database and run the schema setup
	fmt.Println("Initializing the database...")
	conn, err := db.Run(&params)
	if err != nil {
		fmt.Println("error starting db with err: ", err)
		return 
	}
	defer conn.Close()

	// Step 3: Create a new Message instance
	testMessage := models.Message{
		Severity:        123,
		DescriptionText: "blablabla",
		ReceivedDateTime: time.Now(),
	}

	// Step 4: Insert the message using InsertMessage
	insertedID, err := db.InsertMessage(conn, testMessage, &params)
	if err != nil {
		fmt.Println("Error inserting message:", err)
		return
	}

	// Step 5: Print the inserted message ID
	fmt.Printf("Message inserted with ID: %d\n", insertedID)
}
