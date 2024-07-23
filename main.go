package main

import (
	"log"
	"map_lan_ip/config"
	"map_lan_ip/task"
	"os"
)

var LogFile string = "map_lan_ip.log"

func main() {
	c, err := config.ReadConfig()
	if err != nil {
		log.Fatal("读取配置文件失败:", err.Error())
	}
	logFile, err := os.OpenFile(c.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("打开日志文件失败:", err.Error())
	}
	defer logFile.Close()

	// 日志写入文件
	log.SetOutput(logFile)

	task.MapLanIp()
}
