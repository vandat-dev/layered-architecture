package setting

import (
	"time"
)

type Config struct {
	Server ServerSetting `map_structure:"server"`
	System SystemSetting `map_structure:"system"`

	Postgres PostgresSetting `map_structure:"postgres"`
	Redis    RedisSetting    `map_structure:"redis"`
	Kafka    KafkaSetting    `map_structure:"kafka"`
	MinIO    MinIOSetting    `map_structure:"minio"`
	Logger   LoggerSetting   `map_structure:"logger"`
	JWT      JWTSetting      `map_structure:"jwt"`
}

type ServerSetting struct {
	Port int    `map_structure:"port"`
	Mode string `map_structure:"mode"`
}

type SystemSetting struct {
	DefaultPassWord string `map_structure:"default_password"`
}

type PostgresSetting struct {
	Host            string `map_structure:"host"`
	Port            int    `map_structure:"port"`
	UserName        string `map_structure:"username"`
	Password        string `map_structure:"password"`
	DBName          string `map_structure:"dbname"`
	SSLMode         string `map_structure:"sslmode"`
	MaxIdleConn     int    `map_structure:"maxIdleConn"`
	MaxOpenConn     int    `map_structure:"maxOpenConn"`
	ConnMaxLifeTime int    `map_structure:"connMaxLifeTime"`
}

type LoggerSetting struct {
	LogLevel    string `map_structure:"log_level"`
	FileLogName string `map_structure:"file_log_name"`
	MaxBackups  int    `map_structure:"max_backups"`
	MaxSize     int    `map_structure:"max_size"`
	MaxAge      int    `map_structure:"max_age"`
	Compress    bool   `map_structure:"compress"`
}

type RedisSetting struct {
	Host     string `map_structure:"host"`
	Port     int    `map_structure:"port"`
	Password string `map_structure:"password"`
	Database int    `map_structure:"database"`
}

type KafkaSetting struct {
	Host    string   `map_structure:"host"`
	Port    int      `map_structure:"port"`
	Topics  []string `map_structure:"topics"`
	GroupID string   `map_structure:"group_id"`
}

type MinIOSetting struct {
	Endpoint        string `map_structure:"endpoint"`
	AccessKeyID     string `map_structure:"access_key_id"`
	SecretAccessKey string `map_structure:"secret_access_key"`
	BucketName      string `map_structure:"bucket_name"`
	UseSSL          bool   `map_structure:"use_ssl"`
}

type JWTSetting struct {
	SecretKey     string        `map_structure:"secret_key"`
	TokenExpiry   time.Duration `map_structure:"token_expiry"`
	RefreshExpiry time.Duration `map_structure:"refresh_expiry"`
}
