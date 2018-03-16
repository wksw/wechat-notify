package miniprogrammodels

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

type templateNotify struct {
	ToUser          string                  `json:"touser"`
	TemplateId      string                  `json:"template_id"`
	Page            string                  `json:"page"`
	FormId          string                  `json:"form_id"`
	Data            map[string]TemplateData `json:"data"`
	EmphasisKeyword string                  `json:"emphasis_keyword"`
}

type job struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Ntf     Notify `json:"notify"`
}

func Run(n *Notify) error {
	beego.Informational("miniprogram notify start")
	cpunum := runtime.NumCPU()
	for i := 0; i < cpunum; i++ {
		go notify(*n)
	}
	return nil
}

func Produce(appid, templateid string) error {
	beego.Informational("miniprogram start push openid into job queue")
	if appid == "" {
		beego.Error("appid is empty")
		return errors.New("appid is empty")
	}
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
	err = s.Zset(fmt.Sprintf(NOTIFY_TEMPLATES, appid), templateid)
	if err != nil {
		beego.Error(err)
		return err
	}

	size, err := s.Zsize(fmt.Sprintf(UNKNOWN_TEMPLATE_ID, appid))
	if err != nil {
		beego.Error(err)
		return err
	}

	err = s.Set(fmt.Sprintf(PRODUCE_TOTAL, appid, templateid), size, INVALID_TTL)
	if err != nil {
		beego.Error(err)
		return err
	}

	maxpush := config.MaxPush
	if size > maxpush {
		m := 1
		if size%maxpush == 0 {
			m = size / maxpush
		} else {
			m = size/maxpush + 1
		}
		start_name := ""
		for i := 0; i < m; i++ {
			if i == size-1 {
				maxpush = size % maxpush
			}
			openid_arry, err := s.Zscan(fmt.Sprintf(UNKNOWN_TEMPLATE_ID, appid), start_name, maxpush)
			if err != nil {
				beego.Error("get notify openid list fail", err)
				return err
			}
			if openid_arry != nil {
				start_name = openid_arry[len(openid_arry)-1]
				err = s.Qpush(fmt.Sprintf(JOB_QUEUE, appid, templateid), openid_arry)
				if err != nil {
					beego.Error("push openid into job queue fail", err)
					return err
				}
			}
		}

	} else {
		openid_arry, err := s.Zscan(fmt.Sprintf(UNKNOWN_TEMPLATE_ID, appid), "", size)
		if err != nil {
			beego.Error("get notify openid list fail", err)
			return err
		}
		if openid_arry != nil {
			err = s.Qpush(fmt.Sprintf(JOB_QUEUE, appid, templateid), openid_arry)
			if err != nil {
				beego.Error("push openid into job queue fail", err)
				return err
			}
		}
	}

	err = s.Set(fmt.Sprintf(PRODUCE_FINISH, appid, templateid), "ok", INVALID_TTL)
	if err != nil {
		beego.Error("set finish flag fail", err)
		return err
	}
	beego.Informational("Finish push job into queue", fmt.Sprintf(JOB_QUEUE, appid, templateid))
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
	tn.Page = n.Data.Mtn.Page
	tn.EmphasisKeyword = n.Data.Mtn.EmphasisKeyword
	tn.TemplateId = n.Data.Mtn.TemplateId
	tn.Data = make(map[string]TemplateData)
	for _, data := range n.Data.Mtn.Data {
		if data.Keyname != "" {
			tn.Data[data.Keyname] = TemplateData{data.Value, data.Color}
		}
	}
	var user User
	user.AppId = n.Metadata.AppId
	user.Secret = n.Metadata.Secret

	s.Set(fmt.Sprintf(JOB_STATUS, user.AppId, n.Data.Mtn.TemplateId), "running", INVALID_TTL)

	geted_token := false
	for {
		finish, _ := s.Get(fmt.Sprintf(PRODUCE_FINISH, user.AppId, n.Data.Mtn.TemplateId))
		queuecount, _ := s.Qsize(fmt.Sprintf(JOB_QUEUE, user.AppId, n.Data.Mtn.TemplateId))
		if finish == "ok" && queuecount == 0 {
			break
		}

		token, err := s.GetSetToken(&user)
		if err != nil {
			if !geted_token {
				beego.Error(err)
				s.Client.Do("qclear", fmt.Sprintf(JOB_QUEUE, user.AppId, n.Data.Mtn.TemplateId))
				s.Hset(fmt.Sprintf(NOTIFY_FAIL, user.AppId, n.Data.Mtn.TemplateId), "noopenid", "no way to get token")
				break
			} else {
				continue
			}
		}
		geted_token = true

		s_flag := false
		failMsg := ""
		openid, _ := s.Qpop(fmt.Sprintf(JOB_QUEUE, user.AppId, n.Data.Mtn.TemplateId))
		if openid != "" {
			size, _ := s.Qsize(openid)
			if size == 0 {
				s.Zset(fmt.Sprintf(NOTIFY_INVALID, user.AppId, n.Data.Mtn.TemplateId), openid)
				continue
			}
			for i := 0; i < size; i++ {
				formid, _ := s.GetFormIdByOpenId(openid)
				if formid != "" {
					tn.ToUser = openid
					tn.FormId = formid
					err := send(token, &tn)
					if err.ErrCode != 0 {
						failMsg = fmt.Sprintf("user formid %s fail: %s", formid, err.ErrMsg)
						if err.ErrCode != 41028 && err.ErrCode != 41029 {
							break
						}
					} else {
						s.Zset(fmt.Sprintf(NOTIFY_SUCCESS, user.AppId, n.Data.Mtn.TemplateId), openid)
						s_flag = true
						break
					}
				} else {
					failMsg = "formid expired"
				}
			}
			if !s_flag && size != 0 {
				s.Hset(fmt.Sprintf(NOTIFY_FAIL, user.AppId, n.Data.Mtn.TemplateId), openid, failMsg)
			}

			// clear
			size, _ = s.Qsize(openid)
			if size == 0 {
				s.Client.Do("zdel", fmt.Sprintf(UNKNOWN_TEMPLATE_ID, user.AppId), openid)
			}
		}
	}

	// clear job from notice_board

	jobid := fmt.Sprintf("%s_%s", n.Metadata.AppId, n.Data.Mtn.TemplateId)
	s.Client.Do("zdel", NOTIFY_NOTICE_BOARD, jobid)
	s.Client.Do("hdel", NOTIFY_NOTICE_DETAIL, jobid)
	s.Set(fmt.Sprintf(JOB_END_TIME, user.AppId, n.Data.Mtn.TemplateId), time.Now().Unix(), INVALID_TTL)
	s.Set(fmt.Sprintf(JOB_STATUS, user.AppId, n.Data.Mtn.TemplateId), "finish", INVALID_TTL)
	return nil

}

func send(token string, tn *templateNotify) *ErrResponse {
	var ersp ErrResponse

	notify_json, err := json.Marshal(&tn)
	if err != nil {
		beego.Error(err)
		return &JsonParseErr
	}
	beego.Informational("notify to:", string(notify_json))

	body := bytes.NewBuffer([]byte(notify_json))

	u, _ := url.Parse(MINIPROGRAM_TEMPLATE_NOTIFY_API)
	q := u.Query()
	q.Set("access_token", token)
	u.RawQuery = q.Encode()
	resp, err := http.Post(u.String(), "application/json;charset=utf-8", body)
	if err != nil {
		beego.Error(err)
		return &RemoteServerErr
	}

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Error(err)
		return &ResponseErr
	}
	defer resp.Body.Close()
	err = json.Unmarshal([]byte(result), &ersp)
	if err != nil {
		beego.Error(err)
		return &JsonParseErr
	}

	return &ersp
}
