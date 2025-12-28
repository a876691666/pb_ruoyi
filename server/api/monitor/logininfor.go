package monitor

import (
	"pocketbase-ruoyi/tools"

	"github.com/mileusna/useragent"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// UnlockPayload 请求体结构
type UnlockPayload struct {
	UserName string `json:"user_name" form:"user_name"`
}

// RegisterMonitorLogininfor 注册 /api/monitor/logininfor 相关接口
func RegisterMonitorLogininfor(app *pocketbase.PocketBase) {
	app.OnRecordAuthWithPasswordRequest().BindFunc(func(e *core.RecordAuthWithPasswordRequestEvent) error {
		err := e.Next()
		if e.Collection.Name == "_superusers" {
			return err
		}

		collection, _ := e.App.FindCollectionByNameOrId("logininfor")
		newRecord := core.NewRecord(collection)

		newRecord.Set("user_name", e.Identity)

		uaStr := e.Request.Header.Get("User-Agent")
		ua := useragent.Parse(uaStr)

		if ua.Mobile {
			newRecord.Set("client_key", "mobile")
			newRecord.Set("device_type", "mobile")
		} else {
			newRecord.Set("client_key", "pc")
			newRecord.Set("device_type", "pc")
		}

		newRecord.Set("browser", ua.Name)
		newRecord.Set("os", ua.OS)

		newRecord.Set("ipaddr", tools.GetIPAddr(e.Request))
		newRecord.Set("login_location", tools.GetLocationByIP(tools.GetIPAddr(e.Request)))

		if err == nil {
			newRecord.Set("status", "0")
			newRecord.Set("msg", "登录成功")

		} else {
			newRecord.Set("status", "1")
			newRecord.Set("msg", err.Error())
		}

		e.App.Save(newRecord)

		return err
	})

	app.OnRecordAuthWithPasswordRequest().BindFunc(func(e *core.RecordAuthWithPasswordRequestEvent) error {
		if e.Collection.Name == "_superusers" {
			return e.Next()
		}

		if !tools.VerifyCaptcha(e.RequestEvent) {
			return apis.NewBadRequestError("验证码错误或已过期", nil)
		}

		return e.Next()
	})

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// 解锁用户（示例接口）
		se.Router.POST("/api/monitor/logininfor/unlock/{userName}", func(e *core.RequestEvent) error {
			userName := e.Request.PathValue("userName")

			record, _ := e.App.FindFirstRecordByFilter("users", `username = {:userName}`, dbx.Params{"userName": userName})

			if record == nil {
				return tools.JSONSuccess(e, false)
			}

			record.SetVerified(true)

			if e.App.Save(record) != nil {
				return tools.JSONSuccess(e, false)
			}

			return tools.JSONSuccess(e, true)
		})

		se.Router.DELETE("/api/monitor/logininfor/clean", func(e *core.RequestEvent) error {
			collection, err := e.App.FindCollectionByNameOrId("logininfor")
			if err != nil {
				return tools.JSONSuccess(e, false)
			}

			if e.App.TruncateCollection(collection) != nil {
				return tools.JSONSuccess(e, false)
			}

			return tools.JSONSuccess(e, true)
		})

		return se.Next()
	})
}
