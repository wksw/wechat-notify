package controllers

import (
	"github.com/astaxie/beego"
	. "wksw/notify/models/error"
)

type MainController struct {
	beego.Controller
}

type ClusterController struct {
	beego.Controller
}

func (this *MainController) Get() {
	this.Ctx.Output.JSON(NoneErr, true, false)
	return
}

func (this *ClusterController) Info() {

	this.Ctx.Output.JSON(NoneErr, true, false)
	return
}
