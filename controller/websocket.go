package controller

import (
	"danmu/ws"
	"github.com/gin-gonic/gin"
)

func RegisterWebSocket(router *gin.Engine) {
	wsGroup := router.Group("/ws")
	{
		wsGroup.GET("/:channel", ws.WebsocketManager.WsClient)
	}
}
