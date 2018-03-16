package job

import (
	"encoding/json"
	"errors"
	"fmt"
	. "wksw/notify/models"
	. "wksw/notify/models/db"
	. "wksw/notify/models/notify"
	"wksw/notify/models/notify/miniprogrammodels"
	"wksw/notify/models/notify/publicnummodels"
)

type Job struct {
	Id      string `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Ntf     Notify `json:"notify"`
}

func (j *Job) Run() error {
	config, err := GetConfig()
	if err != nil {
		return err
	}

	s, err := NewSsdb(config.Db.WriteHost, config.Db.WritePort)
	if err != nil {
		return err
	}
	defer s.Client.Close()

	job_detail, err := s.Hget(NOTIFY_NOTICE_DETAIL, j.Id)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(job_detail), &j)
	if err != nil {
		return err
	}

	switch j.Type {
	case NOTIFY_TYPE:
		err := j.notify()
		if err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("undefined job type %s", j.Type))
	}
	return nil
}

func (j *Job) notify() error {

	switch j.Ntf.Type {
	case MINIPROGRAM_NOTIFY:
		// TODO miniprogram notify
		miniprogrammodels.Run(&j.Ntf)
	case PUBLICNUM_NOTIFY:
		// TODO publicnum notify
		publicnummodels.Run(&j.Ntf)
	default:
		return errors.New(fmt.Sprintf("unknown notify type '%s'", j.Ntf.Type))

	}
	return nil
}
