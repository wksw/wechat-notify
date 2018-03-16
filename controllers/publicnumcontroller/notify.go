package publicnumcontroller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	. "wksw/notify/models"
	. "wksw/notify/models/db"
	. "wksw/notify/models/error"
	"wksw/notify/models/job"
	"wksw/notify/models/job/queue"
	. "wksw/notify/models/notify"
	. "wksw/notify/models/notify/publicnummodels"
)

type TemplateNotify struct {
	beego.Controller
}

func (this *TemplateNotify) Notify() {
	config, err := GetConfig()
	if err != nil {
		beego.Error(err)
		this.Ctx.Output.JSON(ConfErr, true, false)
		return
	}

	appid := this.Ctx.Input.Header("AppId")
	secret := this.Ctx.Input.Header("secret")

	var req PublicNumTemplateNotify
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	if err != nil {
		beego.Error(err)
		this.Ctx.Output.JSON(BadRequestErr, true, false)
		return
	}
	err = req.Validation()
	if err != nil {
		beego.Error()
		this.Ctx.Output.JSON(TemplateidEmptyErr, true, false)
		return
	}

	ssdb, err := NewSsdb(config.Db.WriteHost, config.Db.WritePort)
	if err != nil {
		beego.Error(err)
		this.Ctx.Output.JSON(ConnectDbErr, true, false)
		return
	}
	defer ssdb.Client.Close()
	size, _ := ssdb.Qsize(fmt.Sprintf(JOB_QUEUE, appid, req.TemplateId))
	if size != 0 {
		this.Ctx.Output.JSON(JobRunningErr, true, false)
		return
	}

	access_token_bs64 := this.Ctx.Input.Header("Access-Token")
	if access_token_bs64 != "" {
		access_token, err := base64.StdEncoding.DecodeString(access_token_bs64)
		if err != nil {
			beego.Error(err)
			this.Ctx.Output.JSON(TokenParseErr, true, false)
			return
		}
		if access_token != nil {
			var token Token
			err = json.Unmarshal([]byte(access_token), &token)
			if err != nil {
				beego.Error(err)
				this.Ctx.Output.JSON(TokenParseErr, true, false)
				return
			}
			ssdb.Set(fmt.Sprintf(TOKEN, appid), token.Access_Token, token.Expires_In)
		}
	}

	err = ssdb.TemplateNotifyInit(appid, req.TemplateId)
	if err != nil {
		beego.Error(err)
		this.Ctx.Output.JSON(TemplateNotifyInitErr, true, false)
		return
	}
	var job job.Job
	job.Id = fmt.Sprintf("%s_%s", appid, req.TemplateId)
	job.Type = NOTIFY_TYPE
	job.Channel = NOTIFY_TYPE
	job.Ntf.Type = PUBLICNUM_NOTIFY
	job.Ntf.Metadata.AppId = appid
	job.Ntf.Metadata.Secret = secret
	job.Ntf.Data.Ptn = &req
	var que queue.Queue
	que.Jb = &job
	err = que.Push()
	if err != nil {
		beego.Error(err)
		this.Ctx.Output.JSON(JobQueueErr, true, false)
		return
	}
	var user User
	user.AppId = appid
	user.Secret = secret
	go Produce(&user, req.TemplateId)
	this.Ctx.Output.JSON(NoneErr, true, false)
	return
}

func (this *TemplateNotify) GetNotifystate() {
	this.Ctx.Output.JSON(NoneErr, true, false)
	return

}
