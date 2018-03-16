package main

import (
	"github.com/astaxie/beego"
	"os"
	_ "wksw/notify/filter"
	_ "wksw/notify/models/project"
	. "wksw/notify/models"
	. "wksw/notify/models/amqp"
	_ "wksw/notify/models/error"
	. "wksw/notify/models/job"
	. "wksw/notify/models/job/queue"
	_ "wksw/notify/routers"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	dbinit()
	subscribe(NOTIFY_TYPE)
	err := JobInit()
	if err != nil {
		beego.Error(err)
		os.Exit(1)
	}
	// subscribe(CLEAR_TYPE)
}

func subscribe(tp string) {
	beego.Informational(tp, "subscribe")
	config, err := GetConfig()
	if err != nil {
		beego.Error(err)
		os.Exit(1)
	}

	switch config.Amqp.Amqp {
	case "redis":
		redis, err := NewRedisPool(config.Amqp.ManageHost, config.Amqp.Password, config.Amqp.ManagePort)
		if err != nil {
			beego.Error(err)
			os.Exit(1)
		}
		c := redis.Pool.Get()
		defer c.Close()
		_, err = c.Do("PING")
		if err != nil {
			beego.Error(err)
			os.Exit(1)
		}
		var job Job
		job.Type = tp
		job.Channel = tp
		go redis.Subscribe(&job)
	case "rabbitmq":
		beego.Error("not implement")
		os.Exit(1)
	default:
		beego.Error("Unknown amqp config")
		os.Exit(1)
	}

}

func dbinit() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	database := beego.AppConfig.String("metadata_db::database")
	if database == "" {
		database = "notify:notify@tcp(172.172.0.11:3306)/notify?charset=utf-8"
	}
	orm.RegisterDataBase("default", "mysql", database)
	orm.RunSyncdb("default", false, true)
}

func main() {
	beego.Run()
}
