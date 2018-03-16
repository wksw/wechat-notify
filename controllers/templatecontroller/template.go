package templatecontroller

import (
	"github.com/astaxie/beego"
	. "wksw/notify/models/error"
	. "wksw/notify/models"
	. "wksw/notify/models/project"
	. "wksw/notify/models/db"
	"github.com/astaxie/beego/orm"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
	"bytes"
)

type TemplateController struct {
	beego.Controller
}

type templateRequest struct {
	OffSet int `json:"offset"`
	Count int `json:"count"`
}

type miniprogramTemplateResponse struct {
	ErrCode int `json:"errcode"`
	ErrMsg string `json:"errmsg"`
	Tl []miniprogramTemplateList `json:"list"`
}

type miniprogramTemplateList struct {
	TemplateId string `json:"template_id"`
	Title string `json:"title"`
	Content string `json:"content"`
	Example string `json:"example"`
}

type publicnumTemplateResponse struct {
	Tl []publicnumTemplateList `json:"template_list"`
}

type publicnumTemplateList struct {
	TemplateId string `json:"template_id"`
	Title string `json:"title"`
	Primary_Industry string `json:"primary_industry"`
	Deputy_Industry string `json:"deputy_industry"`
	Content string `json:"content"`
	Example string `json:"example"`
}

/*
GET /template?appid=
*/
func (this *TemplateController) GetTemplatelist() {
	var ersp ErrResponse
	appid := this.GetString("appid")
	if appid == "" {
		beego.Error("appid is empty")
		ersp.ErrCode = 4001
		ersp.ErrMsg = "appid is must"
		this.Ctx.Output.JSON(ersp, true, false)
		return
	}

	o := orm.NewOrm()
	o.Using("default")

	var pj NotifyProject
	p := new(NotifyProject)
	qs := o.QueryTable(p)
	err := qs.Filter("appid", appid).One(&pj)
	if err != nil {
		beego.Error(err)
		ersp.ErrCode = 4002
		ersp.ErrMsg = "get " + appid + "fail"
		this.Ctx.Output.JSON(ersp, true, false)
		return
	}
	beego.Informational(pj)
	
	config, err := GetConfig()
	if err != nil {
		beego.Error(err)
		this.Ctx.Output.JSON(ConfErr, true, false)
		return
	}
	s, err := NewSsdb(config.Db.WriteHost, config.Db.WritePort)
	if err != nil {
		beego.Error(err)
		this.Ctx.Output.JSON(ConnectDbErr, true, false)
		return
	}
	defer s.Client.Close()

	var user User
	user.AppId = pj.Appid
	user.Secret = pj.Secret
	token, err := s.GetSetToken(&user)
	if err != nil {
		beego.Error(err)
		ersp.ErrCode = 4009
		ersp.ErrMsg = "get token fail"
		this.Ctx.Output.JSON(ersp, true, false)
		return

	}

	beego.Informational("token=", token)
	if pj.Type == 1 {
		template_list, err := getMiniprogramTemplateList(token)
		if err != nil {
			beego.Error(err)
			ersp.ErrCode = 4006
			ersp.ErrMsg = "get template list fail"
			this.Ctx.Output.JSON(ersp, true, false)
			return
		}
		this.Ctx.Output.JSON(template_list, true, false)
	} else if pj.Type == 2 {
		template_list, err := getTempPublicnumlateList(token)
		if err != nil {
			beego.Error(err)
			ersp.ErrCode = 4007
			ersp.ErrMsg = "get template list fail"
			this.Ctx.Output.JSON(ersp, true, false)
			return
		}
		this.Ctx.Output.JSON(template_list, true, false)
	} else {
		beego.Error("Unknown project type")
		ersp.ErrCode = 4008
		ersp.ErrMsg = "unknown project type"
		this.Ctx.Output.JSON(ersp, true, false)
		return
	}

}

func getMiniprogramTemplateList(token string) (*miniprogramTemplateResponse, error) {
	var trsp miniprogramTemplateResponse
	u, _ := url.Parse(MINIPROGRAM_TEMPLATE_API)
	q := u.Query()
	q.Set("access_token", token)
	u.RawQuery = q.Encode()

	var trq templateRequest
	trq.OffSet = 0
	trq.Count = 10

	body_str, err := json.Marshal(&trq)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer([]byte(body_str))

	resp, err := http.Post(u.String(), "application_json:charset=utf-8", body)

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.Unmarshal([]byte(result), &trsp)
	if err != nil {
		return nil, err
	}
	return &trsp, nil

}

func getTempPublicnumlateList(token string) (*publicnumTemplateResponse, error) {
	var trsp publicnumTemplateResponse
	u, _ := url.Parse(PUBLICNUM_TEMPLATE_API)
	q := u.Query()
	q.Set("access_token", token)
	u.RawQuery = q.Encode()


	resp, err := http.Get(u.String())

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.Unmarshal([]byte(result), &trsp)
	if err != nil {
		return nil, err
	}
	return &trsp, nil

}