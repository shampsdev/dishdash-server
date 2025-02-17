package task

import (
	"net/http"

	"dashboard.dishdash.ru/pkg/repo"
	"github.com/gin-gonic/gin"
)

// GetAllTasks godoc
// @Summary Get all tasks
// @Description Retrieve all tasks from the database
// @Tags tasks
// @Accept json
// @Produce json
// @Schemes http https
// @Success 200 {array} domain.Task
// @Failure 500 "Internal Server Error"
// @Security ApiKeyAuth
// @Router /tasks [get]
func GetAllTasks(taskRepo repo.Task) gin.HandlerFunc {
	return func(c *gin.Context) {
		tasks, err := taskRepo.GetAllTasks(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, tasks)
	}
}
