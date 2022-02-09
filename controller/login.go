package controller

import (
	"danmu/app"
	"danmu/model"
	"danmu/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func RegisterLogin(router *gin.Engine) {
	router.POST("/login", func(c *gin.Context) {
		args := new(struct {
			User string `json:"user"`
		})
		err := c.Bind(args)
		if err != nil {
			app.FailureWithErr(c, err)
			return
		}
		user := &model.User{
			Name: args.User,
		}
		id, err := repository.GetMySQLRepository().SaveUser(user)
		if err != nil {
			app.FailureWithErr(c, err)
			return
		}
		c.JSON(http.StatusOK, strconv.FormatUint(uint64(id), 10))
	})
}
