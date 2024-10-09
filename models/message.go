package models

import "time"

type Message struct{
    Severity int
    DescriptionText string
    ReceivedDateTime time.Time
}