package services

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Traitement des nouveaux éléments
func NewTreatment(lastResult, currentResult interface{}, attributeId string, filters []string) ([]interface{}, error) {
	lastData, ok := lastResult.(primitive.A)
	if !ok {
		return nil, fmt.Errorf("lastResult is not a primitive.A")
	}

	currentData, ok := currentResult.(primitive.A)
	if !ok {
		return nil, fmt.Errorf("currentResult is not a primitive.A")
	}

	newItems := make([]interface{}, 0)
	for _, item := range currentData {
		itemMap, ok := item.(primitive.M)
		if !ok {
			return nil, fmt.Errorf("item in currentResult is not a map[string]interface{}")
		}
		lastItemFound := false
		for _, lastItem := range lastData {
			lastItemMap, ok := lastItem.(primitive.M)
			if !ok {
				return nil, fmt.Errorf("item in lastResult is not a map[string]interface{}")
			}
			if lastItemMap[attributeId] == itemMap[attributeId] {
				lastItemFound = true
				break
			}
		}
		if !lastItemFound {
			newItems = append(newItems, itemMap)
		}
	}

	return newItems, nil
}

// Traitement des éléments supprimés
func DeleteTreatment(lastResult, currentResult interface{}, attributeId string, filters []string) ([]interface{}, error) {
	lastData, ok := lastResult.(primitive.A)
	if !ok {
		return nil, fmt.Errorf("lastResult is not a primitive.A")
	}

	currentData, ok := currentResult.(primitive.A)
	if !ok {
		return nil, fmt.Errorf("currentResult is not a primitive.A")
	}

	deletedItems := make([]interface{}, 0)
	for _, lastItem := range lastData {
		lastItemMap, ok := lastItem.(primitive.M)
		if !ok {
			return nil, fmt.Errorf("item in lastResult is not a map[string]interface{}")
		}
		currentItemFound := false
		for _, item := range currentData {
			itemMap, ok := item.(primitive.M)
			if !ok {
				return nil, fmt.Errorf("item in currentResult is not a map[string]interface{}")
			}
			if itemMap[attributeId] == lastItemMap[attributeId] {
				currentItemFound = true
				break
			}
		}
		if !currentItemFound {
			deletedItems = append(deletedItems, lastItemMap)
		}
	}

	return deletedItems, nil
}

// Traitement des éléments mis à jour
func UpdateTreatment(lastResult, currentResult interface{}, attributeId string, filters []string) ([]interface{}, error) {
	lastData, ok := lastResult.(primitive.A)
	if !ok {
		return nil, fmt.Errorf("lastResult is not a primitive.A")
	}

	currentData, ok := currentResult.(primitive.A)
	if !ok {
		return nil, fmt.Errorf("currentResult is not a primitive.A")
	}

	updatedItems := make([]interface{}, 0)
	for _, item := range currentData {
		itemMap, ok := item.(primitive.M)
		if !ok {
			return nil, fmt.Errorf("item in currentResult is not a map[string]interface{}")
		}
		for _, lastItem := range lastData {
			lastItemMap, ok := lastItem.(primitive.M)
			if !ok {
				return nil, fmt.Errorf("item in lastResult is not a map[string]interface{}")
			}
			if lastItemMap[attributeId] == itemMap[attributeId] {
				updated := false
				for _, filter := range filters {
					if lastItemMap[filter] != itemMap[filter] {
						updated = true
						break
					}
				}
				if updated {
					updatedItems = append(updatedItems, itemMap)
				}
				break
			}
		}
	}

	return updatedItems, nil
}
