package monitor

import (
	"fmt"
	"pocketbase-ruoyi/tools"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// RegisterMonitorOnline 注册 /api/monitor/online 相关接口
func RegisterMonitorOnline(app *pocketbase.PocketBase) {

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/api/monitor/online/list", func(e *core.RequestEvent) error {
			collection, _ := e.App.FindCollectionByNameOrId("users")
			result := []*core.Record{}
			userName := e.Request.URL.Query().Get("user_name")
			deptName := e.Request.URL.Query().Get("dept_name")

			q := app.RecordQuery("users").
				InnerJoin("_authOrigins", dbx.NewExp("_authOrigins.recordRef = users.id")).
				AndWhere(dbx.HashExp{"_authOrigins.collectionRef": collection.Id}).
				AndWhere(tools.BuildDataScopeExpression(e, "dept")).
				OrderBy("_authOrigins.created DESC")

			if userName != "" {
				q = q.AndWhere(dbx.NewExp(fmt.Sprintf("users.user_name LIKE '%%%s%%'", userName)))
			}
			if deptName != "" {
				q = q.AndWhere(dbx.NewExp(fmt.Sprintf("users.dept_name LIKE '%%%s%%'", deptName)))
			}

			q.All(&result)

			return tools.JSONSuccess(e, result)
		})
		se.Router.DELETE("/api/monitor/online/{user_id}", func(e *core.RequestEvent) error {
			collection, _ := e.App.FindCollectionByNameOrId("users")

			userID := e.Request.PathValue("user_id")

			record, err := app.FindRecordById("users", userID)
			if err != nil {
				return err
			}
			record.RefreshTokenKey()

			e.App.Save(record)

			e.App.DB().Delete("_authOrigins", dbx.HashExp{
				"collectionRef": collection.Id,
				"recordRef":     userID,
			}).Execute()

			return tools.JSONSuccess(e, true)
		})

		return se.Next()
	})
}
