package db

import (
	"github.com/ideawu/ssdb"
	// "github.com/astaxie/beego/validation"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
	. "wksw/notify/models"
	. "wksw/notify/models/notify"
)

type SsdbClient struct {
	Client *ssdb.Client
}

func NewSsdb(address string, port int) (*SsdbClient, error) {
	if address == "" {
		address = "127.0.0.1"
	}
	if port == 0 {
		port = 8888
	}
	client, err := ssdb.Connect(address, port)
	if err != nil {
		return nil, err
	}
	return &SsdbClient{client}, nil
}

func (s *SsdbClient) Set(key string, value interface{}, ttl int64) error {
	if key == "" {
		return errors.New("key empty")
	}
	client := s.Client
	var resp []string
	var err error
	if ttl == INVALID_TTL {
		resp, err = client.Do("set", key, value)
	} else {
		resp, err = client.Do("setx", key, value, ttl)
	}
	if err != nil {
		return err
	}
	if len(resp) == 0 {
		return errors.New("ssdb no response")
	}
	if resp[0] != "ok" {
		return errors.New(resp[0])
	}
	return nil
}

func (s *SsdbClient) Get(key string) (string, error) {
	if key == "" {
		return "", errors.New("key empty")
	}
	client := s.Client
	resp, err := client.Do("get", key)
	if err != nil {
		return "", err
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return resp[1], nil
	}
	return "", errors.New("unknow error")
}

func (s *SsdbClient) Hset(name, key, value interface{}) error {
	if key == "" || name == "" {
		return errors.New("key or name empty")
	}
	client := s.Client

	resp, err := client.Do("hset", name, key, value)
	if err != nil {
		return err
	}
	if len(resp) == 0 {
		return errors.New("ssdb no response")
	}
	if resp[0] != "ok" {
		return errors.New(resp[0])
	}
	return nil
}

func (s *SsdbClient) Hget(name, key string) (string, error) {
	if key == "" || name == "" {
		return "", errors.New("key or name empty")
	}
	client := s.Client

	resp, err := client.Do("hget", name, key)
	if err != nil {
		return "", err
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return resp[1], nil
	}
	return "", errors.New("unknow error")
}

func (s *SsdbClient) Hsize(key string) (int, error) {
	if key == "" {
		return 0, errors.New("key empty")
	}
	client := s.Client

	size, err := client.Do("hsize", key)
	if err != nil {
		return 0, err
	}

	if len(size) == 2 && size[0] == "ok" {
		d, err := strconv.Atoi(size[1])
		if err != nil {
			return 0, err
		}
		return d, nil
	}
	return 0, errors.New(fmt.Sprintf("ssdb no response: %v", size))
}

func (s *SsdbClient) Hscan(key, start string, limit int) ([]*FailDetail, error) {
	var failDetail_arry []*FailDetail
	if key == "" {
		return nil, errors.New("key empty")
	}
	client := s.Client

	resp, err := client.Do("hscan", key, start, "", limit)
	if err != nil {
		return nil, err
	}
	if len(resp) >= 2 && resp[0] == "ok" {
		for i := 1; i < len(resp); i += 2 {
			failDetail_arry = append(failDetail_arry, &FailDetail{resp[i], resp[i+1]})
		}
		return failDetail_arry, nil
	}
	return nil, errors.New("unknown error")
}

func (s *SsdbClient) Zset(name, key interface{}) error {
	if name == "" {
		return errors.New("name empty")
	}
	client := s.Client

	resp, err := client.Do("zset", name, key, ZSET_WEIGHT)
	if err != nil {
		return err
	}
	if len(resp) == 0 {
		return errors.New("ssdb no response")
	}
	if resp[0] != "ok" {
		return errors.New(resp[0])
	}
	return nil
}

func (s *SsdbClient) Qpush(key string, values ...interface{}) error {
	if key == "" {
		return errors.New("key empty")
	}
	client := s.Client
	values = append(values, "qpush", key)
	values = append(values[len(values)-2:], values[:len(values)-2]...)
	resp, err := client.Do(values...)
	if err != nil {
		return err
	}
	if len(resp) == 0 {
		return errors.New("ssdb no response")
	}
	if resp[0] != "ok" {
		return errors.New(resp[0])
	}
	return err
}

