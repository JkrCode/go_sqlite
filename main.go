package main

import (
	"context"
	"database/sql"
	"fmt"
	"go_sqlite_demo/db"
	"go_sqlite_demo/helper"
	"go_sqlite_demo/models"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	conn := initializeDatabase()
	defer conn.Close()

	testMessages := helper.CreateTestData()
	for id := range testMessages {
		dbId, _ := db.InsertMessage(conn, testMessages[id])
		fmt.Printf("Message inserted with ID: %d\n", dbId)
	}

	// Initialize channels
	ch1 := make(chan models.Message, 100)
	ch2 := make(chan models.Message, 100)

	//start go routines for initial load and pipelines
	go db.GetRowsAndPutInChannel(ctx, conn, ch1)
	go pipeline1(ctx, ch1, ch2)
	go pipeline2(ctx, ch2)	

	handleGracefulShutdown(cancel, ch1, ch2)

	fmt.Println("Server is shutting down...")
	db.ResetDB(conn)
}

func initializeDatabase() *sql.DB {
	fmt.Println("Initializing the database...")
	conn, err := db.Run()
	if err != nil {
		fmt.Println("error starting db with err: ", err)
		os.Exit(1)
	}
	return conn
}

func handleGracefulShutdown(cancel context.CancelFunc, ch1, ch2 chan models.Message) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nReceived interrupt signal...")
	cancel() // Signal all pipelines to stop

	if completed := waitForCompletionWithTimeout(ch1, ch2); completed {
		fmt.Println("All pipelines completed gracefully")
	} else {
		fmt.Println("Shutdown timeout reached, some messages may be unprocessed")
	}
}

func waitForCompletionWithTimeout(ch1, ch2 chan models.Message) bool {
	shutdownTimeout := time.NewTimer(5 * time.Second)
	defer shutdownTimeout.Stop()

	for { //for loop allows to check multiple times if channes are empty and waits 100ms each iteration
		select {
		case <-shutdownTimeout.C:
			return false
		default:
			// Allow some time for channels to drain
			if len(ch1) == 0 && len(ch2) == 0 {
				return true
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}


func pipeline1(ctx context.Context, in <-chan models.Message, out chan<- models.Message) {
	defer close(out)

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
