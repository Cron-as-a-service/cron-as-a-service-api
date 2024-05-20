package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TaskTreatment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	TaskId    string             `bson:"taskid"`
	Timestamp time.Time          `bson:"timestamp"`
	Result    interface{}        `bson:"result"`
}
