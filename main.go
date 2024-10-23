package main

import (
	"context"
	"fmt"
	"go_sqlite_demo/db"
	"go_sqlite_demo/helper"
	"go_sqlite_demo/models"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	//incorporate context
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("Initializing the database...")
	conn, err := db.Run()
	if err != nil {
		fmt.Println("error starting db with err: ", err)
		return
	}
	defer conn.Close()

	testMessages := helper.CreateTestData()

	for id := range testMessages {
		dbId, _ := db.InsertMessage(conn, testMessages[id])
		fmt.Printf("Message inserted with ID: %d\n", dbId)
	}

	var wg sync.WaitGroup

	ch1 := make(chan models.Message, 100)
	ch2 := make(chan models.Message, 100)

	wg.Add(3)

	go func() {
		defer wg.Done()
		defer close(ch1)
		db.GetRowsAndPutInChannel(conn, ch1)
	}()

	go func() {
		defer wg.Done()
		defer close(ch2)
		pipeline1(ch1, ch2)
	}()

	go func() {
		defer wg.Done()
		pipeline2(ch2)
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	fmt.Println("\nServer is shutting down...")
	db.ResetDB(conn)

}

func pipeline1(ch <-chan models.Message, ch2 chan models.Message) {
	for message := range ch {
		fmt.Println("Processing message:", message)
		ch2 <- message
	}

}

func pipeline2(ch2 <-chan models.Message) {
	for message := range ch2 {
		fmt.Println("Processing message:", message)
	}
}
