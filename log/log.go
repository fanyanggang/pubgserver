package pubglog

import (
	"github.com/wonderivan/logger"
	"pubgserver/model"
	"log"
	"os"
)


func init() {
	InitLog()
}

func InitLog() {
	// 日志配置
	err := logger.SetLogger(`{"Console": {"level": "DEBG"}}`)

	// 通过配置文件配置
	var logConf string
	if model.CONTROL{
		logConf = "/Users/fengqingyang/data/app/pubgserver/conf/log_test.json"
	}else {
		logConf = "/data/app/pubgserver/conf/log.json"
	}

	_, err = os.Stat(logConf)
	if err != nil{
		if !os.IsExist(err) {
			log.Printf("init log IsExist false:%v", err)
			logConf = "./conf/log.json"
		}
	}

	log.Print("init log succ")
	err = logger.SetLogger(logConf)
}
