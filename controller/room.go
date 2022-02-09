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
	router.GET("/:room/message", func(c *gin.Context) {
		room := c.Param("room")
		msgList, err := repository.GetMySQLRepository().GetRoomMessageList(room)
		if err != nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}
		c.JSON(http.StatusOK, gin.H{"msgList": msgList})
	})
	router.POST("/room", func(c *gin.Context) {
		room := new(model.Room)
		err := c.Bind(room)
		if err != nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}
		//room.Url = "http://110.42.134.163/live?app=cwww"
		room.Url = "rtmp://110.42.134.163/" + room.Name
		_, err = repository.GetMySQLRepository().SaveRoom(room)
		if err != nil {
			app.Failure(c)
			return
		}
		// nginx
		err = nginxConfig(room.Name)
		if err != nil {
			app.Failure(c)
			return
		}

		app.Success(c, room)
	})
	// TODO  替代 get room
	router.GET("/:user/room", func(c *gin.Context) {
		user := c.Param("user")
		room, err := repository.GetMySQLRepository().GetRoomByUser(user)
		if err != nil {
			app.Failure(c)
			return
		}
		app.Success(c, room)
	})
	router.GET("/liveon/room", func(c *gin.Context) {
		app.Success(c, repository.GetMemoryRepository().GetList())
	})

	// TODO  替代 post /liveon/:user
	router.POST("/:user/on", func(c *gin.Context) {
		user := c.Param("user")
		if len(user) == 0 {
			app.Failure(c)
			return
		}
		room, err := repository.GetMySQLRepository().GetRoomByUser(user)
		if err != nil {
			app.Failure(c)
			return
		}
		repository.GetMemoryRepository().Save(user, room)
		app.Success(c, nil)
	})
	router.POST("/liveon/:user", func(c *gin.Context) {
		user := c.Param("user")
		if len(user) == 0 {
			app.Failure(c)
			return
		}
		room, err := repository.GetMySQLRepository().GetRoomByUser(user)
		if err != nil {
			app.Failure(c)
			return
		}
		repository.GetMemoryRepository().Save(user, room)
		app.Success(c, nil)
	})
	// TODO  替代 post /liveoff/:user
	router.POST("/:user/off", func(c *gin.Context) {
		user := c.Param("user")
		if len(user) == 0 {
			app.Failure(c)
			return
		}
		repository.GetMemoryRepository().Delete(user)
		app.Success(c, nil)
	})
	router.POST("/liveoff/:user", func(c *gin.Context) {
		user := c.Param("user")
		if len(user) == 0 {
			app.Failure(c)
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
