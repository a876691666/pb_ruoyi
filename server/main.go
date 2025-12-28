package main

import (
	"log"
	"os"

	"pocketbase-ruoyi/api/auth"
	"pocketbase-ruoyi/api/monitor"
	"pocketbase-ruoyi/api/system"
	"pocketbase-ruoyi/api/system/menu"
	"pocketbase-ruoyi/api/tenant"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// 自定义在请求上下文中存放权限标识的 key
type ctxKey string

// LoginForm 用于绑定登录请求体（支持 json 与 form）
type LoginForm struct {
	UUID     string `json:"uuid" form:"uuid"`
	Code     string `json:"code" form:"code"`
	TenantID string `json:"tenantId" form:"tenantId"`
}

func main() {
	app := pocketbase.New()

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), true))
		return se.Next()
	})

	// 注册自定义 API 路由
	auth.RegisterRBAC(app)
	auth.RegisterAuth(app)
	monitor.RegisterMonitorLogininfor(app)
	monitor.RegisterMonitorOnline(app)
	menu.RegisterSystemMenu(app)
	system.RegisterSystemDept(app)
	system.RegisterSystemUserProfile(app)
	system.RegisterSystemRole(app)
	system.RegisterSystemDict(app)
	system.RegisterSystemOSS(app)
	system.RegisterSystemTenant(app)
	system.RegisterSystemImportExport(app)
	system.RegisterSystemCollections(app)

	tenant.RegisterTenant(app)

	auth.RegisterDataScope(app)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
