package initialize

import (
	"app/global"
	"app/pkg/logger"
)

func InitLogger() {
	global.Logger = logger.NewLogger(global.Config.Logger)
}
