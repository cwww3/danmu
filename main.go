package main

import (
	"context"
	"danmu/repository"
	"danmu/ws"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	loadConfig()
	go ws.WebsocketManager.Start()
	//go ws.WebsocketManager.SendService()
	go ws.WebsocketManager.SendService()
	go ws.WebsocketManager.SendGroupService()
	//go ws.WebsocketManager.SendGroupService()
	//go ws.WebsocketManager.SendAllService()
	//go ws.WebsocketManager.SendAllService()
	//go ws.TestSendGroup()
	//go ws.TestSendAll()

	repository.SetUpMessageRepository()

	router := gin.Default()
	//router.LoadHTMLGlob("web/*")

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Gin Server")
	})
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	router.GET("/:room/message", func(c *gin.Context) {
		room := c.Param("room")
		msgList,err := repository.GetMessageRepository().GetRoomMessageList(room)
		if err != nil {
			c.JSON(http.StatusBadRequest,nil)
			return
		}
		c.JSON(http.StatusOK,gin.H{"msgList":msgList})
	})

	wsGroup := router.Group("/ws")
	{
		wsGroup.GET("/:channel", ws.WebsocketManager.WsClient)
	}

	srv := &http.Server{
		Addr:    ":8888",
		Handler: router,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server Start Error: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown Error:", err)
	}
	log.Println("Server Shutdown")
}

func loadConfig() {
	viper.AddConfigPath("configs/")
	viper.SetConfigName("settings-prod")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
}
