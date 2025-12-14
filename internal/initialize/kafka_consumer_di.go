package initialize

import (
	"app/global"
	"app/internal/modules/delivery_frame/repo"
	"app/internal/modules/delivery_frame/service"
	"app/internal/third_party/kafka"
	"app/internal/third_party/redis"
	"app/internal/third_party/s3"
)

func InitKafkaConsumer() {
	// Init Redis Provider
	redisProvider := redis.NewRedisProvider()
	s3Provider := s3.NewS3Provider()

	// Init Kafka Delivery
	deliveryFrameRepo := repo.NewScanRepository(global.Postgres)
	scanService := service.NewScanService(deliveryFrameRepo, s3Provider, redisProvider)
	deliveryHandler := kafka.ProcessKafkaFrame(scanService)

	StartKafkaConsumer(deliveryHandler)
}
