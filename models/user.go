package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	. "wksw/notify/models/error"
	"github.com/astaxie/beego/validation"
	// "time"
)

type CustomToken struct {
	AppId string `json:"appid"`
	TokenUrl string `json:"tokenurl"`
}

type User struct {
	AppId  string `json:"appid"`
	Secret string `json:"secret"`
}

type Token struct {
	Access_Token string `json:"access_token"`
	Expires_In   int64  `json:"expires_in"`
}

type rever_token struct {
	Access_token string `json:"access_token"`
	Expire_time  int64 `json:"expire_time"`
}

func (p *CustomToken) Validation() error {
	valid := validation.Validation{}
	if v := valid.Required(p.AppId, "appid"); !v.Ok {
		return errors.New(fmt.Sprintf("%s:%s", v.Error.Key, v.Error.Message))
	}
	return nil
}

func GetToken(user *User) (Token, error) {
	var token Token
	u, _ := url.Parse(GETTOKEN_API)
	q := u.Query()
	q.Set("grant_type", "client_credential")
	q.Set("appid", user.AppId)
	q.Set("Secret", user.Secret)
	u.RawQuery = q.Encode()
	res, err := http.Get(u.String())
	if err != nil {
		return token, err
	}

	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return token, err
	}
	var resp ErrResponse
	err = json.Unmarshal(result, &resp)
	if err != nil {
		return token, err
	}
	if resp.ErrCode != 0 {
		return token, errors.New(fmt.Sprintf("[%d]%s", resp.ErrCode, resp.ErrMsg))
	}

	err = json.Unmarshal(result, &token)
	if err != nil {
		return token, err
	}
	return token, nil
}

func CustomGetToken(tokenurl string) (Token, error) {
	var token Token
	u, _ := url.Parse(tokenurl)
	q := u.Query()
	u.RawQuery = q.Encode()
	res, err := http.Get(u.String())
	if err != nil {
		return token, err
	}

	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return token, err
	}
	var resp ErrResponse
	err = json.Unmarshal(result, &resp)
	if err != nil {
		return token, err
	}
	if resp.ErrCode != 0 {
		return token, errors.New(fmt.Sprintf("[%d]%s", resp.ErrCode, resp.ErrMsg))
	}

	err = json.Unmarshal(result, &token)
	if err != nil {
		return token, err
	}
	return token, nil


	
	// uu, _ := url.Parse(tokenurl)
	// q := uu.Query()

	// uu.RawQuery = q.Encode()

	// res, err := http.Get(uu.String())
	// if err != nil {
	// 	return token, err
	// }

	// result, err := ioutil.ReadAll(res.Body)
	// res.Body.Close()
	// if err == nil {
	// 	var rt rever_token
	// 	err = json.Unmarshal(result, &rt)
	// 	if err != nil {
	// 		return token, err
	// 	}
	// 	if rt.Access_token == "" || rt.Expire_time == 0 {
	// 		return token, errors.New("access_token empty or expire_time is 0")
			
	// 	}

	// 	token.Access_Token = rt.Access_token
	// 	now := time.Now().Unix()
	// 	expires_in := rt.Expire_time - now
	// 	if expires_in < 0 {
	// 		expires_in = 0
	// 	}
	// 	token.Expires_In = expires_in
	// 	return token, nil
	// } else {
	// 	return token, err
	// }
}
