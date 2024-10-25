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
	defer cancel()

	conn := initializeDatabase()
	defer conn.Close()

	testMessages := helper.CreateTestData()
	for id := range testMessages {
		dbId, _ := db.InsertMessage(conn, testMessages[id])
		fmt.Printf("Message inserted with ID: %d\n", dbId)
	}
	// Initialize channels and start pipelines
	ch1 := make(chan models.Message, 100)
	ch2 := make(chan models.Message, 100)

	startPipelines(ctx, conn, ch1, ch2)

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

func startPipelines(ctx context.Context, conn *sql.DB, ch1, ch2 chan models.Message) {
	go func() {
		defer close(ch1)
		db.GetRowsAndPutInChannel(ctx, conn, ch1)
	}()

	go func() {
		defer close(ch2)
		pipeline1(ctx, ch1, ch2)
	}()

	go func() {
		pipeline2(ctx, ch2)
	}()
}

func handleGracefulShutdown(cancel context.CancelFunc, ch1, ch2 chan models.Message) {
	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nReceived interrupt signal...")
	cancel() // Signal all pipelines to stop

	// Wait briefly to allow pipelines to finish processing
	if completed := waitForCompletionWithTimeout(ch1, ch2); completed {
		fmt.Println("All pipelines completed gracefully")
	} else {
		fmt.Println("Shutdown timeout reached, some messages may be unprocessed")
	}
}

func waitForCompletionWithTimeout(ch1, ch2 chan models.Message) bool {
	shutdownTimeout := time.NewTimer(5 * time.Second)
	defer shutdownTimeout.Stop()

	select {
	case <-shutdownTimeout.C:
		return false
	default:
		// Allow some time for channels to drain
		time.Sleep(100 * time.Millisecond)
		return len(ch1) == 0 && len(ch2) == 0
	}
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
