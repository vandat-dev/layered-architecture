package kafka

import (
	"app/global"
	"app/internal/modules/delivery_frame/service"
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type FrameProcessor struct {
	scanService service.IScanService
}

func ProcessKafkaFrame(scanService service.IScanService) *FrameProcessor {
	return &FrameProcessor{
		scanService: scanService,
	}
}

func (p *FrameProcessor) Handle(ctx context.Context, msg kafka.Message) error {
	// Only process frame topics
	// Assuming topic name is configured or we check prefix
	// For now, let's just process everything or check topic
	// But StartKafkaConsumer subscribes to specific topics.

	// Extract headers
	var deviceID, scanID string
	for _, h := range msg.Headers {
		if h.Key == "device_id" {
			deviceID = string(h.Value)
		} else if h.Key == "scan_id" {
			scanID = string(h.Value)
		}
	}

	if deviceID == "" || scanID == "" {
		global.Logger.Error(fmt.Sprintf("[FRAME-PROCESSOR] Missing headers. DeviceID: %s, ScanID: %s", deviceID, scanID))
		return nil // Skip invalid message
	}

	// Call service
	// msg.Value contains the image bytes
	if err := p.scanService.ProcessKafkaFrame(ctx, deviceID, scanID, msg.Value); err != nil {
		global.Logger.Error(fmt.Sprintf("[FRAME-PROCESSOR] Failed to process frame: %v", err))
		return err
	}

	return nil
}
