package projectcontroller

import (
	"github.com/astaxie/beego"
	. "wksw/notify/models"
	. "wksw/notify/models/error"
	. "wksw/notify/models/db"
	. "wksw/notify/models/project"
	. "wksw/notify/models/keystone"
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"fmt"
	. "wksw/notify/models/notify/publicnummodels"
)

type ProjectController struct {
	beego.Controller
}

/*
{
	"appid": "appid",
	"type": "",
	"describe": "",
	"CustomToken": ""
}
*/

func (this *ProjectController) Add() {
	var ersp ErrResponse
	var p NotifyProject
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &p)
	if err != nil {
		beego.Error(err)
		ersp.ErrCode = 3001
		ersp.ErrMsg = "bad request"
		this.Ctx.Output.JSON(ersp, true, false)
		return
	}
	if p.Appid == "" {
		beego.Error("appid empty in request body")
		ersp.ErrCode = 3002
		ersp.ErrMsg = "appid is must in request body"
		this.Ctx.Output.JSON(ersp, true, false)
		return
	}
	if p.Type != 1 && p.Type != 2 {
		beego.Error("project type not miniprogram or publicnum")
		ersp.ErrCode = 3003
		ersp.ErrMsg = "type must be miniprogram or publicnum"
		this.Ctx.Output.JSON(ersp, true, false)
		return
	}
	o := orm.NewOrm()
	o.Using("default")
	
	tokeninfo, err := GetUserInfo(this.Ctx.Input.Header("X-Auth-Token"))
	if err != nil {
		beego.Error(err)
		ersp.ErrCode = 3008
		ersp.ErrMsg = "get user info fail"
		this.Ctx.Output.JSON(ersp, true, false)
		return
	}
	exist := o.QueryTable("notify_project").Filter("appid", p.Appid).Filter("Userid", tokeninfo.Token.User.Id).Exist()
	beego.Informational("exist=", exist)
	if exist {
		beego.Error("project", p.Appid, "already exist")
		ersp.ErrCode = 3008
		ersp.ErrMsg = "project " + p.Appid + " already exist"
		this.Ctx.Output.JSON(ersp, true, false)
		return
	}
	p.Userid = tokeninfo.Token.User.Id
	_, err = o.Insert(&p)
	if err != nil {
		beego.Error(err)
		ersp.ErrCode = 3004
		ersp.ErrMsg = "add project fail"
		this.Ctx.Output.JSON(ersp, true, false)
		return
	}
	if p.CustomToken != "" {
		config, err := GetConfig()
		if err != nil {
			beego.Error(err)
			this.Ctx.Output.JSON(ConfErr, true, false)
			return
		}

		ssdb, err := NewSsdb(config.Db.WriteHost, config.Db.WritePort)
		if err != nil {
			beego.Error(err)
			this.Ctx.Output.JSON(ConnectDbErr, true, false)
			return
		}
		defer ssdb.Client.Close()

		var req []*CustomToken
		req = append(req, &CustomToken{p.Appid, p.CustomToken})

		for _, v := range req {
			err = v.Validation()
			if err != nil {
				beego.Error(err)
				this.Ctx.Output.JSON(BadRequestErr, true, false)
				return
			}
		}

		err = ssdb.AddCustomToken(req)
		if err != nil {
			beego.Error(err)
			this.Ctx.Output.JSON(AddCustomTokenErr, true, false)
			return
		}
	}
	this.Ctx.Output.JSON(ersp, true, false)
	return
}

