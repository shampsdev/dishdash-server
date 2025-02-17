package task

import (
	"net/http"
	"strconv"

	"dashboard.dishdash.ru/pkg/domain"
	"dashboard.dishdash.ru/pkg/repo"
	"github.com/gin-gonic/gin"
)

// UpdateTask godoc
// @Summary Update a task
// @Description Update a task in the database
// @Tags tasks
// @Accept json
// @Produce json
// @Schemes http https
// @Param id path string true "Task ID"
// @Param task body domain.Task true "Task Data"
// @Success 200 {object} domain.Task
// @Failure 400 "Bad Request"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Security ApiKeyAuth
// @Router /tasks/{id} [put]
func UpdateTask(taskRepo repo.Task) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updatedTask domain.Task
		if err := c.ShouldBindJSON(&updatedTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		updatedTask.ID = id

		task, err := taskRepo.UpdateTask(c, &updatedTask)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, task)
	}
}
