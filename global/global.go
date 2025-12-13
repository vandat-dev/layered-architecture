package global

import (
	"app/pkg/logger"
	"app/pkg/setting"

	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Config   setting.Config
	Logger   *logger.LogZap
	Redis    *redis.Client
	MinIO    *minio.Client
	Postgres *gorm.DB
)

/*
Config: Redis, Mysql, Postgres, WebSocket Manager, ...
*/
