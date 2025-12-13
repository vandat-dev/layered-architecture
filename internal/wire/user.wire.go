//go:build wireinject

package wire

import (
	"app/global"
	"app/internal/modules/user/controller"
	"app/internal/modules/user/repo"
	"app/internal/modules/user/service"

	"app/internal/third_party/redis"

	"github.com/google/wire"
	"gorm.io/gorm"
)

func ProvideDB() *gorm.DB {
	return global.Postgres
}

func InitUserRouterHandler() (*controller.UserController, error) {
	wire.Build(
		ProvideDB,
		redis.NewRedisProvider,
		repo.NewUserRepository,
		service.NewUserService,
		controller.NewUserController,
	)
	return new(controller.UserController), nil
}
