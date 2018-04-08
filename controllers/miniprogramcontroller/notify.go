package miniprogramcontroller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	. "wksw/notify/models"
	. "wksw/notify/models/db"
	. "wksw/notify/models/error"
	"wksw/notify/models/job"
	"wksw/notify/models/job/queue"
	. "wksw/notify/models/keystone"
	. "wksw/notify/models/notify"
	. "wksw/notify/models/notify/miniprogrammodels"
	"sort"
)

type TemplateNotify struct {
	beego.Controller
}
type state struct {
	Total int `json:"total"`
	// Hst map[int]History `json:"history"`
	Hst []History `json:"history"`
}


func (this *TemplateNotify) GetState() {
	var st state
	config, err := GetConfig()
	if err != nil {
		beego.Error(err)
		this.Ctx.Output.JSON(ConfErr, true, false)
		return
	}
	appid := this.Ctx.Input.Header("AppId")
	templateid := this.GetString("templateid")
	if templateid == "" {
		beego.Error(MissingTemplateidErr.ErrMsg)
		this.Ctx.Output.JSON(MissingTemplateidErr, true, false)
		return
	}
	ssdb, err := NewSsdb(config.Db.WriteHost, config.Db.WritePort)
	if err != nil {
		beego.Error(err)
		this.Ctx.Output.JSON(ConnectDbErr, true, false)
		return
	}
	defer ssdb.Client.Close()
	size, _ := ssdb.Hsize(fmt.Sprintf(HISTORY, appid, templateid))
	hs, _ := ssdb.Hscan(fmt.Sprintf(HISTORY, appid, templateid), "", size)
	var tmpIndex []int
	for _, h := range(hs) {
		id := h.OpenId 
		d, _ := strconv.Atoi(id)
		tmpIndex = append(tmpIndex, d)
	}
	sort.Ints(tmpIndex)

	for _, key := range(tmpIndex) {
		for _, h := range(hs) {
			id := h.OpenId
			d, _ := strconv.Atoi(id)
			value := h.ErrMsg
			var data History
			json.Unmarshal([]byte(value), &data)
			if key == d {
				st.Hst = append(st.Hst, data)
				break
			}
		}
	}
	total, err := ssdb.Zsize(fmt.Sprintf(UNKNOWN_TEMPLATE_ID, appid))
	st.Total = total
	this.Ctx.Output.JSON(st, true, false)
	return
}

