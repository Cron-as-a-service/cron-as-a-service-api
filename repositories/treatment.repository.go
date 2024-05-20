package repositories

import (
	"context"
	"frderoubaix.me/cron-as-a-service/config/database"
	"frderoubaix.me/cron-as-a-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func StoreTreatmentResult(taskId string, result []interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	treatment := model.TaskTreatment{
		ID:        primitive.NewObjectID(),
		TaskId:    taskId,
		Timestamp: time.Now(),
		Result:    result,
	}

	// Insert the new treatment result
	_, err := database.TasksTreatment.InsertOne(ctx, treatment)
	if err != nil {
		return err
	}

	// Ensure only the 5 latest treatments are kept
	filter := bson.M{"taskid": taskId}
	opts := options.Find().SetSort(bson.D{{"timestamp", 1}})
	cursor, err := database.TasksTreatment.Find(ctx, filter, opts)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var treatments []model.TaskTreatment
	if err := cursor.All(ctx, &treatments); err != nil {
		return err
	}

	if len(treatments) > 5 {
		// Remove the oldest treatments
		for i := 0; i < len(treatments)-5; i++ {
			_, err := database.TasksTreatment.DeleteOne(ctx, bson.M{"_id": treatments[i].ID})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func GetTreatmentResults(taskId string) ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var treatments []model.TaskTreatment
	filter := bson.M{"taskid": taskId}
	cursor, err := database.TasksTreatment.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &treatments); err != nil {
		return nil, err
	}

	results := make([]map[string]interface{}, len(treatments))
	for i, treatment := range treatments {
		result := make(map[string]interface{})
		result["ID"] = treatment.ID.Hex()
		result["TaskId"] = treatment.TaskId
		result["Timestamp"] = treatment.Timestamp

		// Convert the result from primitive.A to JSON object
		convertedResult := []map[string]interface{}{}
		if items, ok := treatment.Result.(primitive.A); ok {
			for _, item := range items {
				convertedItem := make(map[string]interface{})
				if kvPairs, ok := item.(primitive.D); ok {
					for _, kv := range kvPairs {
						key := kv.Key
						value := kv.Value
						convertedItem[key] = value
					}
				}
				convertedResult = append(convertedResult, convertedItem)
			}
		}

		result["Result"] = convertedResult
		results[i] = result
	}

	return results, nil
}
