package task

import (
	"dashboard.dishdash.ru/cmd/config"
	"dashboard.dishdash.ru/pkg/gateways/http/middlewares"
	"dashboard.dishdash.ru/pkg/repo"
	"github.com/gin-gonic/gin"
)

func SetupHandlers(r *gin.RouterGroup, taskRepo repo.Task) {
	taskGroup := r.Group("tasks")
	taskGroup.Use(middlewares.ApiTokenAuth(config.C.Auth.ApiToken))

	taskGroup.GET("", GetAllTasks(taskRepo))
	taskGroup.GET(":id", GetTask(taskRepo))
	taskGroup.POST("", SaveTask(taskRepo))
	taskGroup.PUT(":id", UpdateTask(taskRepo))
	taskGroup.DELETE(":id", DeleteTask(taskRepo))
}
