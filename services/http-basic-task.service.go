package services

import (
	"fmt"
	"frderoubaix.me/cron-as-a-service/model"
	"go.uber.org/zap"
	"net/http"
)

func BasicTask(task model.CronTask) {
	switch task.HttpMethod {
	case "GET":
		result, err := http.Get(task.TargetUrl)
		if err != nil {
			TaskLogger(task, "ERROR", fmt.Sprintf(
				"Task id : %s with Method %s failed with error : %s", task.Id, task.HttpMethod, err.Error(),
			))
			zap.L().Error(
				fmt.Sprintf(
					"Task id : %s with Method %s failed with error : %s", task.Id, task.HttpMethod, err.Error(),
				),
			)
			return
		}
		TaskLogger(task, "OK", fmt.Sprintf(
			"Task id : %s with Method %s success with code : %d", task.Id, task.HttpMethod, result.StatusCode,
		))
		zap.L().Info(
			fmt.Sprintf(
				"Task id : %s with Method %s success with code : %d", task.Id, task.HttpMethod, result.StatusCode,
			),
		)
		break
	case "POST":
		result, err := http.Post(task.TargetUrl, "application/json", nil)
		if err != nil {
			TaskLogger(task, "ERROR", fmt.Sprintf(
				"Task id : %s with Method %s failed with error : %s", task.Id, task.HttpMethod, err.Error(),
			))
			zap.L().Error(
				fmt.Sprintf(
					"Task id : %s with Method %s failed with error : %s", task.Id, task.HttpMethod, err.Error(),
				),
			)
			return
		}
		TaskLogger(task, "OK", fmt.Sprintf(
			"Task id : %s with Method %s success with code : %d", task.Id, task.HttpMethod, result.StatusCode,
		))
		zap.L().Info(
			fmt.Sprintf(
				"Task id : %s with Method %s success with code : %d", task.Id, task.HttpMethod, result.StatusCode,
			),
		)
		break
	case "PUT":
		result, err := http.NewRequest(http.MethodPut, task.TargetUrl, nil)
		if err != nil {
			TaskLogger(task, "ERROR", fmt.Sprintf(
				"Task id : %s with Method %s failed with error : %s", task.Id, task.HttpMethod, err.Error(),
			))
			zap.L().Error(
				fmt.Sprintf(
					"Task id : %s with Method %s failed with error : %s", task.Id, task.HttpMethod, err.Error(),
				),
			)
			return
		}
		TaskLogger(task, "OK", fmt.Sprintf(
			"Task id : %s with Method %s success with code : %d", task.Id, task.HttpMethod, result.Response.StatusCode,
		))
		zap.L().Info(
			fmt.Sprintf(
				"Task id : %s with Method %s success with code : %d", task.Id, task.HttpMethod, result.Response.StatusCode,
			),
		)
		break
	}
}
