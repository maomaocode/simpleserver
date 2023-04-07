package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type (
	QuestAnswer struct {
		Quest  string `json:"quest"`
		Answer string `json:"answer"`
	}

	UploadFormReq struct {
		ID           string        `json:"id"`
		QuestAnswers []QuestAnswer `json:"quest_answers"`
	}

	UploadFormRes struct {
	}

	CheckIsRegisteredReq struct {
		ID string `json:"id"`
	}

	CheckIsRegisteredRes struct {
		Result bool `json:"result"`
	}
)

type Meta struct {
	ErrCode uint32 `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

type Msg struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

func NewMsg(code uint32, message string, data interface{}) *Msg {
	return &Msg{
		Meta: Meta{
			ErrCode: code,
			ErrMsg:  message,
		},
		Data:    data,
	}
}

func main() {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	r.GET("/api/v1/ping", func(c *gin.Context) {
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

	r.POST("/api/v1/uploadForm", func(context *gin.Context) {
		data, _ := ioutil.ReadAll(context.Request.Body)

		req := UploadFormReq{}
		if err := json.Unmarshal(data, &req); err != nil {
			context.JSON(200, gin.H{
				"msg": NewMsg(1, "invalid request", nil),
			})
			return
		}

		if saver.Exist(req.ID) {
			context.JSON(200, gin.H{
				"msg": NewMsg(2, "already registered", nil)})
			return

		}

		afterMarshal, _ := json.Marshal(req)

		saver.write(req.ID, string(afterMarshal))

		context.JSON(200, gin.H{
			"msg": NewMsg(0, "success", &UploadFormRes{}),
		})
	})

	r.POST("/api/v1/checkIsRegistered", func(context *gin.Context) {
		data, _ := ioutil.ReadAll(context.Request.Body)

		req := CheckIsRegisteredReq{}
		if err := json.Unmarshal(data, &req); err != nil {
			context.JSON(200, gin.H{
				"msg": NewMsg(1, "invalid request", nil),
			})
			return
		}

		context.JSON(200, gin.H{
			"msg": NewMsg(0, "success", &CheckIsRegisteredRes{Result: saver.Exist(req.ID)})})
	})

	go r.Run(":8080")

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		select {
		case sig := <-c:
			switch sig {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("receive exit signal ", sig.String(), ",exit...")
				return
			}
		}
	}

}
