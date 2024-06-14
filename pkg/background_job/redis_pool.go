package background_job

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func Pool(host, port string) *redis.Pool {
	return &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
		},
	}
}
