package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
)

type Config struct {
	AppName    string
	HttpPort   int
	RunMode    string
	MaxPush    int
	MaxJob     int
	MaxHistory int
	Amqp       Amqp
	Db         Db
	Auth_url string 
	Admin_username string
	Admin_password string
	Admin_domain_name string
	Admin_project_name string
	Admin_project_domain_name string
	User_project_name  string
	User_domain_name string
	User_enable bool
	Notify_server string
	Max_project int 
	Wechat Wechat
}

type Amqp struct {
	Amqp       string
	UserName   string
	Password   string
	ManageHost string
	ManagePort int
	PublicHost string
	PublicPort int
}

type Db struct {
	UserName  string
	Password  string
	WriteHost string
	WritePort int
	ReadHost  string
	ReadPort  int
}

type Wechat struct {
	AppId string 
	Secret string 
}

func GetConfig() (*Config, error) {
	var config Config
	config.AppName = beego.AppConfig.String("appname")
	port, err := beego.AppConfig.Int("httpport")
	if err != nil {
		port = 8080
	}
	config.HttpPort = port
	config.RunMode = beego.AppConfig.String("runmode")

	maxpush, err := beego.AppConfig.Int("max_push")
	if err != nil {
		maxpush = DEFAULT_MAX_PUSH
	}
	if maxpush == 0 {
		maxpush = DEFAULT_MAX_PUSH
	}
	config.MaxPush = maxpush

	maxjob, err := beego.AppConfig.Int("static::max_job")
	if err != nil {
		maxjob = DEFAULT_MAX_JOB
	}
	if maxjob == 0 {
		maxjob = DEFAULT_MAX_JOB
	}
	config.MaxJob = maxjob

	maxhistory, _ := beego.AppConfig.Int("static::max_history")
	if maxhistory == 0 {
		maxhistory = DEFAULT_MAX_HISTORY
	}
	config.MaxHistory = maxhistory

	config.Amqp.Amqp = beego.AppConfig.String("amqp::amqp")
	if config.Amqp.Amqp == "" {
		config.Amqp.Amqp = "redis"
	}
	if config.Amqp.Amqp != "redis" && config.Amqp.Amqp != "rabbitmq" {
		return nil, errors.New(fmt.Sprintf("parse config faile, unknow amqp %s", config.Amqp.Amqp))
	}
	config.Amqp.UserName = beego.AppConfig.String("amqp::username")
	config.Amqp.Password = beego.AppConfig.String("amqp::password")
	config.Amqp.ManageHost = beego.AppConfig.String("amqp::manage_host")
	if config.Amqp.ManageHost == "" {
		config.Amqp.ManageHost = "localhost"
	}

	port, err = beego.AppConfig.Int("amqp::manage_port")
	if err != nil {
		if config.Amqp.Amqp == "redis" {
			port = 6379
		} else if config.Amqp.Amqp == "rabbitmq" {
			port = 5672

		}
	}
	config.Amqp.ManagePort = port

	config.Amqp.PublicHost = beego.AppConfig.String("amqp::public_host")
	if config.Amqp.PublicHost == "" {
		config.Amqp.PublicHost = config.Amqp.ManageHost
	}

	port, err = beego.AppConfig.Int("amqp::public_port")
	if err != nil {
		port = config.Amqp.ManagePort
	}
	config.Amqp.PublicPort = port

	config.Db.UserName = beego.AppConfig.String("db::username")
	config.Db.Password = beego.AppConfig.String("db::password")
	config.Db.WriteHost = beego.AppConfig.String("db::write_host")
	if config.Db.WriteHost == "" {
		config.Db.WriteHost = "localhost"
	}
	port, err = beego.AppConfig.Int("db::write_port")
	if err != nil {
		port = 8888
	}
	config.Db.WritePort = port

	config.Db.ReadHost = beego.AppConfig.String("db::read_host")
	if config.Db.ReadHost == "" {
		config.Db.ReadHost = config.Db.WriteHost
	}
	port, err = beego.AppConfig.Int("db::read_port")
	if err != nil {
		port = config.Db.WritePort
	}
	config.Db.ReadPort = port


	auth_url := beego.AppConfig.String("identify::auth_url")
	if auth_url == "" {
		auth_url = DEFAULT_AUTH_URL
	}
	config.Auth_url = auth_url

	admin_username := beego.AppConfig.String("identify::admin_username")
	if admin_username == "" {
		admin_username = DEFAULT_ADMIN_USERNAME
	}
	config.Admin_username = admin_username

	admin_password := beego.AppConfig.String("identify::admin_password")
	if admin_password == "" {
		admin_password = DEFAULT_ADMIN_PASSWORD
	}
	config.Admin_password = admin_password

	admin_domain_name := beego.AppConfig.String("identify::admin_domain_name")
	if admin_domain_name == "" {
		admin_domain_name = DEFAULT_ADMIN_DOAMIN_NAME
	}
	config.Admin_domain_name = admin_domain_name

	admin_project_name := beego.AppConfig.String("identify::admin_project_name")
	if admin_project_name == "" {
		admin_project_name = DEFAULT_ADMIN_PROJECT_NAME
	}
	config.Admin_project_name = admin_project_name

	admin_project_domain_name := beego.AppConfig.String("identify::admin_project_domain_name")
	if admin_project_domain_name == "" {
		admin_project_domain_name = DEFAULT_ADMIN_PROJECT_DOMAIN_NAME
	}
	config.Admin_project_domain_name = admin_project_domain_name

	user_project_name := beego.AppConfig.String("identify::user_project_name")
	if user_project_name == "" {
		user_project_name = DEFAULT_USER_PROJECT_NAME
	}
	config.User_project_name = user_project_name

	user_domain_name := beego.AppConfig.String("identify::user_domain_name")
	if user_domain_name == "" {
		user_domain_name = DEFAULT_USER_DOMAIN_NAME
	}
	config.User_domain_name = user_domain_name

	user_enable, err := beego.AppConfig.Bool("identify::user_enable")
	if err != nil {
		user_enable = false
	}
	config.User_enable = user_enable

	max_project, err := beego.AppConfig.Int("max_project")
	if max_project == 0 {
		max_project = DEFAULT_MAX_PROJECT
	}

	config.Max_project = max_project

	appid := beego.AppConfig.String("wechat::appid")
	if appid == "" {
		return &config, errors.New("appid is must")
	}
	config.Wechat.AppId = appid

	secret := beego.AppConfig.String("wechat::secret")
	if secret == "" {
		return &config, errors.New("secret is must")
	}
	config.Wechat.Secret = secret

	

	return &config, nil

}
