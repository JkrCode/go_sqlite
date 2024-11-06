package pipeline

import (
	"bufio"
	"context"
	"fmt"
	"go_sqlite_demo/models"
	"os"
)

func Pipeline1(ctx context.Context, in <-chan models.Message, out chan<- models.Message) {
	defer close(out)

	matchingLines, err := loadMatchingLines("pipeline/example.txt")

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-in:
			if !ok {
				return
			}
			fmt.Println("processing message 1: ", msg)

			mask(scanner, msg, out)

			select {
			case <-ctx.Done():
				return
			case out <- msg:
				// Message forwarded successfully
			}
		}
	}
}

func loadMatchingLines(s string) {
	file, err := os.Open(s)
	if err != nil {
		fmt.Printf("error loading file from os: %s", err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

}

func mask(scanner *bufio.Scanner, msg models.Message, out chan<- models.Message) {
	for {
		if scanner.Text() == msg.DescriptionText {
			out <- msg
		}
		scanner.Scan()
	}
}

func Pipeline2(ctx context.Context, in <-chan models.Message) {
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
