package auth

import (
	"net/http"
	"time"

	"pocketbase-ruoyi/api/system"
	"pocketbase-ruoyi/api/system/menu"
	"pocketbase-ruoyi/tools"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// RegisterAuth 注册自定义的认证相关路由和钩子
func RegisterAuth(app *pocketbase.PocketBase) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/api/auth/code", func(e *core.RequestEvent) error {
			id, b64img, answer, err := tools.GenerateBase64Captcha(6, 120, 40)
			if err != nil {
				return e.JSON(http.StatusInternalServerError, map[string]any{
					"code": http.StatusInternalServerError,
					"msg":  err.Error(),
					"data": nil,
				})
			}

			tools.CacheSet(id, answer, 10*time.Minute)

			return tools.JSONSuccess(e, map[string]any{
				"captchaEnabled": true,
				"uuid":           id,
				"img":            b64img,
			})
		})
		return se.Next()
	})

	app.OnRecordAuthRequest().BindFunc(func(e *core.RecordAuthRequestEvent) error {

		isSuperAdmin := false
		userID := ""
		meta := map[string]any{}

		if e.Auth != nil {
			isSuperAdmin = e.Auth.IsSuperuser()
			userID = e.Auth.Id
		}

		if e.Record != nil {
			isSuperAdmin = e.Record.IsSuperuser()
			userID = e.Record.Id
		}

		if tools.IsRoleSuperuser(app, userID) {
			isSuperAdmin = true
		}

		if isSuperAdmin {
			meta["permissions"] = []string{"*:*:*"}
			meta["roles"] = []string{"superadmin"}
		} else if userID != "" {
			permissions := menu.GetAllPermissionsByUser(e.RequestEvent, userID)
			roles := system.GetAllRolesKeyByUser(e.RequestEvent, userID)

			meta["permissions"] = permissions
			meta["roles"] = roles
		}

		if meta["permissions"] != nil {
			e.Meta = meta
		}

		if e.Token != "" {
			app.Store().Set(e.Token, 1)
		}

		tools.ClearUserTenant(e.RequestEvent)

		return e.Next()
	})

}
