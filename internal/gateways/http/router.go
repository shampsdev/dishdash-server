package http

import (
	"dishdash.ru/docs"
	cardHandler "dishdash.ru/internal/gateways/http/handlers/card"
	"dishdash.ru/internal/gateways/http/handlers/swipes"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const allowOrigin = "*"

func setupRouter(router *gin.Engine, wsServer *socketio.Server, useCases UseCases) {
	router.HandleMethodNotAllowed = true
	v1 := router.Group("/api/v1")
	{
		cardHandler.SetupHandlers(v1, useCases.Card)
	}

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Use(allowOriginMiddleware(allowOrigin))
	router.GET("/socket.io/*any", gin.WrapH(wsServer))
	router.POST("/socket.io/*any", gin.WrapH(wsServer))
	swipes.SetupEcho(wsServer)
}
