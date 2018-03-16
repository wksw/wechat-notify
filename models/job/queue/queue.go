package queue

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	. "wksw/notify/models"
	"wksw/notify/models/amqp"
	. "wksw/notify/models/db"
	. "wksw/notify/models/job"
)

type Queue struct {
	Jb *Job
}

// add job from waiting queue to running queue
func (q *Queue) Push() error {
	config, err := GetConfig()
	if err != nil {
		return err
	}

	s, err := NewSsdb(config.Db.WriteHost, config.Db.WritePort)
	if err != nil {
		return err
	}
	defer s.Client.Close()

	switch q.Jb.Type {
	case NOTIFY_TYPE:
		appid := q.Jb.Ntf.Metadata.AppId
		job_str, err := json.Marshal(q.Jb)
		if err != nil {
			return err
		}

		jobid := ""
		switch q.Jb.Ntf.Type {
		case MINIPROGRAM_NOTIFY:
			s.Set(fmt.Sprintf(JOB_STATUS, appid, q.Jb.Ntf.Data.Mtn.TemplateId), "waiting", INVALID_TTL)
			jobid = fmt.Sprintf("%s_%s", q.Jb.Ntf.Metadata.AppId, q.Jb.Ntf.Data.Mtn.TemplateId)
		case PUBLICNUM_NOTIFY:
			s.Set(fmt.Sprintf(JOB_STATUS, appid, q.Jb.Ntf.Data.Ptn.TemplateId), "waiting", INVALID_TTL)
			jobid = fmt.Sprintf("%s_%s", q.Jb.Ntf.Metadata.AppId, q.Jb.Ntf.Data.Ptn.TemplateId)
		default:
			return errors.New("unknown notify type")

		}

		err = s.Qpush(NOTIFY_NOTICE_QUEUE, jobid)
		if err != nil {
			return err
		}
		err = s.Hset(NOTIFY_NOTICE_DETAIL, jobid, job_str)
		if err != nil {
			return err
		}
	default:
		return errors.New("unknown job type")
	}

	return nil
}

func JobInit() error {
	go templateNotifyPublishInit()
	err := startInit()
	if err != nil {
		return err
	}
	return nil
}

func startInit() error {
	config, err := GetConfig()
	if err != nil {
		return err
	}

	s, err := NewSsdb(config.Db.WriteHost, config.Db.WritePort)
	if err != nil {
		return err
	}
	defer s.Client.Close()

	job_arry, err := s.Zscan(NOTIFY_NOTICE_BOARD, "", config.MaxJob)
	for _, v := range job_arry {
		var job Job
		job.Id = v
		err = job.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func templateNotifyPublishInit() {
	timer := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-timer.C:
			templateNotifyPublish()
		}
	}
}

func templateNotifyPublish() error {
	config, _ := GetConfig()
	s, err := NewSsdb(config.Db.WriteHost, config.Db.WritePort)
	if err != nil {
		return err
	}
	defer s.Client.Close()
	lock, _ := s.Get(NOTIFY_JOB_LOCK)
	size, _ := s.Zsize(NOTIFY_NOTICE_BOARD)
	maxjob := config.MaxJob

	if size < maxjob {
		if lock != "locked" {
			s.Set(NOTIFY_JOB_LOCK, "locked", INVALID_TTL)
			for j := 0; j < maxjob-size; j++ {
				jobid, _ := s.Qpop(NOTIFY_NOTICE_QUEUE)
				if jobid != "" {
					// job, err := s.Hget(NOTIFY_NOTICE_DETAIL, jobid)
					fmt.Println("push job", jobid, "into queue")
					err := publish(config, jobid)
					if err != nil {
						return err
					}
					s.Zset(NOTIFY_NOTICE_BOARD, jobid)
				} else {
					break
				}
			}
			s.Set(NOTIFY_JOB_LOCK, "unlocked", INVALID_TTL)
		}
	}
	return nil
}

func publish(config *Config, jobid string) error {
	switch config.Amqp.Amqp {
	case "redis":
		redis, err := amqp.NewRedisPool(config.Amqp.ManageHost, config.Amqp.Password, config.Amqp.ManagePort)
		if err != nil {
			return err
		}
		c := redis.Pool.Get()
		defer c.Close()
		_, err = c.Do("PING")
		if err != nil {
			return err
		}
		err = redis.Publish(NOTIFY_TYPE, jobid)
		if err != nil {
			return err
		}
	case "rabbitmq":
		return errors.New("unimplemented")
	default:
		return errors.New("parse config fail")
	}
	return nil
}
