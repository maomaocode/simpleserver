package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

type WRequest struct {
	Address string `json:"address"`
	Data    string `json:"data"`
}

type WResponse struct {
	Code    uint32      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewResponse(code uint32, message string, data interface{}) *WResponse {
	return &WResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": NewResponse(0, "success", "pong"),
		})
	})

	saver := NewSaver("./server/data")
	defer saver.Close()

	_, err := saver.GetFile()
	if err != nil {
		panic(err)
	}

	r.POST("/uploadForm", func(context *gin.Context) {
		data, _ := ioutil.ReadAll(context.Request.Body)

		req := WRequest{}
		if err := json.Unmarshal(data, &req); err != nil {
			context.JSON(200, gin.H{
				"msg": NewResponse(1, "invalid request", nil),
			})
			return
		}

		saver.write(fmt.Sprintf("%s,%s\n", req.Address, req.Data))

		context.JSON(200, gin.H{
			"msg": NewResponse(0, "success", nil),
		})
	})

	go r.Run()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		select {
		case sig := <-c:
			switch sig {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				// 这里做一些清理操作或者输出相关说明，比如 断开数据库连接
				fmt.Println("receive exit signal ", sig.String(), ",exit...")
				return
			}
		}
	}

}
