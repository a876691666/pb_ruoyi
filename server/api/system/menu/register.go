package menu

import (
	"pocketbase-ruoyi/tools"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// RegisterSystemMenu 注册 /api/system/menu/getRouters
func RegisterSystemMenu(app *pocketbase.PocketBase) {
	app.OnRecordDeleteExecute("menu").BindFunc(syncDeleteRoleMenu)
	app.OnRecordAfterDeleteSuccess("menu").BindFunc(syncMenuDeleteAfter)
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/api/system/menu/getRouters", func(e *core.RequestEvent) error {
			ri, err := e.RequestInfo()
			if err != nil {
				return e.BadRequestError("请求信息错误", err)
			}
			if ri.Auth == nil {
				return e.UnauthorizedError("未登录或无权限", nil)
			}

			menus := []Menu{}

			if tools.IsRoleSuperuser(app, ri.Auth.Id) {
				q := e.App.DB().Select("*").From("menu").
					Where(dbx.In("menu_type", "M", "C")).
					AndWhere(dbx.HashExp{"status": "0"}).
					OrderBy("parent_id ASC", "order_num ASC")
				if err := q.All(&menus); err != nil {
					return e.InternalServerError("查询菜单失败", err)
				}
			} else {
				uid := ri.Auth.Id

				q := e.App.DB().
					Select("m.*").
					From("menu as m").
					InnerJoin("role_menu as rm", dbx.NewExp("rm.menu = m.id")).
					InnerJoin("user_role as ur", dbx.NewExp("ur.role = rm.role")).
					Where(dbx.HashExp{"ur.user": uid}).
					AndWhere(dbx.In("m.menu_type", "M", "C")).
					AndWhere(dbx.HashExp{"m.status": "0"}).
					OrderBy("m.parent_id ASC", "m.order_num ASC")

				if err := q.All(&menus); err != nil {
					return e.InternalServerError("查询用户菜单失败", err)
				}
			}

			menus = uniqueMenus(menus)

			tree := buildMenuTree(menus)
			routers := buildRouters(tree)

			return tools.JSONSuccess(e, routers)
		})
		return se.Next()
	})
}