func (s *SsdbClient) Qpop(key string) (string, error) {
	if key == "" {
		return "", errors.New("key empty")
	}
	client := s.Client

	resp, err := client.Do("qpop", key)
	if err != nil {
		return "", err
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return resp[1], nil
	}

	return "", errors.New("not want value")
}

func (s *SsdbClient) Zsize(key string) (int, error) {
	if key == "" {
		return 0, errors.New("key empty")
	}
	client := s.Client

	size, err := client.Do("zsize", key)
	if err != nil {
		return 0, err
	}

	if len(size) == 2 && size[0] == "ok" {
		d, err := strconv.Atoi(size[1])
		if err != nil {
			return 0, err
		}
		return d, nil
	}
	return 0, errors.New("ssdb no response")
}

func (s *SsdbClient) Zscan(key, start string, limit int) ([]string, error) {
	if key == "" {
		return nil, errors.New("key empty")
	}
	client := s.Client

	resp, err := client.Do("zscan", key, start, "", "", limit)
	if err != nil {
		return nil, err
	}
	if len(resp) >= 2 && resp[0] == "ok" {
		var openids []string
		for i := 1; i < len(resp); i += 2 {
			openids = append(openids, resp[i])
		}
		return openids, nil
	}
	if len(resp) >= 1 && resp[0] == "ok" {
		// empty no values
		return nil, nil
	}

	return nil, errors.New("unknown error")
}

func (s *SsdbClient) Qsize(key string) (int, error) {
	if key == "" {
		return 0, errors.New("key empty")
	}

	client := s.Client

	size, err := client.Do("qsize", key)
	if err != nil {
		return 0, err
	}

	if len(size) == 2 && size[0] == "ok" {
		d, err := strconv.Atoi(size[1])
		if err != nil {
			return 0, err
		}
		return d, nil
	}

	return 0, errors.New("ssdb no response")
}

func (s *SsdbClient) GetFormIdByOpenId(openid string) (string, error) {
	if openid == "" {
		return "", errors.New("openid empty")
	}

	// get unexpired formid
	size, err := s.Qsize(openid)
	if err == nil {
		for i := 0; i < size; i++ {
			formid, err := s.Qpop(openid)
			if err != nil {
				return "", err
			}
			var fi FormIdInfo
			err = json.Unmarshal([]byte(formid), &fi)
			if err != nil {
				return "", err
			}

			timestamp := fi.TimeStamp
			loc, _ := time.LoadLocation("Local")
			timelayout := "2006-01-02 15:04:05"
			timestampstr := time.Unix(timestamp, 0).Format(timelayout)
			thetime, _ := time.ParseInLocation(timelayout, timestampstr, loc)
			m, _ := time.ParseDuration(FORMID_EXPIRE_IN)
			m1 := thetime.Add(m)
			now := time.Now()
			if now.Before(m1) {
				return fi.FormId, nil
			}

		}
		return "", errors.New("no valid formid")
	} else {
		return "", err
	}

}

