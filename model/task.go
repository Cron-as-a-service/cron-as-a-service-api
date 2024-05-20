package model

import (
	"github.com/robfig/cron/v3"
	"time"
)

type CronTask struct {
	Id           string
	Name         string
	UserId       string
	Cron         string
	CronId       cron.EntryID
	FunctionType string
	HttpMethod   string
	TargetUrl    string

	// Attributs optionnels
	Filters      []string
	ObjectId     *string
	Differential *string

	// Champs de date
	CreatedAt time.Time
	UpdatedAt time.Time
}
