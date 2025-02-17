package task

import (
	"net/http"
	"strconv"

	"dashboard.dishdash.ru/pkg/repo"
	"github.com/gin-gonic/gin"
)

// DeleteTask godoc
// @Summary Delete a task
// @Description Delete a task with the given ID from the database
// @Tags tasks
// @Accept json
// @Produce json
// @Schemes http https
// @Param id path string true "Task ID"
// @Success 200
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Security ApiKeyAuth
// @Router /tasks/{id} [delete]
func DeleteTask(taskRepo repo.Task) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = taskRepo.DeleteTask(c, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
	}
}
