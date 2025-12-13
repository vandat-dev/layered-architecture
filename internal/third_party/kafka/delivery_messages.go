package kafka

import (
	"app/global"
	"app/internal/modules/user/service"
	"context"
	"fmt"
	"strings"

	"github.com/segmentio/kafka-go"
)

type KafkaDeliveryMessages struct {
	userService service.IUserService
}

func NewKafkaDeliveryMessages(userService service.IUserService) *KafkaDeliveryMessages {
	return &KafkaDeliveryMessages{
		userService: userService,
	}
}

func (k *KafkaDeliveryMessages) Handle(ctx context.Context, msg kafka.Message) error {
	global.Logger.Info(fmt.Sprintf("[DELIVERY] Received message on topic %s", msg.Topic))

	if strings.HasPrefix(msg.Topic, "user_") {
		return k.userService.ReceiveMessages(msg.Value)
	} else if strings.HasPrefix(msg.Topic, "worker_") {
		//global.Logger.Info(fmt.Sprintf("[DELIVERY] Worker topic received: %s", msg.Topic))
		// Add worker service logic here
		return nil
	} else {
		global.Logger.Warn(fmt.Sprintf("[DELIVERY] Unknown topic: %s", msg.Topic))
		return nil
	}
}
