package services

import "frderoubaix.me/cron-as-a-service/model"

func ServiceFactory(task model.CronTask) {
	switch task.FunctionType {
	case "differential":
		HttpDataTask(task)
	default:
		BasicTask(task)
	}
}
