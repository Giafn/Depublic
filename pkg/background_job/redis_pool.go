package background_job

import (
	"fmt"

	"github.com/Giafn/Depublic/configs"
	"github.com/gomodule/redigo/redis"
)

func Pool(configs *configs.Config) *redis.Pool {
	return &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%s", configs.Redis.Host, configs.Redis.Port))
		},
	}
}
