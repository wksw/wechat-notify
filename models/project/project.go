package project


import (
	"github.com/astaxie/beego/orm"
	
)
type NotifyProject struct {
	Id int `json:"id"`
	Userid string `json:"userid"`
	Appid string `json:"appid"`
	Type int `json:"type"`
	Describe string `json:"describe"`
	CustomToken string `json:"customToken"`
	Secret string `json:"secret"`
}

func init() {
	orm.RegisterModel(new(NotifyProject))
}