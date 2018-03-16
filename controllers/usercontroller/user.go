package usercontroller

import (
	"encoding/json"
	"github.com/astaxie/beego"
	. "wksw/notify/models"
	. "wksw/notify/models/db"
	. "wksw/notify/models/error"
	. "wksw/notify/models/keystone"
	"fmt"
	"time"
)

type UserController struct {
	beego.Controller
}

type login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type project struct {
	AppId string `json:"appid"`
	Type string `json:"type"`
	TokenUrl string `json:"tokenurl"`
	Templates []string `json:"templates"`
}

/*
Post /custometoken
[
	{
		"appid": "",
		"tokenurl": ""
	}
]
*/

func (this *UserController) Post() {
	var req []*CustomToken
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	if err != nil {
		beego.Error(err)
		this.Ctx.Output.JSON(BadRequestErr, true, false)
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
	this.Ctx.Output.JSON(NoneErr, true, false)
	return

}

func (this *UserController) Delete() {
	var req []string
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	if err != nil {
		beego.Error(err)
		this.Ctx.Output.JSON(BadRequestErr, true, false)
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
	err = ssdb.DeleteCustomToken(req)
	if err != nil {
		beego.Error(err)
		this.Ctx.Output.JSON(DeleteCustomTokenErr, true, false)
		return
	}
	this.Ctx.Output.JSON(NoneErr, true, false)
	return
}

type user struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Email string `json:"email"`
	Description string `json:"description"`
}

type response struct {
	Token string `json:"token"`
	Expires_in int64 `json:"expires_in"`
}

func (this *UserController) Login() {
	var u user
	var resp ErrResponse
	var t K_UserToken

	config, _ := GetConfig()

	err := json.Unmarshal(this.Ctx.Input.RequestBody, &u)
	if err != nil {
		beego.Error(err)
		resp.ErrCode = 1
		resp.ErrMsg = fmt.Sprintf("Bad Request [%s]", err)
		this.Ctx.Output.JSON(resp, true, false)
		return
	}
	t.At.Idt.Methods = append(t.At.Idt.Methods, "password")
	t.At.Idt.Pwd.Us.Name = u.UserName
	t.At.Idt.Pwd.Us.Dm.Name = config.User_domain_name
	t.At.Idt.Pwd.Us.Password = u.Password
	token, err := t.GetUserToken()
	if err != nil {
		beego.Error(err)
		resp.ErrCode = 2
		resp.ErrMsg = fmt.Sprintf("login fail [%s]", err)
		this.Ctx.Output.JSON(resp, true, false)
		return
	}
	var t_resp response
	t_resp.Token = token
	tokeninfo, err := GetUserInfo(token)
	if err != nil {
		beego.Error(err)
		resp.ErrCode = 3008
		resp.ErrMsg = "get user info fail"
		this.Ctx.Output.JSON(resp, true, false)
		return
	}
	expires_at := tokeninfo.Token.Expires_at
	now := time.Now().Unix()
	format := "2006-01-02T15:04:05.000000Z"
	expires_at_timestamp, _ := time.Parse(format, expires_at)
	t_timestamp := expires_at_timestamp.Unix()
	t_resp.Expires_in = t_timestamp - now 
	this.Ctx.Output.JSON(t_resp, true, false)
	return 
}

func (this *UserController) Registe() {
	var u user 
	var resp ErrResponse
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &u)
	if err != nil {
		beego.Error(err)
		resp.ErrCode = 3
		resp.ErrMsg = fmt.Sprintf("bad request [%s]", err)
		this.Ctx.Output.JSON(resp, true, false)
		return
	}

	var uc UserCreate
	uc.Uinfo.Name = u.UserName
	uc.Uinfo.Password = u.Password
	uc.Uinfo.Default_project_id = DEFAULT_USER_PROJECT_NAME
	uc.Uinfo.Domain_id = DEFAULT_USER_DOMAIN_NAME
	uc.Uinfo.Email = u.Email
	uc.Uinfo.Description = u.Description

	err = uc.Create()
	if err != nil {
		beego.Error(err)
		resp.ErrCode = 4
		resp.ErrMsg = fmt.Sprintf("registe fail [%s]", err)
		this.Ctx.Output.JSON(resp, true, false)
		return 
	}

	this.Ctx.Output.JSON(resp, true, false)
	return
}


