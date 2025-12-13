package initialize

import (
	"app/global"
	"app/internal/modules/user/repo"
	"app/internal/modules/user/service"
	"app/internal/third_party/kafka"
	"app/internal/third_party/redis"
)

func InitKafkaConsumer() {
	// Init Redis Provider
	redisProvider := redis.NewRedisProvider()

	// Init Kafka Delivery
	userRepo := repo.NewUserRepository(global.Postgres)
	userService := service.NewUserService(userRepo, redisProvider)
	deliveryHandler := kafka.NewKafkaDeliveryMessages(userService)

	StartKafkaConsumer(deliveryHandler)
}
