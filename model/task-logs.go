package model

import "time"

type TaskLogs struct {
	TaskId    string
	Status    string
	Timestamp time.Time
	Message   string
}
