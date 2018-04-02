package routers

import (
	"github.com/astaxie/beego"
	"wksw/notify/controllers"
	"wksw/notify/controllers/miniprogramcontroller"
	"wksw/notify/controllers/publicnumcontroller"
	"wksw/notify/controllers/usercontroller"
	"wksw/notify/controllers/projectcontroller"
	"wksw/notify/controllers/templatecontroller"
	"wksw/notify/controllers/quotacontroller"
)

func init() {
	//health check
	beego.Router("/", &controllers.MainController{})

	beego.Router("/clusterinfo", &controllers.ClusterController{}, "get:Info")

	beego.Router("/miniprogram/templatenotify", &miniprogramcontroller.TemplateNotify{}, "post:Notify")
	beego.Router("/miniprogram/templatenotify/postformid", &miniprogramcontroller.TemplateNotify{}, "post:PostFormid")

	beego.Router("/templatenotify/notifystate", &miniprogramcontroller.TemplateNotify{}, "get:GetNotifystate")
	beego.Router("/templatenotify/state", &miniprogramcontroller.TemplateNotify{}, "get:GetState")
	beego.Router("/templatenotify/templates",&miniprogramcontroller.TemplateNotify{},  "get:GetTemplates")

	beego.Router("/publicnum/templatenotify", &publicnumcontroller.TemplateNotify{}, "post:Notify")

	beego.Router("/custometoken", &usercontroller.UserController{}, "post:Post;delete:Delete")

	beego.Router("/login", &usercontroller.UserController{}, "post:Login")
	beego.Router("/wxlogin", &usercontroller.UserController{}, "post:WxLogin")
    beego.Router("/registe", &usercontroller.UserController{}, "post:Registe")
    beego.Router("/project", &projectcontroller.ProjectController{}, "post:Add;delete:Delete;get:Get")
    beego.Router("/project/:appid", &projectcontroller.ProjectController{}, "get:GetByAppId")
    beego.Router("/template", &templatecontroller.TemplateController{}, "get:GetTemplatelist")
    beego.Router("/quota", &quotacontroller.QuotaController{}, "get:Get;post:Update")
}
