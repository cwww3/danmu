package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var SuccessResponse = Response{
	Code: "00000",
	Msg:  "Success",
}

var FailureResponse = Response{
	Code: "00001",
	Msg:  "Failure",
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: "00000",
		Msg:  "Success",
		Data: data,
	})
}

func Failure(c *gin.Context) {
	c.JSON(http.StatusOK, FailureResponse)
}

func FailureWithErr(c *gin.Context, err error) {
	c.JSON(http.StatusOK, Response{
		Code: "00001",
		Msg:  err.Error(),
	})
}
