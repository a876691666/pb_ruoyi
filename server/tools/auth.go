package tools

import (
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// IsRoleSuperuser 检查用户是否拥有超级管理员角色
func IsRoleSuperuser(app *pocketbase.PocketBase, userID string) bool {
	if record, _ := app.FindFirstRecordByFilter(
		"user_role",
		"user = {:userID} && role = {:roleID}",
		dbx.Params{"userID": userID, "roleID": "1"}); record != nil {
		return true
	}
	return false
}

// GetAuthTokenFromRequest 从请求事件中提取认证令牌
func GetAuthTokenFromRequest(e *core.RequestEvent) string {
	token := e.Request.Header.Get("Authorization")
	if token != "" {
		// the schema prefix is not required and it is only for
		// compatibility with the defaults of some HTTP clients
		token = strings.TrimPrefix(token, "Bearer ")
	}
	return token
}
