package main

import (
	"fmt"
	"go_sqlite_demo/db"
	"go_sqlite_demo/helper"
	"go_sqlite_demo/models"
	"os"
	"os/signal"
	"syscall"
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
	testMessages := helper.CreateTestData()

	for id := range testMessages {
		dbId, _ := db.InsertMessage(conn, testMessages[id], &params)
		fmt.Printf("Message inserted with ID: %d\n", dbId)
	}

	ch :=make(chan models.Message)
	go db.GetRowsAndPutInChannel(conn, ch)

	go ProcessMessages(ch)


	
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	fmt.Println("\nServer is shutting down...")
	db.ResetDB(conn, &params)

}

func ProcessMessages(ch <-chan models.Message) {
    for message := range ch {
        fmt.Println("Processing message:", message)
    }
}
