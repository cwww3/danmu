package controller

import (
	"danmu/app"
	"danmu/model"
	"danmu/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os/exec"
)

func RegisterRoom(router *gin.Engine) {
	// 获取直播间历史消息
	router.GET("/:room/message", func(c *gin.Context) {
		room := c.Param("room")
		msgList, err := repository.GetMySQLRepository().GetRoomMessageList(room)
		if err != nil {
			app.FailureWithErr(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"msgList": msgList})
	})
	// 开通直播间
	router.POST("/room", func(c *gin.Context) {
		room := new(model.Room)
		err := c.Bind(room)
		if err != nil {
			app.FailureWithErr(c, err)
			return
		}
		room.Url = "rtmp://110.42.134.163/" + room.Name
		_, err = repository.GetMySQLRepository().SaveRoom(room)
		if err != nil {
			app.FailureWithErr(c, err)
			return
		}
		// 创建nginx配置
		err = nginxConfig(room.Name)
		if err != nil {
			app.FailureWithErr(c, err)
			return
		}

		app.Success(c, room)
	})
	// 获取直播间
	router.GET("/:user/room", func(c *gin.Context) {
		user := c.Param("user")
		if len(user) == 0 {
			app.FailureWithErr(c, app.ParamErr)
		}
		room, err := repository.GetMySQLRepository().GetRoomByUser(user)
		if err != nil {
			app.FailureWithErr(c, err)
			return
		}
		app.Success(c, room)
	})
	// 获取开播的直播间
	router.GET("/liveon/room", func(c *gin.Context) {
		app.Success(c, repository.GetMemoryRepository().GetList())
	})
	// 开始直播
	router.POST("/:user/on", func(c *gin.Context) {
		user := c.Param("user")
		if len(user) == 0 {
			app.FailureWithErr(c, app.ParamErr)
			return
		}
		room, err := repository.GetMySQLRepository().GetRoomByUser(user)
		if err != nil {
			app.FailureWithErr(c, err)
			return
		}
		repository.GetMemoryRepository().Save(user, room)
		app.Success(c, nil)
	})
	// 关闭直播
	router.POST("/:user/off", func(c *gin.Context) {
		user := c.Param("user")
		if len(user) == 0 {
			app.FailureWithErr(c, app.ParamErr)
			return
		}
		repository.GetMemoryRepository().Delete(user)
		app.Success(c, nil)
	})
}

func nginxConfig(room string) error {
	//创建nginx文件
	text := "\\                application %v {\\n                    live on;\\n                    gop_cache on; #打开 GOP 缓存，减少首屏等待时间\\n                }"
	text = fmt.Sprintf(text, room)
	text = fmt.Sprintf("3a %v", text)
	c := exec.Command("sed", "-i", text, "/etc/nginx/conf.d/live/cwww3.conf")
	data, err := c.CombinedOutput()
	if err != nil {
		log.Println("update nginx file failed err", string(data))
		return err
	}

	cc := exec.Command("service", "nginx", "reload")
	data, err = cc.CombinedOutput()
	if err != nil {
		log.Println("nginx reload failed err", string(data))
	}
	return err
}
