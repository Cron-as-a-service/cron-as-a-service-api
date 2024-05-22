package routes

import (
	"fmt"
	"frderoubaix.me/cron-as-a-service/cron"
	"frderoubaix.me/cron-as-a-service/model"
	"frderoubaix.me/cron-as-a-service/repositories"
	"frderoubaix.me/cron-as-a-service/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func TaskRoute(router *gin.Engine) {
	//All routes related to tasks comes here
	router.GET("/test", func(context *gin.Context) {
		println("test endpoint trigger")
		context.IndentedJSON(http.StatusOK, "test")
	})

	router.GET("/api/v1/task/user/:userId", func(context *gin.Context) {
		userId := context.Param("userId")
		tasks, err := services.GetByUserId(userId)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, tasks)
	})

	router.POST("/api/v1/task", func(context *gin.Context) {
		var task model.CronTask

		if err := context.ShouldBindJSON(&task); err != nil {
			zap.L().Error(fmt.Sprintf("Impossible to bind request body to CronTask: %v", err))
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		// Generate Id
		task.Id = uuid.New().String()

		// Initialize date fields
		task.CreatedAt = time.Now()
		task.UpdatedAt = time.Now()

		result, err := services.CreateTask(task)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Impossible to create task : %v", err))
			return
		}
		// Add task to cron job
		id := cron.AddCronTask(task)
		_, err = services.UpdateCronId(task.Id, id)
		if err != nil {
			return
		}
		context.JSON(http.StatusCreated, result)
	})

	router.PUT("/api/v1/task/:id", func(context *gin.Context) {
		id := context.Param("id")
		var task model.CronTask

		if err := context.ShouldBindJSON(&task); err != nil {
			zap.L().Error(fmt.Sprintf("Impossible to bind request body to CronTask: %v", err))
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		task.Id = id
		task.UpdatedAt = time.Now()

		result, taskUpdate, err := services.UpdateTask(task)
		cron.UpdateCronTask(*taskUpdate)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Impossible to update task: %v", err))
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		context.JSON(http.StatusOK, result)
	})

	router.GET("/api/v1/task/:taskId/logs", func(context *gin.Context) {
		taskId := context.Param("taskId")
		logs, err := services.GetTaskLogs(taskId)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, logs)
	})

	router.DELETE("/api/v1/task/:id", func(context *gin.Context) {
		id := context.Param("id")

		result, task, err := services.DeleteTask(id)
		cron.DeleteCronTask(*task)
		_, err = services.PurgeTaskLogs(task.Id)
		if err != nil {
			return
		}
		if err != nil {
			zap.L().Error(fmt.Sprintf("Impossible to delete task: %v", err))
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		context.JSON(http.StatusOK, result)
	})

	router.GET("/api/v1/task/:taskId/treatments", func(context *gin.Context) {
		taskId := context.Param("taskId")
		treatments, err := repositories.GetTreatmentResults(taskId)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, treatments)
	})

	router.GET("/api/v1/task/:taskId/run", func(context *gin.Context) {
		taskId := context.Param("taskId")
		task, err := services.GetTaskById(taskId)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		services.RunTask(*task)
		context.Status(http.StatusOK)
	})
}
