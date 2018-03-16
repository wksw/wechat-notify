package notify

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/validation"
	. "wksw/notify/models"
)

/*
{
	"type": "miniprogram|publicnumber",
	"metadata": {
		"appid": "",
		"secret": ""
	},
	"token": {
		"access_token": "",
		"expires_in": 7200
	}
	"custome_token": true | false
	"data": {
		"limit": 0,
		"miniprogram_template_notify": {
			"templateid": "",
			"page": "",
			"data": [
				{
					"keyname": "",
					"value": "",
					"color": ""
				}
			]

		},
		"publicnumber_template_notify": {

		}
	}
}
*/

type Notify struct {
	Type     string `json:"type"`
	Metadata User   `json:"metadata"`
	Data     Data   `json:"data"`
}
type Data struct {
	Mtn *MiniprogramTemplateNotify `json:"miniprogram_template_notify"`
	Ptn *PublicNumTemplateNotify   `json:"publicnumber_template_notify"`
}

type MiniprogramTemplateNotify struct {
	TemplateId      string          `json:"templateid"`
	Page            string          `json:"page"`
	Data            []*templateData `json:"data"`
	EmphasisKeyword string          `json:"emphasis_keyword"`
}

type PublicNumTemplateNotify struct {
	TemplateId string          `json:"templateid"`
	Url        string          `json:"url"`
	Mp         *Miniprogram    `json:"miniprogram"`
	Data       []*templateData `json:"data"`
}

type Miniprogram struct {
	AppId    string `json:"appid"`
	PagePath string `json:"pagepath"`
}

type templateData struct {
	Keyname string `json:"keyname"`
	Value   string `json:"value"`
	Color   string `json:"color"`
}

type TemplateData struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

type PostFormid struct {
	Data []*PostFormidData `json:"data"`
}

type PostFormidData struct {
	OpenId  string        `json:"openid"`
	FormIds []*FormIdInfo `json:"formids"`
}

type FormIdInfo struct {
	FormId     string `json:"formid"`
	TemplateId string `json:"templateid"`
	TimeStamp  int64  `json:"timestamp"`
}

type History struct {
	Total   int    `json:"total"`
	Success int    `json:"success"`
	Fail    int    `json:"fail"`
	Invalid int    `json:"invalid"`
	StartAt string `json:"start_at"`
	EndAt   string `json:"end_at"`
}

type NotifyState struct {
	Total         int           `json:"total"`
	Queue         int           `json:"queue"`
	Success       int           `json:"success"`
	Fail          int           `json:"fail"`
	Invalid       int           `json:"invalid"`
	Status        string        `json:"status"`
	Duration      int64         `json:"duration"`
	Page          Page          `json:"page"`
	FailDetail    []*FailDetail `json:"fail_detail"`
	InvalidDetail []string      `json:"invalid_detail"`
}

type Page struct {
	InvalidTotalPage int `json:"invalid_total_page"`
	FailTotalPage    int `json:"fail_total_page"`
	CurrentPage      int `json:"current_page"`
	PerPage          int `json:"perpage"`
}

type FailDetail struct {
	OpenId string `json:"openid"`
	ErrMsg string `json:"errmsg"`
}

func (f *PostFormid) Validation() error {
	valid := validation.Validation{}
	for _, data := range f.Data {
		// not empty
		if v := valid.Required(data.OpenId, "openid"); !v.Ok {
			return errors.New(fmt.Sprintf("%s:%s", v.Error.Key, v.Error.Message))
		}
		// check length
		if v := valid.MinSize(data.OpenId, OPENID_MINSIZE, "openid"); !v.Ok {
			return errors.New(fmt.Sprintf("%s:%s", v.Error.Key, v.Error.Message))
		}

		for _, formid := range data.FormIds {
			// not empty
			if v := valid.Required(formid.FormId, "formid"); !v.Ok {
				return errors.New(fmt.Sprintf("%s:%s", v.Error.Key, v.Error.Message))
			}
			// check length
			if v := valid.MinSize(formid.FormId, FORMID_MINSIZE, "formid"); !v.Ok {
				return errors.New(fmt.Sprintf("%s:%s", v.Error.Key, v.Error.Message))
			}
		}
	}
	return nil
}

func (m *MiniprogramTemplateNotify) Validation() error {
	valid := validation.Validation{}
	if v := valid.Required(m.TemplateId, "templateid"); !v.Ok {
		return errors.New(fmt.Sprintf("%s:%s", v.Error.Key, v.Error.Message))
	}
	return nil
}

func (p *PublicNumTemplateNotify) Validation() error {
	valid := validation.Validation{}
	if v := valid.Required(p.TemplateId, "templateid"); !v.Ok {
		return errors.New(fmt.Sprintf("%s:%s", v.Error.Key, v.Error.Message))
	}
	return nil
}
