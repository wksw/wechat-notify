package keystone

import (
	"encoding/json"
	"net/http"
	"net/url"
	"bytes"
	"errors"
	"io/ioutil"
	"fmt"
	. "wksw/notify/models"
)

type K_Token struct {
	At K_Auth `json:"auth"`
}

type K_UserToken struct {
	At K_UserAuth `json:"auth"`
}

type K_UserAuth struct {
	Idt K_Identity `json:"identity"`
}

type K_Auth struct {
	Idt K_Identity `json:"identity"`
	Scp K_Scope `json:"scope"`
}

type K_Scope struct {
	Pj K_Project `json:"project"`
}

type K_Project struct {
	Dm K_Domain `json:"domain"`
	Name string `json:"name"`
}

type K_Identity struct {
	Methods []string `json:"methods"`
	Pwd K_Password `json:"password"`
}

type K_Password struct {
	Us K_User `json:"user"`
}

type K_User struct {
	Name string `json:"name"`
	Dm K_Domain `json:"domain"`
	Password string `json:"password"`
}

type K_Domain struct {
	Name string `json:"name"`
}


type UserCreate struct {
	Uinfo UserInfo `json:"user"`

}
type UserInfo struct {
	Default_project_id string `json:"default_project_id"`
	Domain_id string `json:"domain_id"`
	Enable bool `json:"enable"`
	Name  string `json:"name"`
	Password string `json:"password"`
	Description string `json:"description"`
	Email string `json:"email"`
}

type KeystoneErr struct {
	Einfo ErrInfo `json:"error"`
}

type ErrInfo struct {
	Message string `json:"message"`
	Code int `json:"code"errg`
	Title string `json:"title"`
}

type TokenDetail struct {
	Token TokenInfo `json:"token"`
}

type TokenInfo struct {
	Issued_at string `json:"issued_at"`
	Audit_ids []string `json:"audit_ids"`
	Methods []string `json:"methods"`
	Expires_at string `json:"expires_at"`
	User TokenInfoUser `json:"user"`
}

type TokenInfoUser struct {
	Password_exires_at string `json:"password_expires_at"`
	Domain TokenInfoUserDomain `json:"domain"`
	Id string `json:"id"`
	Name string `json:"name"`
}

type TokenInfoUserDomain struct {
	Id string `json:"id"`
	Name string `json:"name"`
}


func (u *UserCreate) Create() error {
	config, _ := GetConfig()

	var t K_Token
	t.At.Idt.Methods = append(t.At.Idt.Methods, "password")
	t.At.Idt.Pwd.Us.Name = config.Admin_username
	t.At.Idt.Pwd.Us.Dm.Name = config.Admin_domain_name
	t.At.Idt.Pwd.Us.Password = config.Admin_password
	t.At.Scp.Pj.Dm.Name = config.Admin_project_domain_name
	t.At.Scp.Pj.Name = config.Admin_project_name

	token, err := t.GetToken()
	if err != nil {
		return err
	}
	fmt.Println(token)

	req_str, err := json.Marshal(u)
	if err != nil {
		return err
	}
	fmt.Println(string(req_str))
	body := bytes.NewBuffer([]byte(req_str))

	client := &http.Client{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", config.Auth_url, USERS_API), body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", token)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var kerr KeystoneErr

	err = json.Unmarshal([]byte(result), &kerr)
	if err != nil {
		return err
	}
	if kerr.Einfo.Code != 0 {
		return errors.New(kerr.Einfo.Title)
	}
	return nil

}


func (t *K_Token) GetToken() (string, error){
	config, _ := GetConfig()
	auth_url := config.Auth_url
	req_str, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	fmt.Println(string(req_str))
	body := bytes.NewBuffer([]byte(req_str))

	u, _ := url.Parse(fmt.Sprintf("%s%s",auth_url, TOKENS_API))
	q := u.Query()
	u.RawQuery = q.Encode()
	resp, err := http.Post(u.String(), "application/json;charset=utf-8", body)
	if err != nil {
		return "", err
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Println(string(result))
	defer resp.Body.Close()

	token := resp.Header.Get("X-Subject-Token")
	if token == "" {
		return "", errors.New("get token empty")
	}
	return token, nil
}


func ValidToken(usertoken string) error {
	config, _ := GetConfig()

	var t K_Token
	t.At.Idt.Methods = append(t.At.Idt.Methods, "password")
	t.At.Idt.Pwd.Us.Name = config.Admin_username
	t.At.Idt.Pwd.Us.Dm.Name = config.Admin_domain_name
	t.At.Idt.Pwd.Us.Password = config.Admin_password
	t.At.Scp.Pj.Dm.Name = config.Admin_project_domain_name
	t.At.Scp.Pj.Name = config.Admin_project_name

	token, err := t.GetToken()
	if err != nil {
		return err
	}
	// body := bytes.NewBuffer([]byte(""))

	client := &http.Client{}

	req, err := http.NewRequest("HEAD", fmt.Sprintf("%s%s", config.Auth_url, TOKENS_API), nil)
	if err != nil {
		return err
	}
	// req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", token)
	req.Header.Set("X-Subject-Token", usertoken)
	fmt.Println(usertoken)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	status := resp.StatusCode
	if status != 200 {
		return errors.New(fmt.Sprintf("%d", status))
	}
	return nil
}

func (t *K_UserToken) GetUserToken() (string, error){
	config, _ := GetConfig()
	auth_url := config.Auth_url
	req_str, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	fmt.Println(string(req_str))
	body := bytes.NewBuffer([]byte(req_str))

	u, _ := url.Parse(fmt.Sprintf("%s%s",auth_url, TOKENS_API))
	q := u.Query()
	u.RawQuery = q.Encode()
	resp, err := http.Post(u.String(), "application/json;charset=utf-8", body)
	if err != nil {
		return "", err
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Println(string(result))
	defer resp.Body.Close()

	token := resp.Header.Get("X-Subject-Token")
	if token == "" {
		return "", errors.New("get token empty")
	}
	return token, nil
}


func GetUserInfo(usertoken string) (*TokenDetail, error) {
	config, _ := GetConfig()

	var t K_Token
	t.At.Idt.Methods = append(t.At.Idt.Methods, "password")
	t.At.Idt.Pwd.Us.Name = config.Admin_username
	t.At.Idt.Pwd.Us.Dm.Name = config.Admin_domain_name
	t.At.Idt.Pwd.Us.Password = config.Admin_password
	t.At.Scp.Pj.Dm.Name = config.Admin_project_domain_name
	t.At.Scp.Pj.Name = config.Admin_project_name

	token, err := t.GetToken()
	if err != nil {
		return nil, err
	}
	// body := bytes.NewBuffer([]byte(""))

	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", config.Auth_url, TOKENS_API), nil)
	if err != nil {
		return nil, err
	}
	// req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", token)
	req.Header.Set("X-Subject-Token", usertoken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var tokeninfo TokenDetail
	err = json.Unmarshal([]byte(result), &tokeninfo)
	if err != nil {
		return nil, err
	}

	fmt.Println("=========", string(result))
	return &tokeninfo, nil
}

