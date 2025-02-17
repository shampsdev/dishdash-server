package task

import (
	"net/http"
	"strconv"

	"dashboard.dishdash.ru/pkg/repo"
	"github.com/gin-gonic/gin"
)

// GetTask godoc
// @Summary Get a task by ID
// @Description Retrieve a task from the database
// @Tags tasks
// @Accept json
// @Produce json
// @Schemes http https
// @Param id path string true "Task ID"
// @Success 200 {object} domain.Task
// @Failure 400 "Bad Request"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Security ApiKeyAuth
// @Router /tasks/{id} [get]
func GetTask(taskRepo repo.Task) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		task, err := taskRepo.GetTaskByID(c, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}

		c.JSON(http.StatusOK, task)
	}
}
