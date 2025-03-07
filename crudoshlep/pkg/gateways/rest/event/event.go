package event

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/shampsdev/dishdash-server/crudoshlep/pkg/usecase"
)

func SetupHandlers(r *gin.RouterGroup, useCases usecase.Cases) {
	r.POST("/events", CreateEvent(useCases.Event))
}

// CreateEvent godoc
// @Summary Save event
// @Description Save event, you can type any valid json in data field
// @Tags events
// @Accept json
// @Produce json
// @Schemes http https
// @Param event body usecase.SaveEventInput true "Event data"
// @Success 200
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Security ApiKeyAuth
// @Router /events [post]
func CreateEvent(useCase *usecase.Event) gin.HandlerFunc {
	return func(c *gin.Context) {
		var event usecase.SaveEventInput
		if err := c.ShouldBindJSON(&event); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if !json.Valid(event.Data) {
			c.JSON(400, gin.H{"error": "invalid json"})
			return
		}
		err := useCase.SaveEvent(c, event)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Status(200)
	}
}