func (s *SsdbClient) TemplateNotifyInit(appid, templateid string) error {
	// init job start time

	var history History
	total_str, _ := s.Get(fmt.Sprintf(PRODUCE_TOTAL, appid, templateid))
	if total_str != "" {
		d, err := strconv.Atoi(total_str)
		if err == nil {
			history.Total = d
		}
	}
	success_size, _ := s.Zsize(fmt.Sprintf(NOTIFY_SUCCESS, appid, templateid))
	fail_size, _ := s.Hsize(fmt.Sprintf(NOTIFY_FAIL, appid, templateid))
	invalid_size, _ := s.Zsize(fmt.Sprintf(NOTIFY_INVALID, appid, templateid))
	start_at, _ := s.Get(fmt.Sprintf(JOB_START_TIME, appid, templateid))
	end_at, err := s.Get(fmt.Sprintf(JOB_END_TIME, appid, templateid))
	history.Success = success_size
	history.Fail = fail_size
	history.Invalid = invalid_size
	history.StartAt = start_at
	history.EndAt = end_at

	config, err := GetConfig()
	if err != nil {
		return err
	}

	history_str, _ := json.Marshal(history)
	max_history := config.MaxHistory
	next_job_num := 0
	next_job_num_str, _ := s.Get(fmt.Sprintf(NEXT_JOB_NUM, appid, templateid))
	if next_job_num_str != "" {
		d, err := strconv.Atoi(next_job_num_str)
		if err == nil {
			next_job_num = d
		}
	}

	err = s.Set(fmt.Sprintf(JOB_START_TIME, appid, templateid), time.Now().Unix(), INVALID_TTL)
	if err != nil {
		return errors.New("init time start fail")
	}

	size, _ := s.Hsize(fmt.Sprintf(HISTORY, appid, templateid))
	if size != 0 && size > max_history {
		s.Client.Do("hdel", fmt.Sprintf(HISTORY, appid, templateid), next_job_num-1-max_history)
	}
	if next_job_num != 0 {
		s.Hset(fmt.Sprintf(HISTORY, appid, templateid), next_job_num, history_str)
	}

	s.Set(fmt.Sprintf(NEXT_JOB_NUM, appid, templateid), next_job_num+1, INVALID_TTL)

	// init job end time
	err = s.Set(fmt.Sprintf(JOB_END_TIME, appid, templateid), 0, INVALID_TTL)
	if err != nil {
		return errors.New("init time end fail")
	}

	// init fail record
	_, err = s.Client.Do("hclear", fmt.Sprintf(NOTIFY_FAIL, appid, templateid))
	if err != nil {
		return errors.New("init fail record fail")
	}

	// init invalid record
	_, err = s.Client.Do("zclear", fmt.Sprintf(NOTIFY_INVALID, appid, templateid))
	if err != nil {
		return errors.New("init invalid record fail")
	}

	// init success record
	_, err = s.Client.Do("zclear", fmt.Sprintf(NOTIFY_SUCCESS, appid, templateid))
	if err != nil {
		return errors.New("init success record fail")
	}

	// init finish flag
	err = s.Set(fmt.Sprintf(PRODUCE_FINISH, appid, templateid), "", INVALID_TTL)
	if err != nil {
		return errors.New("init finish flag fail")
	}

	// init total record
	err = s.Set(fmt.Sprintf(PRODUCE_TOTAL, appid, templateid), 0, INVALID_TTL)
	if err != nil {
		return errors.New("init total record fail")
	}
	return nil
}

func (s *SsdbClient) StoreFormId(f *PostFormid, appid string) error {
	err := f.Validation()
	if err != nil {
		return err
	}

	// defer s.Client.Close()

	for _, data := range f.Data {
		openid := data.OpenId
		for _, formid := range data.FormIds {
			if formid.TemplateId == "" {
				err := s.Zset(fmt.Sprintf(UNKNOWN_TEMPLATE_ID, appid), openid)
				if err != nil {
					return err
				}

			} else {
				err := s.Zset(fmt.Sprintf(TEMPLATE_ID_LIST, appid), formid.TemplateId)
				if err != nil {
					return err
				}
				err = s.Zset(formid.TemplateId, openid)
				if err != nil {
					return err
				}

			}
			formid.TimeStamp = time.Now().Unix()
			fmd, err := json.Marshal(formid)
			if err != nil {
				return err
			}
			err = s.Qpush(openid, string(fmd))
			if err != nil {
				return err
			}

		}
	}
	return nil
}

func (s *SsdbClient) AddCustomToken(users []*CustomToken) error {
	for _, user := range users {
		err := s.Hset(CUSTOME_TOKEN, user.AppId, user.TokenUrl)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SsdbClient) DeleteCustomToken(users []string) error {
	for _, user := range users {
		if user != "" {
			_, err := s.Client.Do("hdel", CUSTOME_TOKEN, user)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *SsdbClient) GetSetToken(user *User) (string, error) {
	var t Token
	var err error
	token, _ := s.Get(fmt.Sprintf(TOKEN, user.AppId))
	if token == "" {
		if user.Secret != "" {
			t, err = GetToken(user)
		} else {
			tokenurl, err := s.Hget(CUSTOME_TOKEN, user.AppId)
			if err != nil {
				return "", err
			}
			t, err = CustomGetToken(tokenurl)
		}
		if err != nil {
			return "", err
		}
		token := t.Access_Token
		expires_in := t.Expires_In
		if token == "" || expires_in == 0 {
			return "", errors.New("get token from remote server is empty")
		}
		err = s.Set(fmt.Sprintf(TOKEN, user.AppId), token, expires_in)
		if err != nil {
			return "", err
		}
		return token, nil
	} else {
		return token, nil
	}
}
