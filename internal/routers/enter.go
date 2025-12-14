package routers

import (
	deliveryRouter "app/internal/modules/delivery_frame/router"
	"app/internal/modules/user/router"
)

type RouterGroup struct {
	User          router.UsersRouterGroup
	DeliveryFrame deliveryRouter.ScanRouterGroup
}

var RouterGroupApp = new(RouterGroup)
