package tools

import (
	"strings"

	"github.com/pocketbase/pocketbase/core"
)

// SetUserTenant 设置当前用户临时租户ID
func SetUserTenant(e *core.RequestEvent, tenantID string) {
	if e.Auth == nil {
		return
	}

	token := GetAuthTokenFromRequest(e)
	e.App.Store().Set("temp_tenant_id_for_user_"+token, tenantID)
}

// ClearUserTenant 清理当前用户临时租户ID
func ClearUserTenant(e *core.RequestEvent) {
	if e.Auth == nil {
		return
	}

	token := GetAuthTokenFromRequest(e)
	e.App.Store().Remove("temp_tenant_id_for_user_" + token)
}

// GetUserTenant 获取当前用户的租户ID
func GetUserTenant(e *core.RequestEvent) string {
	if e.Auth == nil {
		return ""
	}
	token := GetAuthTokenFromRequest(e)
	temp := e.App.Store().Get("temp_tenant_id_for_user_" + token)
	if tempStr, ok := temp.(string); ok && strings.TrimSpace(tempStr) != "" {
		return tempStr
	}
	return e.Auth.GetString("tenant_id")
}
