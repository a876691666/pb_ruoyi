package system

import (
	"pocketbase-ruoyi/tools"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// RegisterSystemTenant exposes system/tenant related custom endpoints.
// Adds: PUT /api/system/tenant/dynamic/{tenantId} to set a temporary tenant context for current user.
func RegisterSystemTenant(app *pocketbase.PocketBase) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// Set temporary tenant for current authenticated user
		se.Router.GET("/api/system/tenant/dynamic/{tenantId}", func(e *core.RequestEvent) error {
			ri, err := e.RequestInfo()
			if err != nil {
				return e.BadRequestError("请求信息错误", err)
			}
			if ri.Auth == nil {
				return e.UnauthorizedError("未登录或无权限", nil)
			}

			tenantID := e.Request.PathValue("tenantId")
			if tenantID == "" {
				return e.BadRequestError("缺少租户ID", nil)
			}

			// 校验租户是否存在
			if _, err := e.App.FindRecordById("tenant", tenantID); err != nil {
				return e.BadRequestError("租户不存在", err)
			}

			// 权限控制：
			// - 超级管理员可切换到任意租户
			// - 普通用户仅允许切换到自身租户
			if !isSuperuserByEvent(e) {
				selfTenant := ri.Auth.GetString("tenant_id")
				if selfTenant == "" || selfTenant != tenantID {
					return e.UnauthorizedError("无权切换到指定租户", nil)
				}
			}

			// 设置临时租户上下文
			tools.SetUserTenant(e, tenantID)
			return tools.JSONSuccess(e, map[string]any{
				"tenantId": tenantID,
			})
		})

		return se.Next()
	})
}

// isSuperuserByEvent 判断当前请求用户是否为超级管理员（包含 e.Auth.IsSuperuser 与 role_key=superadmin）
func isSuperuserByEvent(e *core.RequestEvent) bool {
	if e.Auth == nil {
		return false
	}
	if e.Auth.IsSuperuser() {
		return true
	}
	// Check role_key=superadmin for the current user
	count := 0
	q := e.App.DB().Select("count(*)").From("role").
		InnerJoin("user_role", dbx.NewExp("user_role.role = role.id")).
		Where(dbx.HashExp{
			"role.role_key":  "superadmin",
			"user_role.user": e.Auth.Id,
		})
	q.Row(&count)
	return count > 0
}
