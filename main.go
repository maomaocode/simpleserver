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

type (
	QuestAnswer struct {
		Quest  string `json:"quest"`
		Answer string `json:"answer"`
	}

	UploadFormReq struct {
		Address      string        `json:"address"`
		QuestAnswers []QuestAnswer `json:"quest_answers"`
	}

	UploadFormRes struct {
	}

	CheckIsRegisteredReq struct {
		Address string `json:"address"`
	}

	CheckIsRegisteredRes struct {
		Result bool `json:"result"`
	}
)

type Msg struct {
	Code    uint32      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewMsg(code uint32, message string, data interface{}) *Msg {
	return &Msg{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func main() {
	r := gin.Default()
	r.Use(Cors())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": NewMsg(0, "success", "pong"),
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

		req := UploadFormReq{}
		if err := json.Unmarshal(data, &req); err != nil {
			context.JSON(200, gin.H{
				"msg": NewMsg(1, "invalid request", nil),
			})
			return
		}

		if saver.Exist(req.Address) {
			context.JSON(200, gin.H{
				"msg": NewMsg(2, "already registered", nil)})
			return

		}

		afterMarshal, _ := json.Marshal(req)

		saver.write(req.Address, string(afterMarshal))

		context.JSON(200, gin.H{
			"msg": NewMsg(0, "success", &UploadFormRes{}),
		})
	})

	r.POST("/checkIsRegistered", func(context *gin.Context) {
		data, _ := ioutil.ReadAll(context.Request.Body)

		req := CheckIsRegisteredReq{}
		if err := json.Unmarshal(data, &req); err != nil {
			context.JSON(200, gin.H{
				"msg": NewMsg(1, "invalid request", nil),
			})
			return
		}

		context.JSON(200, gin.H{
			"msg": NewMsg(0, "success", &CheckIsRegisteredRes{Result: saver.Exist(req.Address)})})
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
