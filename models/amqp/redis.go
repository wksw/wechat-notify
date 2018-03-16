package amqp

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"time"
	. "wksw/notify/models"
	"wksw/notify/models/job"
)

type RedisPool struct {
	Pool *redis.Pool
}

func NewRedisPool(host, password string, port int) (*RedisPool, error) {

	if host == "" {
		host = "localhost"
	}

	if port == 0 {
		port = 6379
	}

	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			server := fmt.Sprintf("%s:%d", host, port)
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}

			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return &RedisPool{pool}, nil
}

func (r *RedisPool) Publish(channel, jobid string) error {
	c := r.Pool.Get()
	defer c.Close()
	_, err := c.Do("PUBLISH", channel, jobid)
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisPool) Subscribe(j *job.Job) {
	c := r.Pool.Get()
	defer c.Close()

	psc := redis.PubSubConn{c}
	psc.Subscribe(j.Channel)
	for {
		// _, err := c.Do("PING")
		// if err != nil {
		// 	beego.Error("connect to redis fail, will reconnect")
		// 	c = r.Pool.Get()
		// 	psc := redis.PubSubConn{c}
		// 	psc.Subscribe(j.Channel)
		// }

		switch v := psc.Receive().(type) {
		case redis.Message:
			switch j.Type {
			case NOTIFY_TYPE:
				beego.Informational("job", j.Type, "start")
				j.Id = string(v.Data)
				err := j.Run()
				if err != nil {
					beego.Error(err)
				}
			}
		// case redis.Subscription:
		// 	//
		case error:
			beego.Error(v)
		}
	}
}
