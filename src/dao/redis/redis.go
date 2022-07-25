package redis

import (
	"fmt"
	"goScaffold/settings"

	"github.com/go-redis/redis"
)

var rdb *redis.Client

func Init(cfg *settings.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Network: "tcp", // tcp适合长连接
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host,
			cfg.Port,
		),
		Password: cfg.Password,
		DB:       cfg.Db,
		PoolSize: cfg.PoolSize,
	})

	status := rdb.Ping()
	if status.Err() != nil {
		return status.Err()
	}

	return
}

func Close() {
	_ = rdb.Close()
}
