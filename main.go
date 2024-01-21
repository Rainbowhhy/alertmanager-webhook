package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"golangcode/alertmanager-webhook/conf"
	"golangcode/alertmanager-webhook/model"
	"golangcode/alertmanager-webhook/sender"
	"io"
	"log"
	"os"
	"time"
)

var h bool
var workdir string
var configFile string
var qywechatKey string
var feishuKey string
var dingdingKey string
var logFilePath string
var port string
var host string
var server string
var redisServer string
var redisPort string
var redisPassword string

func initConf() {
	if dir, err := os.Getwd(); err == nil {
		workdir = dir
	}
	var webhookConf conf.Config
	webhookConf = conf.GetConf(workdir, configFile)
	qywechatKey = webhookConf.QyWechatKey
	feishuKey = webhookConf.FeishuKey
	dingdingKey = webhookConf.DingdingKey
	logFilePath = webhookConf.LogFilePath
	port = webhookConf.Port
	host = webhookConf.Host
	server = fmt.Sprintf("%s:%s", host, port)
	redisServer = webhookConf.RedisServer
	redisPort = webhookConf.RedisPort
	redisPassword = webhookConf.RedisPassword
}

func main() {
	flag.BoolVar(&h, "h", false, "display this help and exit")
	flag.StringVar(&configFile, "c", "", "please input webhook configuration file path!")
	flag.Parse()
	if h {
		flag.Usage()
		return
	}
	initConf()
	gin.DisableConsoleColor()
	var f *os.File
	_, err := os.Lstat(logFilePath)
	if err != nil {
		f, _ = os.Create(logFilePath)
	} else {
		f, _ = os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	}
	gin.DefaultWriter = io.MultiWriter(f)
	router := gin.New()

	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// the client access log format
		return fmt.Sprintf("%s - - [%s] \"%s %s %s\" %d %s \"%s\" \"%s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())

	router.POST("/qywechat", func(c *gin.Context) {
		var notification model.Notification
		err := c.BindJSON(&notification)
		if err != nil {
			log.Println("request error:", err.Error())
			return
		}
		sender.SendToQywechat(notification, qywechatKey, redisServer, redisPort, redisPassword)

	})

	router.POST("/feishu", func(c *gin.Context) {
		var notification model.Notification
		err := c.BindJSON(&notification)
		if err != nil {
			log.Println("request error:", err.Error())
			return
		}
		sender.SendToFeishu(notification, feishuKey, redisServer, redisPort, redisPassword)

	})

	router.POST("/dingding", func(c *gin.Context) {
		var notification model.Notification
		err := c.BindJSON(&notification)
		if err != nil {
			log.Println("request error:", err.Error())
			return
		}
		sender.SendToDingding(notification, dingdingKey, redisServer, redisPort, redisPassword)

	})

	log.Println("the Process is Running")
	router.Run(server)
}
