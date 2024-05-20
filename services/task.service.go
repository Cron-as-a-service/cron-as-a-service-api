package services

import (
	"context"
	"fmt"
	"frderoubaix.me/cron-as-a-service/config/database"
	"frderoubaix.me/cron-as-a-service/model"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	_ "net/http"
	"time"
)

func CreateTask(task model.CronTask) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := database.Tasks.InsertOne(ctx, task)
	if err != nil {
		return result, err
	}
	TaskLogger(task, "CREATED", fmt.Sprintf("Task with id : %s was correctly created", task.Id))
	return result, err
}

func GetAllTasks() ([]model.CronTask, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var tasks []model.CronTask
	filter := bson.M{}
	result, err := database.Tasks.Find(ctx, filter)
	if err != nil {
		return tasks, err
	}
	if err := result.All(context.Background(), &tasks); err != nil {
		return tasks, err
	}
	return tasks, err
}

func GetByUserId(userId string) ([]model.CronTask, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var tasks []model.CronTask
	filter := bson.M{"userid": userId}
	result, err := database.Tasks.Find(ctx, filter)
	if err != nil {
		return tasks, err
	}
	if err := result.All(context.Background(), &tasks); err != nil {
		return tasks, err
	}
	return tasks, err
}

func UpdateTask(task model.CronTask) (*mongo.UpdateResult, *model.CronTask, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"id": task.Id}
	update := bson.M{
		"$set": bson.M{
			"name":         task.Name,
			"userid":       task.UserId,
			"cron":         task.Cron,
			"functiontype": task.FunctionType,
			"httpmethod":   task.HttpMethod,
			"targeturl":    task.TargetUrl,
			"filters":      task.Filters,
			"objectid":     task.ObjectId,
			"differential": task.Differential,
			"updatedat":    time.Now(),
		},
	}

	result, err := database.Tasks.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, nil, err
	}
	TaskLogger(task, "UPDATED", fmt.Sprintf("Task with id : %s was correctly updated", task.Id))

	// Obtenir la tâche avant de la supprimer
	taskUpdate, err := GetTaskById(task.Id)
	if err != nil {
		return nil, nil, err
	}

	return result, taskUpdate, err
}

func DeleteTask(id string) (*mongo.DeleteResult, *model.CronTask, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	task, err := GetTaskById(id)
	if err != nil {
		return nil, nil, err
	}

	filter := bson.M{"id": id}
	result, err := database.Tasks.DeleteOne(ctx, filter)
	if err != nil {
		return result, nil, err
	}

	return result, task, nil
}

// Méthode pour mettre à jour uniquement le champ CronId
func UpdateCronId(taskId string, cronId cron.EntryID) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"id": taskId}
	update := bson.M{
		"$set": bson.M{
			"cronid": cronId,
		},
	}

	result, err := database.Tasks.UpdateOne(ctx, filter, update)
	if err != nil {
		return result, fmt.Errorf("could not update CronId: %v", err)
	}
	return result, nil
}

func GetTaskById(id string) (*model.CronTask, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var task model.CronTask
	filter := bson.M{"id": id}
	err := database.Tasks.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}
