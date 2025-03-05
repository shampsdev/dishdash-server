package task

import (
	"net/http"

	"dashboard.dishdash.ru/pkg/domain"
	"dashboard.dishdash.ru/pkg/repo"
	"github.com/gin-gonic/gin"
)

// SaveTask godoc
// @Summary Create a new task
// @Description Create a new task in the database
// @Tags tasks
// @Accept json
// @Produce json
// @Schemes http https
// @Param task body domain.Task true "Task Data"
// @Success 201 {object} domain.Task
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Security ApiKeyAuth
// @Router /tasks [post]
func SaveTask(taskRepo repo.Task) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newTask domain.Task
		if err := c.ShouldBindJSON(&newTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id, err := taskRepo.CreateTask(c, &newTask)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		newTask.ID = id
		c.JSON(http.StatusCreated, newTask)
	}
}
