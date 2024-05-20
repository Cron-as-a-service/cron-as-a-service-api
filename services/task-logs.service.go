package services

import (
	"context"
	"frderoubaix.me/cron-as-a-service/config/database"
	"frderoubaix.me/cron-as-a-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func TaskLogger(task model.CronTask, status string, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	taskLog := model.TaskLogs{
		TaskId:    task.Id,
		Status:    status,
		Timestamp: time.Now(),
		Message:   message,
	}
	_, err := database.TasksLogs.InsertOne(ctx, taskLog)
	if err != nil {
		// log
		return
	}
}

func GetTaskLogs(taskId string) ([]model.TaskLogs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var logs []model.TaskLogs
	filter := bson.M{"taskid": taskId}
	cursor, err := database.TasksLogs.Find(ctx, filter)
	if err != nil {
		return logs, err
	}
	if err = cursor.All(ctx, &logs); err != nil {
		return logs, err
	}

	return logs, nil
}

func PurgeTaskLogs(taskId string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"taskid": taskId}
	result, err := database.TasksLogs.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
