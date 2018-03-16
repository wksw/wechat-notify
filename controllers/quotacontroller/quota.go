package quotacontroller


import (
	"github.com/astaxie/beego"
	. "wksw/notify/models/error"
)


type QuotaController struct {
	beego.Controller
}

/*
GET /project
response boy:
{
	"project_limit": 10
}
*/

func (this *QuotaController) Get() {
	var ersp ErrResponse
	this.Ctx.Output.JSON(ersp, true, false)
	return
}

/*
POST /project

request body:
{
	"project_limit": 10
}
*/

func (this *QuotaController) Update() {
	var ersp ErrResponse
	this.Ctx.Output.JSON(ersp, true, false)
	return 
}