func (this *ProjectController) Delete() {
	var ersp ErrResponse
	o := orm.NewOrm()
	o.Using("default")
	appid := this.GetString("appid")
	if appid == "" {
		beego.Error("appid empty in paramater")
		ersp.ErrCode = 3005
		ersp.ErrMsg = "appid paramater not found"
		this.Ctx.Output.JSON(ersp, true, false)
		return
	}
	tokeninfo, err := GetUserInfo(this.Ctx.Input.Header("X-Auth-Token"))
	if err != nil {
		beego.Error(err)
		ersp.ErrCode = 3008
		ersp.ErrMsg = "get user info fail"
		this.Ctx.Output.JSON(ersp, true, false)
		return
	}
	if _, err := o.Raw("DELETE FROM notify_project where appid=? and userid=?", appid, tokeninfo.Token.User.Id).Exec(); err != nil {
		beego.Error(err)
		ersp.ErrCode = 3006
		ersp.ErrMsg = "delete fail"
		this.Ctx.Output.JSON(ersp, true, false)
		return
	}

	this.Ctx.Output.JSON(ersp, true, false)
	return
}

func (this *ProjectController) Get() {
	var ersp ErrResponse
	tokeninfo, err := GetUserInfo(this.Ctx.Input.Header("X-Auth-Token"))
	if err != nil {
		beego.Error(err)
		ersp.ErrCode = 3008
		ersp.ErrMsg = "get user info fail"
		this.Ctx.Output.JSON(ersp, true, false)
		return
	}
	beego.Informational("userinfo=", tokeninfo.Token.User.Id)
	o := orm.NewOrm()
	o.Using("default")
	var maps []orm.Params
	_, err = o.QueryTable("NotifyProject").Filter("Userid", tokeninfo.Token.User.Id).Values(&maps)
	if err != nil {
		beego.Error(err)
		ersp.ErrCode = 3007
		ersp.ErrMsg = "get projects fail"
		this.Ctx.Output.JSON(ersp, true, false)
		return

	}
	config, err := GetConfig()
	if err != nil {
		beego.Error(err)
		this.Ctx.Output.JSON(ConfErr, true, false)
		return
	}
	ssdb, err := NewSsdb(config.Db.WriteHost, config.Db.WritePort)
	if err != nil {
		beego.Error(err)
		this.Ctx.Output.JSON(ConnectDbErr, true, false)
		return
	}
	defer ssdb.Client.Close()
	for index, m := range(maps){
		typee := m["Type"]
		if typee.(int64) ==  1{
			currentUser, _ := ssdb.Zsize(fmt.Sprintf(UNKNOWN_TEMPLATE_ID, m["Appid"])) 
			maps[index]["currentUser"] = currentUser
		} else if typee.(int64) == 2 {
			var user User
			user.AppId = m["Appid"].(string)
			user.Secret = m["Secret"].(string)
			token, _ := ssdb.GetSetToken(&user)
			u, err := GetUsers(token, "")
			if err != nil {
				beego.Error(err)
				maps[index]["currentUser"] = 0
			} else {
				maps[index]["currentUser"] = u.Total
			}

		}


		m["Secret"] = "******"
	}
	this.Ctx.Output.JSON(maps, true, false)
	return
}

func (this *ProjectController) GetByAppId() {
	var ersp ErrResponse
	appid := this.Ctx.Input.Param(":appid")
	tokeninfo, err := GetUserInfo(this.Ctx.Input.Header("X-Auth-Token"))
	if err != nil {
		beego.Error(err)
		ersp.ErrCode = 3008
		ersp.ErrMsg = "get user info fail"
		this.Ctx.Output.JSON(ersp, true, false)
		return
	}

	o := orm.NewOrm()
	o.Using("default")
	var maps []orm.Params
	_, err = o.QueryTable("NotifyProject").Filter("Userid", tokeninfo.Token.User.Id).Filter("Appid", appid).Values(&maps)
	if err != nil {
		beego.Error(err)
		ersp.ErrCode = 3007
		ersp.ErrMsg = "get project fail"
		this.Ctx.Output.JSON(ersp, true, false)
		return
	}
	for _, m := range(maps) {
		m["Secret"] = "******"
	}
	this.Ctx.Output.JSON(maps, true, false)
	return

}

