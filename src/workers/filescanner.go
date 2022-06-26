package workers

import (
	"time"

	"github.com/gocelery/gocelery"
	"github.com/gomodule/redigo/redis"
)

type directoryScanTask struct {
	directory string
	scanTime  int
}

func (t *directoryScanTask) RunTask() (interface{}, error) {
	var i int = 0

	for i <= t.scanTime {
		//time.Sleep(10 * time.Second)
		println("stuff")
		i++
	}
	return i, nil
}

func StartWorkerDirectoryScan() {
	// create redis connection pool
	redisPool := &redis.Pool{
		MaxIdle:     3,                 // maximum number of idle connections in the pool
		MaxActive:   0,                 // maximum number of connections allocated by the pool at a given time
		IdleTimeout: 240 * time.Second, // close connections after remaining idle for this duration
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL("redis://redis:eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81@172.20.0.5:6379")
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	// initialize celery client
	cli, _ := gocelery.NewCeleryClient(
		gocelery.NewRedisBroker(redisPool),
		&gocelery.RedisCeleryBackend{Pool: redisPool},
		5, // number of workers
	)

	cli.Register("directoryScan_reflect", &directoryScanTask{})
	cli.StartWorker()

	println("Directory Scan Worker Init")
}
