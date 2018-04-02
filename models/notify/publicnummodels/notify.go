package publicnummodels

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"time"
	. "wksw/notify/models"
	. "wksw/notify/models/db"
	. "wksw/notify/models/error"
	. "wksw/notify/models/notify"
)

type job struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Ntf     Notify `json:"notify"`
}

type templateNotify struct {
	ToUser     string                  `json:"touser"`
	TemplateId string                  `json:"template_id"`
	Url        string                  `json:"url"`
	Mp         Miniprogram            `json:"miniprogram"`
	Data       map[string]TemplateData `json:"data"`
}

type users struct {
	Total       int    `json:"total"`
	Count       int    `json:"count"`
	Data        openid `json:"data"`
	Next_openid string `json:next_openid`
}

type openid struct {
	Openid []string `json:"openid"`
}

func Run(n *Notify) error {
	beego.Informational("publicnum notify start")
	cpunum := runtime.NumCPU()
	for i := 0; i < cpunum; i++ {
		go notify(*n)
	}
	return nil
}

func Produce(user *User, templateid string) error {
	beego.Informational("publicnum start push openid into job queue")
	config, err := GetConfig()
	if err != nil {
		beego.Error(err)
		return err
	}

	s, err := NewSsdb(config.Db.WriteHost, config.Db.WritePort)
	if err != nil {
		beego.Error(err)
		return err
	}
	defer s.Client.Close()

	err = s.Zset(fmt.Sprintf(NOTIFY_TEMPLATES, user.AppId), templateid)
	if err != nil {
		beego.Error(err)
		return err
	}

	next_openid := ""
	for {
		token, err := s.GetSetToken(user)
		if err != nil {
			beego.Error(err)
			return err
		}
		u, err := GetUsers(token, next_openid)
		if err != nil {
			beego.Error("get users from remote server fail", err)
			break
		}
		if len(u.Data.Openid) != 0 {
			s.Set(fmt.Sprintf(PRODUCE_TOTAL, user.AppId, templateid), u.Total, INVALID_TTL)

			err = s.Qpush(fmt.Sprintf(JOB_QUEUE, user.AppId), u.Data.Openid)
			if err != nil {
				beego.Error(err)
				break
			}
		}
		if u.Next_openid == "" {
			break
		}
		next_openid = u.Next_openid
	}

	err = s.Set(fmt.Sprintf(PRODUCE_FINISH, user.AppId, templateid), "ok", INVALID_TTL)
	if err != nil {
		beego.Error(err)
		return err
	}
	beego.Informational("publicnum finish push job into queue", fmt.Sprintf(JOB_QUEUE, user.AppId))
	return nil
}

func notify(n Notify) error {
	config, err := GetConfig()
	if err != nil {
		beego.Error(err)
		return err
	}

	s, err := NewSsdb(config.Db.WriteHost, config.Db.WritePort)
	if err != nil {
		beego.Error(err)
		return err
	}
	defer s.Client.Close()

	var tn templateNotify
	tn.TemplateId = n.Data.Ptn.TemplateId
	tn.Url = n.Data.Ptn.Url
	tn.Mp.AppId = n.Data.Ptn.Mp.AppId
	tn.Mp.PagePath = n.Data.Ptn.Mp.PagePath
	tn.Data = make(map[string]TemplateData)
	for _, data := range n.Data.Ptn.Data {
		if data.Keyname != "" {
			tn.Data[data.Keyname] = TemplateData{data.Value, data.Color}
		}
	}

	var user User
	user.AppId = n.Metadata.AppId
	user.Secret = n.Metadata.Secret

	s.Set(fmt.Sprintf(JOB_STATUS, user.AppId, n.Data.Ptn.TemplateId), "running", INVALID_TTL)

	geted_token := false
	for {
		finish, _ := s.Get(fmt.Sprintf(PRODUCE_FINISH, user.AppId, n.Data.Ptn.TemplateId))
		queuecount, _ := s.Qsize(fmt.Sprintf(JOB_QUEUE, user.AppId))
		if finish == "ok" && queuecount == 0 {
			break
		}
		token, err := s.GetSetToken(&user)
		if err != nil {
			if !geted_token {
				beego.Error(err)
				s.Client.Do("qclear", fmt.Sprintf(JOB_QUEUE, user.AppId))
				s.Hset(fmt.Sprintf(NOTIFY_FAIL, user.AppId, n.Data.Ptn.TemplateId), "noopenid", "no way to get token")
				break
			} else {
				continue
			}
		}

		geted_token = true
		openid, _ := s.Qpop(fmt.Sprintf(JOB_QUEUE, user.AppId))
		if openid != "" {
			tn.ToUser = openid
			ersp := send(token, &tn)
			if ersp.ErrCode != 0 {
				s.Hset(fmt.Sprintf(NOTIFY_FAIL, user.AppId, n.Data.Ptn.TemplateId), openid, ersp.ErrMsg)
			} else {
				s.Zset(fmt.Sprintf(NOTIFY_SUCCESS, user.AppId, n.Data.Ptn.TemplateId), openid)
			}
		}
	}

	jobid := fmt.Sprintf("%s_%s", n.Metadata.AppId, n.Data.Ptn.TemplateId)
	s.Client.Do("zdel", NOTIFY_NOTICE_BOARD, jobid)
	s.Client.Do("hdel", NOTIFY_NOTICE_DETAIL, jobid)
	s.Set(fmt.Sprintf(JOB_END_TIME, user.AppId, n.Data.Ptn.TemplateId), time.Now().Unix(), INVALID_TTL)
	s.Set(fmt.Sprintf(JOB_STATUS, user.AppId, n.Data.Ptn.TemplateId), "finish", INVALID_TTL)
	return nil

}

func send(token string, tn *templateNotify) *ErrResponse {
	var ersp ErrResponse
	notify_json, err := json.Marshal(&tn)
	if err != nil {
		beego.Error(err)
		return &JsonParseErr
	}
	beego.Informational("notify to", string(notify_json))
	body := bytes.NewBuffer([]byte(notify_json))

	u, _ := url.Parse(PUBLICNUMBER_TEMPLATE_NOTIFY_API)
	q := u.Query()

	q.Set("access_token", token)
	u.RawQuery = q.Encode()

	resp, err := http.Post(u.String(), "application/json;charset=utf-8", body)
	if err != nil {
		beego.Error(err)
		return &RemoteServerErr
	}

	result, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		beego.Error(err)
		return &ResponseErr
	}

	err = json.Unmarshal([]byte(result), &ersp)
	if err != nil {
		beego.Error(err)
		return &JsonParseErr
	}
	return &ersp
}

func GetUsers(token, next_openid string) (*users, error) {
	u, _ := url.Parse(PUBLICNUMBER_GETUSERS_API)
	q := u.Query()

	q.Set("access_token", token)
	if next_openid != "" {
		q.Set("next_openid", next_openid)
	}
	u.RawQuery = q.Encode()

	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	var ersp ErrResponse
	err = json.Unmarshal(result, &ersp)
	if err != nil {
		return nil, err
	}
	if ersp.ErrCode != 0 {
		return nil, errors.New(ersp.ErrMsg)
	}

	var usr users
	err = json.Unmarshal(result, &usr)
	if err != nil {
		return nil, err
	}
	return &usr, nil
}
