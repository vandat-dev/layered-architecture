package initialize

import (
	"app/global"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var ctx = context.Background()

func Redis() {
	global.Logger.Info("Redis connecting!")
	r := global.Config.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", r.Host, r.Port),
		Password: r.Password,
		DB:       r.Database,
		PoolSize: 10,
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		global.Logger.Error("Redis connect failed!", zap.Error(err))
		panic(err)
	}
	global.Redis = rdb
	global.Logger.Info("Redis connect success!")
}
