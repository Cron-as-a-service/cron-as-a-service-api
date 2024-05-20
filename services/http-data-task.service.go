package services

import (
	"context"
	"encoding/json"
	"fmt"
	"frderoubaix.me/cron-as-a-service/config/database"
	"frderoubaix.me/cron-as-a-service/model"
	"frderoubaix.me/cron-as-a-service/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

func HttpDataTask(task model.CronTask) {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if task.HttpMethod == "GET" {

		// Get Data
		result, err := httpCall(task)
		if err != nil {
			//log
			return
		}

		lastResult, err := getLastTaskResult(task)
		if err != nil {
			//log
			err := addDataInDB(task, result)
			if err != nil {
				TaskLogger(task, "FAILED", fmt.Sprintf(
					"Task id : %s with Method %s cannot upload result in database", task.Id, task.HttpMethod,
				))
				return
			}
			TaskLogger(task, "OK", fmt.Sprintf(
				"Task id : %s with Method %s success with code : %d for the first time", task.Id, task.HttpMethod, result.StatusCode,
			))
			zap.L().Info(
				fmt.Sprintf(
					"Task id : %s with Method %s success with code : %d for the first time", task.Id, task.HttpMethod, result.StatusCode,
				),
			)
			return
		}

		print(lastResult)

		// Insert plat in DB (replace previous call by current call)
		err = addDataInDB(task, result)
		if err != nil {
			TaskLogger(task, "FAILED", fmt.Sprintf(
				"Task id : %s with Method %s cannot upload result in database", task.Id, task.HttpMethod,
			))
			return
		}

		currentResult, err := getLastTaskResult(task)
		if err != nil {
			//log
			return
		}

		print(currentResult)

		//si pas de previous pas de traitement
		// TODO : Treatment
		// Get treatment function
		treatmentFunc := TreatmentFactory(*task.Differential)
		if treatmentFunc == nil {
			//log
			return
		}

		treatmentResult, err := treatmentFunc(lastResult["result"], currentResult["result"], *task.ObjectId, task.Filters)
		if err != nil {
			return
		}

		err = repositories.StoreTreatmentResult(task.Id, treatmentResult)
		if err != nil {
			zap.L().Info(
				fmt.Sprintf(
					"Task id : %s Impossible to store treatment in database", task.Id,
				),
			)
			return
		}

		zap.L().Info(
			fmt.Sprintf(
				"Task id : %s with Method %s success with code : %d", task.Id, task.HttpMethod, result.StatusCode,
			),
		)
	}
}

func httpCall(task model.CronTask) (*http.Response, error) {
	result, err := http.Get(task.TargetUrl)
	if err != nil {
		zap.L().Error(
			fmt.Sprintf(
				"Task id : %s with Method %s failed with code : %d and error : %s", task.Id, task.HttpMethod, result.StatusCode, err.Error(),
			),
		)
		return nil, err
	}

	return result, err
}

func getLastTaskResult(task model.CronTask) (bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"cronTaskId": task.Id}
	var result bson.M
	err := database.TasksResult.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func addDataInDB(task model.CronTask, result *http.Response) error {
	body, err := io.ReadAll(result.Body)
	if err != nil {
		return err
	}

	var jsonData interface{}
	if err := json.Unmarshal(body, &jsonData); err != nil {
		return err
	}

	data := bson.M{
		"cronTaskId": task.Id,
		"timestamp":  time.Now(),
		"result":     jsonData,
	}

	filter := bson.M{"cronTaskId": task.Id}
	update := bson.M{"$set": data}
	opts := options.Update().SetUpsert(true)

	_, err = database.TasksResult.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}
