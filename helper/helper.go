package helper

import (
	"go_sqlite_demo/models"
	"time"
)

func createTestData() []models.Message {
	return []models.Message{
		{Severity: 123,
			DescriptionText:  "blablabla",
			ReceivedDateTime: time.Now()},
		{Severity: 12,
			DescriptionText:  "this is ab very long text, i dont really know ",
			ReceivedDateTime: time.Now()},
		{Severity: 8382,
			DescriptionText:  "short text",
			ReceivedDateTime: time.Now()},
		{Severity: 1,
			DescriptionText:  "hack hackers etc",
			ReceivedDateTime: time.Now()},
		{Severity: 4,
			DescriptionText:  "",
			ReceivedDateTime: time.Now()},
		{Severity: 123,
			DescriptionText:  "blablabla",
			ReceivedDateTime: time.Now()},
		{Severity: 12,
			DescriptionText:  "this is ab very long text, i dont really know ",
			ReceivedDateTime: time.Now()},
		{Severity: 8382,
			DescriptionText:  "short text",
			ReceivedDateTime: time.Now()},
		{Severity: 1,
			DescriptionText:  "hack hackers etc",
			ReceivedDateTime: time.Now()},
		{Severity: 4,
			DescriptionText:  "",
			ReceivedDateTime: time.Now()},
	}
}
