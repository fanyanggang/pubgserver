package main

import (
	pubgsql "github.com/fanyanggang/BackendPlatform/fsql"
	"github.com/wonderivan/logger"
	game "pubgserver/common"
	log "pubgserver/log"
	"pubgserver/model"
)

func main() {

	var pubg, financial pubgsql.SQLGroupConfig

	if model.CONTROL {
		pubg = pubgsql.SQLGroupConfig{
			Name:   "pubg",
			Master: "root:pubg123456@tcp(localhost:3306)/pubg?charset=utf8&max_idle=100&max_active=500&max_lifetime_sec=1800",
			Slaves: []string{"root:pubg123456@tcp(localhost:3306)/pubg?charset=utf8&max_idle=100&max_active=500&max_lifetime_sec=1800"},
		}

		financial = pubgsql.SQLGroupConfig{
			Name:   "financial",
			Master: "root:pubg123456@tcp(localhost:3306)/financial?charset=utf8&max_idle=100&max_active=500&max_lifetime_sec=1800",
			Slaves: []string{"root:pubg123456@tcp(localhost:3306)/financial?charset=utf8&max_idle=100&max_active=500&max_lifetime_sec=1800"},
		}
		game.InitRedis("127.0.0.1:6379")
	} else {
		pubg = pubgsql.SQLGroupConfig{
			Name:   "pubg",
			Master: "pubg_rw:BcPIzg42B/qko9akLJFi5H0cDn4=@tcp(rummy-allinone.c9jqvxqnoxoi.ap-south-1.rds.amazonaws.com:3306)/pubg?charset=utf8&max_idle=100&max_active=500&max_lifetime_sec=1800",
			Slaves: []string{"admin:pubg.2020@tcp(pubg-db.cv3gvs0tnix0.ap-south-1.rds.amazonaws.com:3306)/pubg?charset=utf8&max_idle=100&max_active=500&max_lifetime_sec=1800"},
		}

		financial = pubgsql.SQLGroupConfig{
			Name:   "financial",
			Master: "pubg_rw:BcPIzg42B/qko9akLJFi5H0cDn4=@tcp(rummy-allinone.c9jqvxqnoxoi.ap-south-1.rds.amazonaws.com:3306)/financial?charset=utf8&max_idle=100&max_active=500&max_lifetime_sec=1800",
			Slaves: []string{"admin:pubg.2020@tcp(pubg-db.cv3gvs0tnix0.ap-south-1.rds.amazonaws.com:3306)/financial?charset=utf8&max_idle=100&max_active=500&max_lifetime_sec=1800"},
		}
		game.InitRedis("rummy-redis.oigiga.ng.0001.aps1.cache.amazonaws.com:6379")
	}

	var conf []pubgsql.SQLGroupConfig = []pubgsql.SQLGroupConfig{pubg, financial}
	err := pubgsql.InitSQLClient(conf)
	if err != nil {
		logger.Error("InitSQLClient err:%v", err)
	}

	log.InitLog()
	game.InitService()

	//go game.CheckTables()
}
