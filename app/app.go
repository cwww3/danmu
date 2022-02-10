package app

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

var ParamErr = errors.New("参数有误")

type Response struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: "00000",
		Msg:  "Success",
		Data: data,
	})
}

func FailureWithErr(c *gin.Context, err error) {
	c.JSON(http.StatusOK, Response{
		Code: "00001",
		Msg:  err.Error(),
	})
}
