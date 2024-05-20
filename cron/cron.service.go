package cron

import (
	"fmt"
	"frderoubaix.me/cron-as-a-service/model"
	"frderoubaix.me/cron-as-a-service/services"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// Cron instance
var cronInstance *cron.Cron

func InitCron() {
	zap.L().Info(fmt.Sprintf("Cron service going to start"))
	cronInstance = cron.New()
	tasks, err := services.GetAllTasks()
	if err != nil {
		zap.L().Error(fmt.Sprintf("Impossible to create task : %v", err))
		return
	}

	for _, task := range tasks {
		id := AddCronTask(task)
		_, err = services.UpdateCronId(task.Id, id)
		if err != nil {
			return
		}
	}

	cronInstance.Start()
	zap.L().Info(fmt.Sprintf("Cron service correctly start"))
}

func AddCronTask(task model.CronTask) cron.EntryID {
	id, cronServiceError := cronInstance.AddFunc(task.Cron, func() { services.ServiceFactory(task) })
	if cronServiceError != nil {
		zap.L().Error(fmt.Sprintf("Impossible to add Task with id : %s to cron service action error : %v", task.Id, cronServiceError))
	}
	return id
}

func DeleteCronTask(task model.CronTask) {
	cronInstance.Remove(task.CronId)
}

func UpdateCronTask(task model.CronTask) {
	DeleteCronTask(task)
	AddCronTask(task)
}
