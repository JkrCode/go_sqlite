package helper

import (
	"go_sqlite_demo/models"
	"time"
)

func CreateTestData() []models.Message {
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
			DescriptionText:  "blablabla2",
			ReceivedDateTime: time.Now()},
		{Severity: 12,
			DescriptionText:  "this is ab very long text, i dont really know 2",
			ReceivedDateTime: time.Now()},
		{Severity: 8382,
			DescriptionText:  "short text2",
			ReceivedDateTime: time.Now()},
		{Severity: 1,
			DescriptionText:  "hack hackers etc2",
			ReceivedDateTime: time.Now()},
		{Severity: 4,
			DescriptionText:  "2",
			ReceivedDateTime: time.Now()},
	}
}
