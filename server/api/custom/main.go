package custom

import (
	"pocketbase-ruoyi/tools"

	"pocketbase-ruoyi/api/auth"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// RegisterCustom 注册自定义的路由
func RegisterCustom(app *pocketbase.PocketBase) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/api/custom/api", func(e *core.RequestEvent) error {

			return tools.JSONSuccess(e, map[string]any{
				"custom_api": true,
			})

		}).BindFunc(
			auth.
				// check rbac_whitelist.yml 配置的权限标识
				RBAC("custom:api:get"),
		)
		return se.Next()
	})

}