func (this *TemplateNotify) GetTemplates() {
	config, err := GetConfig()
	if err != nil {
		beego.Error(err)
		this.Ctx.Output.JSON(ConfErr, true, false)
		return
	}
	appid := this.Ctx.Input.Header("AppId")
	ssdb, err := NewSsdb(config.Db.WriteHost, config.Db.WritePort)
	if err != nil {
		beego.Error(err)
		this.Ctx.Output.JSON(ConnectDbErr, true, false)
		return
	}
	defer ssdb.Client.Close()
	size, _ := ssdb.Zsize(fmt.Sprintf(NOTIFY_TEMPLATES, appid))
	templates, _ := ssdb.Zscan(fmt.Sprintf(NOTIFY_TEMPLATES, appid), "", size)
	this.Ctx.Output.JSON(templates, true, false)
	return
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

	var req MiniprogramTemplateNotify
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	if err != nil {
		beego.Error(err)
		this.Ctx.Output.JSON(BadRequestErr, true, false)
		return
	}
	err = req.Validation()
	if err != nil {
		beego.Error(err)
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
	job.Ntf.Type = MINIPROGRAM_NOTIFY
	job.Ntf.Metadata.AppId = appid
	job.Ntf.Metadata.Secret = secret
	job.Ntf.Data.Mtn = &req
	var que queue.Queue
	que.Jb = &job
	err = que.Push()
	if err != nil {
		beego.Error(err)
		this.Ctx.Output.JSON(JobQueueErr, true, false)
		return
	}
	go Produce(appid, req.TemplateId)

	this.Ctx.Output.JSON(NoneErr, true, false)
	return
}

func (this *TemplateNotify) PostFormid() {
	var ersp ErrResponse
	config, err := GetConfig()
	if err != nil {
		beego.Error(err)
		this.Ctx.Output.JSON(ConfErr, true, false)
		return
	}

	appid := this.Ctx.Input.Header("AppId")
	wxagent := this.Ctx.Input.Header("X-Wechat-Agent")

	var req PostFormid
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	if err != nil {
		beego.Error(err)
		this.Ctx.Output.JSON(BadRequestErr, true, false)
		return
	}

	ssdb, err := NewSsdb(config.Db.WriteHost, config.Db.WritePort)
	if err != nil {
		beego.Error(err)
		this.Ctx.Output.JSON(ConnectDbErr, true, false)
		return
	}
	defer ssdb.Client.Close()

	if wxagent != "" {
		tokeninfo, err := GetUserInfo(this.Ctx.Input.Header("X-Auth-Token"))
		if err != nil {
			beego.Error(err)
			ersp.ErrCode = 3008
			ersp.ErrMsg = "get user info fail"
			this.Ctx.Output.JSON(ersp, true, false)
			return
		}
		if len(req.Data) >0 {
			req.Data[0].OpenId = tokeninfo.Token.User.Name
			appid = config.Wechat.AppId
		}
	}

	err = req.Validation()
	if err != nil {
		beego.Error(err)
		this.Ctx.Output.JSON(RequestValidErr, true, false)
		return
	}
	err = ssdb.StoreFormId(&req, appid)
	if err != nil {
		beego.Error(err)
		this.Ctx.Output.JSON(PostFormIdErr, true, false)
		return
	}
	this.Ctx.Output.JSON(NoneErr, true, false)
	return
}



func (this *TemplateNotify) GetNotifystate() {
	config, err := GetConfig()
	if err != nil {
		beego.Error(err)
		this.Ctx.Output.JSON(ConfErr, true, false)
		return
	}

	appid := this.Ctx.Input.Header("AppId")
	current_page, err := this.GetInt("current_page")
	if err != nil {
		current_page = 1
	}
	if current_page == 0 {
		current_page = 1
	}
	perpage, err := this.GetInt("perpage")
	if err != nil {
		perpage = 5
	}
	if perpage == 0 {
		perpage = 5
	}
	templateid := this.GetString("templateid")
	if templateid == "" {
		beego.Error(MissingTemplateidErr.ErrMsg)
		this.Ctx.Output.JSON(MissingTemplateidErr, true, false)
		return
	}

	ssdb, err := NewSsdb(config.Db.WriteHost, config.Db.WritePort)
	if err != nil {
		beego.Error(err)
		this.Ctx.Output.JSON(ConnectDbErr, true, false)
		return
	}
	defer ssdb.Client.Close()

	var ns NotifyState

	total_str, _ := ssdb.Get(fmt.Sprintf(PRODUCE_TOTAL, appid, templateid))
	if total_str != "" {
		total, err := strconv.Atoi(total_str)
		if err == nil {
			ns.Total = total
		}
	}

	status, _ := ssdb.Get(fmt.Sprintf(JOB_STATUS, appid, templateid))
	if status != "" {
		ns.Status = status
	} else {
		ns.Status = "fail"
	}

	queue, _ := ssdb.Qsize(fmt.Sprintf(JOB_QUEUE, appid, templateid))
	ns.Queue = queue

	fail_total_page := 0
	invalid_total_page := 0

	fail, _ := ssdb.Hsize(fmt.Sprintf(NOTIFY_FAIL, appid, templateid))
	ns.Fail = fail
	if fail != 0 {
		if perpage < fail {
			if fail%perpage != 0 {
				fail_total_page = fail/perpage + 1
			} else {
				fail_total_page = fail / perpage
			}
		}
		var fail_arry []*FailDetail
		start_name := ""
		for i := 0; i < current_page; i++ {
			fail_arry, _ = ssdb.Hscan(fmt.Sprintf(NOTIFY_FAIL, appid, templateid), start_name, perpage)
			if len(fail_arry) == 0 {
				break
			}
			start_name = fail_arry[len(fail_arry)-1].OpenId
		}
		ns.FailDetail = fail_arry

	}
	ns.Page.FailTotalPage = fail_total_page

	invalid, _ := ssdb.Zsize(fmt.Sprintf(NOTIFY_INVALID, appid, templateid))
	ns.Invalid = invalid
	if invalid != 0 {
		if perpage < invalid {
			if invalid%perpage != 0 {
				invalid_total_page = invalid/perpage + 1
			} else {
				invalid_total_page = invalid / perpage
			}
		}
		var invalid_arry []string
		start_name := ""
		for i := 0; i < current_page; i++ {
			invalid_arry, _ = ssdb.Zscan(fmt.Sprintf(NOTIFY_INVALID, appid, templateid), start_name, perpage)
			if len(invalid_arry) == 0 {
				break
			}
			start_name = invalid_arry[len(invalid_arry)-1]
		}
		ns.InvalidDetail = invalid_arry
	}
	ns.Page.InvalidTotalPage = invalid_total_page
	success, _ := ssdb.Zsize(fmt.Sprintf(NOTIFY_SUCCESS, appid, templateid))
	ns.Success = success

	var start_time int64
	timestart_str, _ := ssdb.Get(fmt.Sprintf(JOB_START_TIME, appid, templateid))
	if timestart_str != "" {
		d, err := strconv.ParseInt(timestart_str, 10, 0)
		if err == nil {
			start_time = d
		}
	}
	var end_time int64
	timeend_str, _ := ssdb.Get(fmt.Sprintf(JOB_END_TIME, appid, templateid))
	if timeend_str != "" {
		d, err := strconv.ParseInt(timeend_str, 10, 0)
		if err == nil {
			end_time = d
		}
	}
	if start_time != 0 && end_time != 0 {
		ns.Duration = end_time - start_time
	}
	this.Ctx.Output.JSON(ns, true, false)
	return
}
