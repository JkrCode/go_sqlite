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
	ctx, cancel := context.WithCancel(context.Background())
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
		db.GetRowsAndPutInChannel(ctx, conn, ch1)
	}()

	go func() {
		defer wg.Done()
		defer close(ch2)
		pipeline1(ctx, ch1, ch2)
	}()

	go func() {
		defer wg.Done()
		pipeline2(ctx, ch2)
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	fmt.Println("\nServer is shutting down...")
	db.ResetDB(conn)

}

func pipeline1(ctx context.Context, in <-chan models.Message, out chan<- models.Message) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-in:
			if !ok {
				return
			}

			fmt.Println("Pipeline 1 processing message:", msg)

			select {
			case <-ctx.Done():
				return
			case out <- msg:
				// Message forwarded successfully
			}
		}
	}
}

func pipeline2(ctx context.Context, in <-chan models.Message) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-in:
			if !ok {
				return
			}
			fmt.Println("Pipeline 2 processing message:", msg)
		}
	}
}
