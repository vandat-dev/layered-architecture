package routers

import (
	"app/internal/modules/user/router"
)

type RouterGroup struct {
	User router.UsersRouterGroup
}

var RouterGroupApp = new(RouterGroup)
